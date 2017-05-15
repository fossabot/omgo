package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	"github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
	"log"
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

func (s *server) QueryUser(ctx context.Context, in *proto.DB_UserKey) (*proto_common.UserBasicInfo, error) {
	return nil, nil
}
