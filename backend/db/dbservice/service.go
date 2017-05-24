package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"regexp"
	"strings"
	"time"
)

const (
	gravatarUrl = "http://www.gravatar.com/avatar/"
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

func (s *server) init(mcfg *mgo.DialInfo, rcfg *redisConfig) {
	s.driver.init(mcfg, rcfg)
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
			userBasicInfo.Avatar = gravatarUrl + utils.GetStringMD5Hash(email)
		}

		// TODO get a token here
		extra := &proto.DB_UserExtraInfo{Secret: request.Secret}
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
	return nil, nil
}

// logout
func (s *server) UserLogout(ctx context.Context, key *proto.DB_UserKey) (*pc.RspHeader, error) {
	return nil, nil
}
