package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

type driver struct {
	redisClient  *redis.Pool
	mongoSession *mgo.Session
}

type redisConfig struct {
	host        string
	db          int
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
}

func (d *driver) init(minfo *mgo.DialInfo, rcfg *redisConfig) {
	// init mongodb client
	var err error
	d.mongoSession, err = mgo.DialWithInfo(minfo)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	d.mongoSession.SetMode(mgo.Monotonic)

	// init redis client with pool
	d.redisClient = &redis.Pool{
		MaxIdle:     rcfg.maxIdle,
		MaxActive:   rcfg.maxActive,
		IdleTimeout: rcfg.idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcfg.host)
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
			// select redis db
			c.Do("SELECT", rcfg.db)
			return c, nil
		},
	}
}

func (d *driver) queryUser(key *proto.DB_UserKey) (*proto_common.UserBasicInfo, error) {
	var userInfo proto_common.UserBasicInfo
	var err error

	if key.Usn != 0 {
		// a valid usn, query in redis first
		redisConn := d.redisClient.Get()
		values, err := redisConn.Do("HGETALL", fmt.Printf("user:%d", key.Usn))
		if err != nil {
			log.Error(err)
		}
		err = redis.ScanStruct(values, &userInfo)
		if err != nil {
			log.Error(err)
		}
	}

	return &userInfo, err
}
