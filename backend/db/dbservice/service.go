package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"time"
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

// query user info
func (s *server) UserQuery(ctx context.Context, key *proto.DB_UserKey) (*proto.DB_UserQueryResult, error) {
	var queryResult proto.DB_UserQueryResult
	setRspHeader(queryResult.Result)

	if key.Usn == 0 && key.Uid == 0 && key.Email == "" {
		queryResult.Result.Status = pc.ResultCode_RESULT_INVALID
		return queryResult, nil
	}

	userInfo, err := s.driver.queryUser(key)
	queryResult.Info = userInfo

	if err != nil {
		log.Errorf("error while query user:%v", err)
		queryResult.Result.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		queryResult.Result.Msg = fmt.Sprintf("error:%v", err)
	}

	return queryResult, nil
}

// update user info
func (s *server) UserUpdateInfo(ctx context.Context, userInfo *pc.UserBasicInfo) (*pc.RspHeader, error) {

}

// register
func (s *server) UserRegister(ctx context.Context, userInfo *pc.UserBasicInfo) (*proto.DB_UserLoginResult, error) {

}

// logout
func (s *server) UserLogout(ctx context.Context, key *proto.DB_UserKey) (*pc.RspHeader, error) {

}
