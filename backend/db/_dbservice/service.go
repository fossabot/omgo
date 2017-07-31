package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/db/dbservice/driver"
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
	letterSet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// gRPC
type server struct {
	db database
}

func setRspHeader(header *pc.RspHeader) *pc.RspHeader {
	header.Status = int32(pc.ResultCode_RESULT_OK)
	header.Timestamp = utils.Timestamp()
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

func genToken() string {
	raw := uuid.NewV4()
	token := make([]byte, len(raw))
	for i, v := range raw {
		token[i] = letterSet[int(v)%len(letterSet)]
	}
	return string(token)
}

func (s *server) init(mcfg *mgo.DialInfo, concurrent int, timeout time.Duration, rcfg *driver.Config) {
	s.db.init(mcfg, concurrent, timeout, rcfg)
}

// query user info
func (s *server) UserQuery(ctx context.Context, key *proto.DB_UserKey) (ret *proto.DB_UserQueryResponse, err error) {
	ret = &proto.DB_UserQueryResponse{}
	ret.Result = &pc.RspHeader{}
	setRspHeader(ret.Result)

	err = regulateUserKey(key)
	if err != nil {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Result.Msg = fmt.Sprintf("err:%v", err)
		return
	}

	userInfo, err := s.db.queryUserInfo(key)
	ret.Info = userInfo

	if err != nil {
		log.Errorf("error while query user:%v", err)
		ret.Result.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Result.Msg = fmt.Sprintf("error:%v", err)
	}

	return
}

// update user info
func (s *server) UserUpdateInfo(ctx context.Context, userInfo *pc.UserInfo) (ret *pc.RspHeader, err error) {
	ret = &pc.RspHeader{}
	setRspHeader(ret)
	err = s.db.updateUserInfoMongoDB(userInfo)
	if err == nil {
		err = s.db.updateUserInfoRedis(userInfo)
	}

	if err != nil {
		log.Errorf("error while update user basic:%v", err)
		ret.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Msg = fmt.Sprintf("error:%v", err)
	}

	return ret, nil
}

// register
func (s *server) UserRegister(ctx context.Context, request *proto.DB_UserRegisterRequest) (ret *proto.DB_UserRegisterResponse, err error) {
	log.Infof("user register:%v", request)

	ret = &proto.DB_UserRegisterResponse{}
	ret.Result = &pc.RspHeader{}
	setRspHeader(ret.Result)
	// check for existed user by email address
	email, valid := checkEmail(request.GetInfo().GetEmail())
	if !valid {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Result.Msg = fmt.Sprintf("user:%v email invalid", request.Info)
		log.Info(ret.Result.Msg)
		return
	}
	ret.Info, err = s.db.queryUserInfo(&proto.DB_UserKey{Email: email})
	// user already existed
	if ret.Info.Usn != 0 {
		// email already registered
		ret.Result.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Result.Msg = fmt.Sprintf("user:%v already registered", ret.Info.Email)
		ret.Info = nil
		log.Info(ret.Result.Msg)
		return
	}

	// allocate new usn and uid
	usn, uid, err := s.db.getUniqueID()
	if err != nil {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Result.Msg = fmt.Sprintf("error while get uniqueID:%v", err)
		log.Errorf(ret.Result.Msg)
		return
	}

	ret.Info = request.Info
	ret.Info.Usn = usn
	ret.Info.Uid = uid
	ret.Info.Since = utils.Timestamp()
	ret.Info.Email = email
	ret.Info.LastLogin = utils.Timestamp()

	if ret.Info.GetAvatar() == "" {
		ret.Info.Avatar = gravatarURL + utils.GetStringMD5Hash(email)
	}

	token := genToken()
	ret.Token = token

	extra := &proto.DB_UserExtraInfo{Usn: usn, Secret: request.Secret, Token: token}
	s.db.updateUserExtraMongoDB(extra)
	s.db.updateUserExtraRedis(extra)
	s.db.updateUserInfoMongoDB(ret.Info)
	s.db.updateUserInfoRedis(ret.Info)

	return
}

// login
func (s *server) UserLogin(ctx context.Context, request *proto.DB_UserLoginRequest) (ret *proto.DB_UserLoginResponse, err error) {
	log.Infof("user login:%v", request)

	ret = &proto.DB_UserLoginResponse{}
	ret.Result = &pc.RspHeader{}
	setRspHeader(ret.Result)
	// basic check
	if request.GetInfo().GetEmail() == "" || len(request.GetSecret()) == 0 {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INVALID)
		log.Errorf("incoming invalid login request:%v", request)
		return
	}

	// query user
	ret.Info, err = s.db.queryUserInfo(&proto.DB_UserKey{Email: request.GetInfo().GetEmail()})
	if err != nil {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Info = nil
		log.Errorf("query user failed")
		return
	}

	// query user extra info
	userExtra, err := s.db.queryUserExtraInfo(ret.Info.Usn)
	if err != nil {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Info = nil
		log.Errorf("query user extra failed:%v", ret.Info)
		return
	}

	if bytes.Compare(userExtra.Secret, request.GetSecret()) != 0 {
		ret.Result.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Info = nil
		log.Info("login with invalid credentials")
		return
	}

	token := genToken()
	ret.Token = token

	// update token
	userExtra.Token = token
	s.db.updateUserExtraMongoDB(userExtra)
	s.db.updateUserExtraRedis(userExtra)
	// update last time login
	ret.Info.LastLogin = utils.Timestamp()
	s.db.updateUserInfoMongoDB(ret.Info)
	s.db.updateUserInfoRedis(ret.Info)

	return
}

// logout
func (s *server) UserLogout(ctx context.Context, request *proto.DB_UserLogoutRequest) (ret *pc.RspHeader, err error) {
	log.Infof("user logout:%v", request)

	ret = &pc.RspHeader{}
	setRspHeader(ret)

	if request.Usn == 0 || len(request.GetToken()) == 0 {
		ret.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Msg = "session invalid"
		return
	}

	userExtra, err := s.db.queryUserExtraInfo(request.Usn)
	if err != nil {
		log.Errorf("unable to find user extra info:%v", err)
		ret.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		ret.Msg = "internal error"
		return
	}
	if strings.Compare(userExtra.GetToken(), request.GetToken()) != 0 {
		ret.Status = int32(pc.ResultCode_RESULT_INVALID)
		ret.Msg = "session invalid"
		return
	}
	s.db.deleteUserExtraRedis(request.GetUsn())

	return
}

// verify token
func (s *server) UserExtraInfoQuery(ctx context.Context, request *proto.DB_UserKey) (ret *proto.DB_UserExtraInfo, err error) {
	log.Infof("query user extra info:%v", request.GetUsn())

	ret = &proto.DB_UserExtraInfo{}
	if request.GetUsn() == 0 {
		log.Error("invalid user serial number")
		return
	}
	extraInfo, err := s.db.queryUserExtraInfo(request.GetUsn())
	if err != nil {
		log.Errorf("error while query user extra info:%v", err)
		return
	}

	ret.Usn = extraInfo.GetUsn()
	ret.Token = extraInfo.GetToken()
	ret.Secret = extraInfo.GetSecret()
	return
}