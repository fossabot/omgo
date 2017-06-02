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
	header.Timestamp = uint64(time.Now().Unix())
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
	return u[:]
}

func (s *server) init(mcfg *mgo.DialInfo, rcfg *redisConfig) {
	s.driver.init(mcfg, rcfg)
}

// query user info
func (s *server) UserQuery(ctx context.Context, key *proto.DB_UserKey) (ret *proto.DB_UserQueryResponse, err error) {
	setRspHeader(ret.Result)

	err = regulateUserKey(key)
	if err != nil {
		ret.Result.Status = pc.ResultCode_RESULT_INVALID
		ret.Result.Msg = fmt.Sprintf("err:%v", err)
		return
	}

	userInfo, err := s.driver.queryUserBasicInfo(key)
	ret.Info = userInfo

	if err != nil {
		log.Errorf("error while query user:%v", err)
		ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		ret.Result.Msg = fmt.Sprintf("error:%v", err)
	}

	return
}

// update user info
func (s *server) UserUpdateInfo(ctx context.Context, userBasicInfo *pc.UserBasicInfo) (ret *pc.RspHeader, err error) {
	setRspHeader(ret)
	err = s.driver.updateUserInfoMongoDB(userBasicInfo)
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
func (s *server) UserRegister(ctx context.Context, request *proto.DB_UserRegisterRequest) (ret *proto.DB_UserRegisterResponse, err error) {
	setRspHeader(ret.Result)
	// check for existed user by email address
	email, valid := checkEmail(request.GetInfo().GetEmail())
	if !valid {
		ret.Result.Status = pc.ResultCode_RESULT_INVALID
		ret.Result.Msg = fmt.Sprintf("user:%v email invalid", request.Info)
		log.Info(ret.Result.Msg)
		return
	}
	userBasicInfo, err := s.driver.queryUserBasicInfo(proto.DB_UserKey{Email: email})
	if err != nil {
		log.Errorf("error while register user:%v", err)
		return
	}

	// user already existed
	if userBasicInfo.Usn != 0 {
		// email already registered
		ret.Result.Status = pc.ResultCode_RESULT_INVALID
		ret.Result.Msg = fmt.Sprintf("user:%v already registered", userBasicInfo)
		log.Info(ret.Result.Msg)
		return
	}

	// allocate new usn and uid
	usn, uid, err := s.driver.getUniqueID()
	if err != nil {
		ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		ret.Result.Msg = fmt.Sprintf("error while get uniqueID:%v", err)
		log.Errorf(ret.Result.Msg)
		return
	}

	userBasicInfo = request.Info
	userBasicInfo.Usn = usn
	userBasicInfo.Uid = uid
	userBasicInfo.Since = uint64(time.Now().Unix())
	userBasicInfo.Email = email
	if userBasicInfo.GetAvatar() == "" {
		userBasicInfo.Avatar = gravatarURL + utils.GetStringMD5Hash(email)
	}

	extra := &proto.DB_UserExtraInfo{Secret: request.Secret, Token: genToken()}
	s.driver.updateUserExtraMongoDB(usn, extra)
	s.driver.updateUserExtraRedis(usn, extra)
	s.driver.updateUserInfoMongoDB(userBasicInfo)
	s.driver.updateUserInfoRedis(userBasicInfo)

	return
}

// login
func (s *server) UserLogin(ctx context.Context, request *proto.DB_UserLoginRequest) (ret *proto.DB_UserLoginResponse, err error) {
	setRspHeader(ret.Result)
	// basic check
	if request.GetInfo().GetEmail() == "" || len(request.GetSecret()) == 0 {
		ret.Result.Status = pc.ResultCode_RESULT_INVALID
		log.Errorf("incoming invalid login request:%v", request)
		return
	}

	// query user
	userInfo, err := s.driver.queryUserBasicInfo(&proto.DB_UserKey{Email: request.GetInfo().GetEmail()})
	if err != nil {
		ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		log.Errorf("query user failed")
		return
	}

	// query user extra info
	userExtra, err := s.driver.queryUserExtraInfo(userInfo.Usn)
	if err != nil {
		ret.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		log.Errorf("query user extra failed")
		return
	}

	if bytes.Compare(userExtra.Secret, request.GetSecret()) != 0 {
		ret.Result.Status = pc.ResultCode_RESULT_INVALID
		log.Info("login with invalid credentials")
		return
	}

	// update token
	userExtra.Token = genToken()
	s.driver.updateUserExtraMongoDB(userInfo.Usn, userExtra)
	s.driver.updateUserExtraRedis(userInfo.Usn, userExtra)
	// update last time login
	userInfo.LastLogin = time.Now().Unix()
	s.driver.updateUserInfoMongoDB(userInfo)
	s.driver.updateUserInfoRedis(userInfo)

	return
}

// logout
func (s *server) UserLogout(ctx context.Context, request *proto.DB_UserLogoutRequest) (ret *pc.RspHeader, err error) {
	setRspHeader(ret)

	if request.Usn == 0 || len(request.GetToken()) == 0 {
		ret.Status = pc.ResultCode_RESULT_INVALID
		ret.Msg = "session invalid"
		return
	}

	userExtra, err := s.driver.queryUserExtraInfo(request.Usn)
	if err != nil {
		ret.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		ret.Msg = err
		return
	}
	if bytes.Compare(userExtra.GetToken(), request.GetToken()) != 0 {
		ret.Status = pc.ResultCode_RESULT_INVALID
		ret.Msg = "session invalid"
		return
	}
	s.driver.deleteUserExtraRedis(request.GetUsn())

	return
}
