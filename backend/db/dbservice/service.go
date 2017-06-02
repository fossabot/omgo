package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

const (
	gravatarURL = "http://www.gravatar.com/avatar/"
)

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// gRPC
type server struct {
	driver driver
}

func setRspHeader(header *pc.RspHeader) *pc.RspHeader {
	header.Status = pc.ResultCode_RESULT_OK
	header.Timestamp = time.Now().Unix()
	return header
}

func checkEmail(email string) (trimmed string, valid bool) {
	trimmed = strings.TrimSpace(email)
	valid = emailRegexp.MatchString(trimmed)
	return
}

func regulateUserKey(key *proto.DB_UserKey) error {
	err := errors.New("invalid user key")
	if key.Usn == 0 && key.Uid == 0 && key.Email == "" {
		return err
	}

	if key.Email != "" {
		ok := false
		key.Email, ok = checkEmail(key.Email)
		if !ok {
			return err
		}
	}

	return nil
}

func genToken() []byte {
	u := uuid.NewV4()
	return u
}

func (s *server) init(mcfg *mgo.DialInfo, rcfg *redisConfig) {
	s.driver.init(mcfg, rcfg)
}

// query user info
func (s *server) UserQuery(ctx context.Context, key *proto.DB_UserKey) (*proto.DB_UserQueryResponse, error) {
	var queryResult proto.DB_UserQueryResponse
	setRspHeader(queryResult.Result)

	err := regulateUserKey(key)
	if err != nil {
		queryResult.Result.Status = pc.ResultCode_RESULT_INVALID
		queryResult.Result.Msg = fmt.Sprintf("err:%v", err)
		return queryResult, nil
	}

	userInfo, err := s.driver.queryUserBasicInfo(key)
	queryResult.Info = userInfo

	if err != nil {
		log.Errorf("error while query user:%v", err)
		queryResult.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		queryResult.Result.Msg = fmt.Sprintf("error:%v", err)
	}

	return queryResult, nil
}

// update user info
func (s *server) UserUpdateInfo(ctx context.Context, userBasicInfo *pc.UserBasicInfo) (*pc.RspHeader, error) {
	var ret pc.RspHeader
	err := s.driver.updateUserInfoMongoDB(userBasicInfo)
	if err == nil {
		err = s.driver.updateUserInfoRedis(userBasicInfo)
	}

	if err != nil {
		log.Errorf("error while update user basic:%v", err)
		ret.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		ret.Msg = fmt.Sprintf("error:%v", err)
	}

	return ret, nil
}

// register
func (s *server) UserRegister(ctx context.Context, request *proto.DB_UserRegisterRequest) (*proto.DB_UserRegisterResponse, error) {
	var ret proto.DB_UserRegisterResponse
	for {
		// check for existed user by email address
		email, valid := checkEmail(request.GetInfo().GetEmail())
		if !valid {
			ret.Result.Status = pc.ResultCode_RESULT_INVALID
			ret.Result.Msg = fmt.Sprintf("user:%v email invalid", request.Info)
			log.Info(ret.Result.Msg)
			break
		}
		userBasicInfo, err := s.driver.queryUserBasicInfo(email)
		if err != nil {
			log.Errorf("error while register user:%v", err)
			break
		}

		// user already existed
		if userBasicInfo.Usn != 0 {
			// email already registered
			ret.Result.Status = pc.ResultCode_RESULT_INVALID
			ret.Result.Msg = fmt.Sprintf("user:%v already registered", userBasicInfo)
			log.Info(ret.Result.Msg)
			break
		}

		// allocate new usn and uid
		usn, uid, err := s.driver.getUniqueID()
		if err != nil {
			ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
			ret.Result.Msg = fmt.Sprintf("error while get uniqueID:%v", err)
			log.Errorf(ret.Result.Msg)
			break
		}

		userBasicInfo = request.Info
		userBasicInfo.Usn = usn
		userBasicInfo.Uid = uid
		userBasicInfo.Since = time.Now().Unix()
		userBasicInfo.Email = email
		if userBasicInfo.GetAvatar() == "" {
			userBasicInfo.Avatar = gravatarURL + utils.GetStringMD5Hash(email)
		}

		extra := &proto.DB_UserExtraInfo{Secret: request.Secret, Token: genToken()}
		s.driver.updateUserExtraMongoDB(usn, extra)
		s.driver.updateUserExtraRedis(usn, extra)
		s.driver.updateUserInfoMongoDB(userBasicInfo)
		s.driver.updateUserInfoRedis(userBasicInfo)

		break
	}

	return ret, nil
}

// login
func (s *server) UserLogin(ctx context.Context, request *proto.DB_UserLoginRequest) (*proto.DB_UserLoginResponse, error) {
	var ret proto.DB_UserLoginResponse
	setRspHeader(ret.Result)
	for {
		// basic check
		if request.GetInfo().GetEmail() == "" || len(request.GetSecret()) == 0 {
			ret.Result.Status = pc.ResultCode_RESULT_INVALID
			log.Errorf("incoming invalid login request:%v", request)
			break
		}

		// query user
		userInfo, err := s.driver.queryUserBasicInfo(&proto.DB_UserKey{Email: request.GetInfo().GetEmail()})
		if err != nil {
			ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
			log.Errorf("query user failed")
			break
		}

		// query user extra info
		userExtra, err := s.driver.queryUserExtraInfo(userInfo.Usn)
		if err != nil {
			ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
			log.Errorf("query user extra failed")
			break
		}

		if bytes.Compare(userExtra.Secret, request.GetSecret()) != 0 {
			ret.Result.Status = pc.ResultCode_RESULT_INVALID
			log.Info("login with invalid credentials")
			break
		}

		// update token
		userExtra.Token = genToken()
		s.driver.updateUserExtraMongoDB(userInfo.Usn, userExtra)
		s.driver.updateUserExtraRedis(userInfo.Usn, userExtra)
		// update last time login
		userInfo.LastLogin = time.Now().Unix()
		s.driver.updateUserInfoMongoDB(userInfo)
		s.driver.updateUserInfoRedis(userInfo)
		break
	}

	return ret, nil
}

// logout
func (s *server) UserLogout(ctx context.Context, key *proto.DB_UserKey) (*pc.RspHeader, error) {

	return nil, nil
}
