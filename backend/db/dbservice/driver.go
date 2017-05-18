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

var (
	mongoDBInvalidError = errors.New("no such db or collection")
)

// init both redis and mongodb client
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

// query user basic info in both redis and mongodb
func (d *driver) queryUserBasicInfo(key *proto.DB_UserKey) (*proto_common.UserBasicInfo, error) {
	var userInfo proto_common.UserBasicInfo
	var err error

	if key.Usn != 0 {
		// a valid usn, query in redis first
		err = d.queryUserBasicInfoRedis(key.Usn, &userInfo)
		if err == nil && userInfo.Usn == key.Usn {
			// found in redis
			return &userInfo, err
		}
	}

	// query in mongodb
	err = d.queryUserBasicInfoMongoDB(key, &userInfo)
	if err != nil {
		// found in mongodb, update to redis
		d.updateUserInfoRedis(&userInfo)
	}

	return &userInfo, err
}

func (d *driver) queryUserBasicInfoRedis(usn uint64, userInfo *proto_common.UserBasicInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("user:%v", usn)))
	if err == nil && len(values) > 0 {
		err = redis.ScanStruct(values, userInfo)
	}

	return err
}

func (d *driver) queryUserBasicInfoMongoDB(key *proto.DB_UserKey, userInfo *proto_common.UserBasicInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C("users")
	if c == nil {
		return mongoDBInvalidError
	}
	err := c.Find(bson.M{"usn": key.Usn, "email": key.Email, "uid": key.Uid}).One(userInfo)
	if err != nil {
		// not found in mongodb
		return err
	}

	return nil
}

// query user extra info in both redis and mongodb
func (d *driver) queryUserExtraInfo(usn uint64) (*proto.DB_UserExtraInfo, error) {
	var extraInfo proto.DB_UserExtraInfo
	var err error

	if usn != 0 {
		// a valid usn, query in redis first
		err = d.queryUserExtraRedis(usn, &extraInfo)
		if err == nil && len(extraInfo.Secret) != 0 {
			// found in redis
			return &extraInfo, err
		}
	}

	// query in mongodb
	err = d.queryUserExtraMongoDB(usn, &extraInfo)
	if err != nil {
		// found in mongodb, update to redis
		d.updateUserExtraRedis(usn, &extraInfo)
	}

	return &extraInfo, err
}

func (d *driver) queryUserExtraRedis(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("userExtra:%v", usn)))
	if err == nil && len(values) > 0 {
		err = redis.ScanStruct(values, extraInfo)
	}

	return err
}

func (d *driver) queryUserExtraMongoDB(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C("userExtra")
	if c == nil {
		return mongoDBInvalidError
	}
	err := c.Find(bson.M{"usn": usn}).One(extraInfo)
	if err != nil {
		// not found in mongodb
		return err
	}

	return nil
}

func (d *driver) updateUserExtraRedis(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	// store result to redis
	_, err := d.redisClient.Get().Do("HMSET", redis.Args{}.Add("userExtra:", usn).AddFlat(extraInfo))
	if err != nil {
		log.Error(err)
	}
	return err
}

// update user basic info in redis
func (d *driver) updateUserInfoRedis(userInfo *proto_common.UserBasicInfo) error {
	// store result to redis
	_, err := d.redisClient.Get().Do("HMSET", redis.Args{}.Add("user:", userInfo.Usn).AddFlat(userInfo))
	if err != nil {
		log.Error(err)
	}
	return err
}

func (d *driver) updateUserInfoMongoDB(userInfo *proto_common.UserBasicInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C("users")
	if c == nil {
		return mongoDBInvalidError
	}
	_, err := c.Upsert(bson.M{"usn": userInfo.Usn}, userInfo)
	if err != nil {
		log.Errorf("error while upsert userinfo:%v error:%v", userInfo, err)
	}

	return err
}
