package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// DBUserStatus is used as a database user status index
type DBUserStatus struct {
	key string
	usn uint64
	uid uint64
}

var (
	errMongoDBInvalid = errors.New("no such db or collection")
)

const (
	// 0.5 day
	expireDuration = 60 * 60 * 12
	keyUser        = "user"
	keyUserExtra   = "userExtra"
)

func redisKey(key string, usn uint64) string {
	return fmt.Sprintf("%v:%v", key, usn)
}

func redisFlat(key string, value interface{}) redis.Args {
	args := redis.Args{}.Add(key).AddFlat(value)
	return args
}

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

////////////////////////////////////////////////////////////////////////////////
// usn, uid
////////////////////////////////////////////////////////////////////////////////
func (d *driver) getUniqueID() (usn, uid uint64, err error) {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	usn = 0
	uid = 0

	c := sessionCpy.DB("master").C("status")
	if c == nil {
		err = errMongoDBInvalid
		return
	}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"usn": 1, "uid": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	dbStatus := DBUserStatus{}
	_, err = c.Find(bson.M{"key": keyUser}).Apply(change, &dbStatus)
	if err != nil {
		// not found in mongodb
		return
	}

	usn = dbStatus.usn
	uid = dbStatus.uid

	return
}

////////////////////////////////////////////////////////////////////////////////
// Basic Info
////////////////////////////////////////////////////////////////////////////////

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
	if err == nil {
		// found in mongodb, update to redis
		d.updateUserInfoRedis(&userInfo)
	}

	return &userInfo, err
}

func (d *driver) queryUserBasicInfoRedis(usn uint64, userInfo *proto_common.UserBasicInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", redisKey(keyUser, usn)))
	if err == nil && len(values) > 0 {
		err = redis.ScanStruct(values, userInfo)
	}

	return err
}

func (d *driver) queryUserBasicInfoMongoDB(key *proto.DB_UserKey, userInfo *proto_common.UserBasicInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	index := mgo.Index{
		Key:        []string{"usn", "uid", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	c := sessionCpy.DB("master").C(keyUser)
	if c == nil {
		return errMongoDBInvalid
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Error(err)
	}
	err = c.Find(bson.M{"usn": key.Usn, "email": key.Email, "uid": key.Uid}).One(userInfo)
	if err != nil {
		// not found in mongodb
		return err
	}

	return nil
}

// update user basic info in redis
func (d *driver) updateUserInfoRedis(userInfo *proto_common.UserBasicInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	key := redisKey(keyUser, userInfo.Usn)
	// store result to redis
	_, err := conn.Do("HMSET", redisFlat(key, userInfo)...)
	if err != nil {
		log.Error(err)
	}
	_, err = conn.Do("EXPIRE", key, expireDuration)
	if err != nil {
		log.Error(err)
	}

	return err
}

func (d *driver) updateUserInfoMongoDB(userInfo *proto_common.UserBasicInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C(keyUser)
	if c == nil {
		return errMongoDBInvalid
	}
	_, err := c.Upsert(bson.M{"usn": userInfo.Usn}, userInfo)
	if err != nil {
		log.Errorf("error while upsert userinfo:%v error:%v", userInfo, err)
	}

	return err
}

func (d *driver) deleteUserInfoRedis(usn uint64) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redisKey(keyUser, usn))
	if err != nil {
		log.Errorf("error while removing userinfo:%v error%v", usn, err)
	}

	return err
}

////////////////////////////////////////////////////////////////////////////////
// Extra Info
////////////////////////////////////////////////////////////////////////////////

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

	values, err := redis.Values(conn.Do("HGETALL", redisKey(keyUserExtra, usn)))
	if err == nil && len(values) > 0 {
		err = redis.ScanStruct(values, extraInfo)
	}

	return err
}

func (d *driver) queryUserExtraMongoDB(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C(keyUserExtra)
	if c == nil {
		return errMongoDBInvalid
	}
	err := c.Find(bson.M{"usn": usn}).One(extraInfo)
	if err != nil {
		// not found in mongodb
		return err
	}

	return nil
}

func (d *driver) updateUserExtraRedis(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	conn := d.redisClient.Get()
	defer conn.Close()
	// store result to redis
	key := redisKey(keyUserExtra, usn)
	_, err := conn.Do("HMSET", redisFlat(key, extraInfo)...)
	if err != nil {
		log.Error(err)
	}
	_, err = conn.Do("EXPIRE", key, time.Hour.Seconds()*12)
	if err != nil {
		log.Error(err)
	}

	return err
}

func (d *driver) updateUserExtraMongoDB(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	sessionCpy := d.mongoSession.Copy()
	defer sessionCpy.Close()

	c := sessionCpy.DB("master").C(keyUserExtra)
	if c == nil {
		return errMongoDBInvalid
	}
	_, err := c.Upsert(bson.M{"usn": usn}, extraInfo)
	if err != nil {
		log.Errorf("error while upsert userextra:%v error:%v", extraInfo, err)
	}

	return err
}

func (d *driver) deleteUserExtraRedis(usn uint64) error {
	conn := d.redisClient.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redisKey(keyUserExtra, usn))
	if err != nil {
		log.Errorf("error while removing userextra:%v error%v", usn, err)
	}

	return err
}
