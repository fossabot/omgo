package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	d.mongoSession.SetMode(mgo.Monotonic, true)

	// init redis client with pool
	d.redisClient = &redis.Pool{
		MaxIdle:     rcfg.maxIdle,
		MaxActive:   rcfg.maxActive,
		IdleTimeout: rcfg.idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcfg.host)
			if err != nil {
				log.Error(err)
			} else {
				// select redis db
				c.Do("SELECT", rcfg.db)
			}
			return c, err
		},
	}
}

func (d *driver) queryUser(key *proto.DB_UserKey) (*proto_common.UserBasicInfo, error) {
	var userInfo proto_common.UserBasicInfo
	var err error

	if key.Usn != 0 {
		// a valid usn, query in redis first
		err = d.queryUserInRedis(key.Usn, &userInfo)
		if err == nil && userInfo.Usn == key.Usn {
			// found in redis
			return &userInfo, err
		}
	}

	// query in mongodb
	err = d.queryUserInMongoDB(&key, &userInfo)
	if err != nil {
		// found in mongodb, update to redis
		d.updateUserInfoRedis(&userInfo)
	}

	return &userInfo, err
}

func (d *driver) queryUserInRedis(usn uint64, userInfo *proto_common.UserBasicInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", fmt.Printf("user:%d", usn)))
	if err == nil && len(values) > 0 {
		err = redis.ScanStruct(values, userInfo)
	}

	return err
}

func (d *driver) queryUserInMongoDB(key *proto.DB_UserKey, userInfo *proto_common.UserBasicInfo) error {
	c := d.mongoSession.DB("master").C("users")
	if c == nil {
		return errors.New("no such db or collection")
	}
	err := c.Find(bson.M{"usn": key.Usn, "email": key.Email, "uid": key.Uid}).One(userInfo)
	if err != nil {
		// no found in mongodb
		return err
	}

	return nil
}

func (d *driver) updateUserInfoRedis(userInfo *proto_common.UserBasicInfo) error {
	// store result to redis
	_, err := d.redisClient.Get().Do("HMSET", redis.Args{}.Add("user:", userInfo.Usn).AddFlat(userInfo))
	if err != nil {
		log.Error(err)
	}
	return err
}
