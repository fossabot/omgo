package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

// gRPC
type server struct {
	driver driver
}

func (s *server) init(mcfg *mgo.DialInfo, rcfg *redisConfig) {
	s.driver.init(mcfg, rcfg)
}

func (s *server) QueryUser(ctx context.Context, in *proto.DB_UserKey) (ret *proto.DB_UserQueryResult, err error) {
	// get redis connection from pool
	conn := s.driver.redisClient.Get()

	// query user information
	var values interface{}
	switch {
	case in.Usn != 0:
		values, err = redis.Values(conn.Do("HGETALL", "user:"+in.Usn))
	case in.Uid != 0:

	}
	conn.Close()

	var userInfo proto_common.UserBasicInfo
	status := proto_common.ResultCode_RESULT_OK

	if err = redis.ScanStruct(values, &userInfo); err != nil {
		log.Error(err)
		status = proto_common.ResultCode_RESULT_INVALID
	}

	return &proto.DB_UserQueryResult{
		Status: status,
		Info:   &userInfo,
	}, err
}
