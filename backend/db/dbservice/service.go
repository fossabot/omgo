package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
	"time"
)

// gRPC
type server struct {
	redisClient *redis.Pool
}

func (s *server) init(host string, db, maxIdle, maxActive int, idleTimeout time.Duration) {
	// init redis client with pool
	s.redisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Panic(err)
				return nil, err
			}
			// select redis db
			c.Do("SELECT", db)
			return c, nil
		},
	}
}

func (s *server) QueryUser(ctx context.Context, in *proto.DB_UserKey) (info *proto_common.UserBasicInfo, err error) {
	// get redis connection from pool
	conn := s.redisClient.Get()

	// query user information
	var values interface{}
	switch {
	case in.Usn != 0:
		values, err = conn.Do("HGETALL", fmt.Printf("user:%d", in.Usn))
	case in.Uid != 0:

	}
	conn.Close()

	if err = redis.ScanStruct(values, &info); err != nil {
		log.Error(err)
	}

	return info, err
}
