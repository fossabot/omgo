package main

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	driver "github.com/master-g/omgo/backend/db/dbservice/driver"
	"github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DBUserStatus is used as a database user status index
type DBUserStatus struct {
	Key string `json:"key,omitempty"`
	Usn uint64 `json:"usn,omitempty"`
	Uid uint64 `json:"uid,omitempty"`
}

var (
	errMongoDBInvalid = errors.New("no such db or collection")
)

const (
	// 0.5 day
	expireDuration = 60 * 60 * 12
)

func redisKey(key string, usn uint64) string {
	return fmt.Sprintf("%v:%v", key, usn)
}

func redisFlat(key string, value interface{}) redis.Args {
	args := redis.Args{}.Add(key).AddFlat(value)
	return args
}

type database struct {
	mongo driver.MongoDriver
	redis driver.RedisDriver
}

func (db *database) init(dialInfo *mgo.DialInfo, concurrent int, timeout time.Duration, cfg *driver.Config) {
	db.mongo.Init(dialInfo, concurrent, timeout)
	db.redis.Init(cfg)
}

////////////////////////////////////////////////////////////////////////////////
// usn, uid
////////////////////////////////////////////////////////////////////////////////
func (db *database) getUniqueID() (usn, uid uint64, err error) {
	err = db.mongo.Execute(func(sess *mgo.Session) error {
		usn = 0
		uid = 0

		c := sess.DB("master").C("status")
		if c == nil {
			return errMongoDBInvalid
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
			return err
		}

		log.Info(dbStatus)

		usn = dbStatus.Usn
		uid = dbStatus.Uid

		return nil
	})

	if err != nil {
		log.Errorf("error while getUniqueID:%v", err)
	}

	return
}

////////////////////////////////////////////////////////////////////////////////
// Basic Info
////////////////////////////////////////////////////////////////////////////////

// query user basic info in both redis and mongodb
func (db *database) queryUserBasicInfo(key *proto.DB_UserKey) (*pc.UserBasicInfo, error) {
	var userInfo pc.UserBasicInfo
	var err error

	if key.Usn != 0 {
		// a valid usn, query in redis first
		err = db.queryUserBasicInfoRedis(key.Usn, &userInfo)
		if err == nil && userInfo.Usn == key.Usn {
			// found in redis
			return &userInfo, err
		}
	}

	// query in mongodb
	err = db.queryUserBasicInfoMongoDB(key, &userInfo)
	if err == nil {
		// found in mongodb, update to redis
		db.updateUserInfoRedis(&userInfo)
	}

	return &userInfo, err
}

func (db *database) queryUserBasicInfoRedis(usn uint64, userInfo *pc.UserBasicInfo) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
		values, err := redis.Values(conn.Do("HGETALL", redisKey(keyUser, usn)))
		if err == nil && len(values) > 0 {
			err = redis.ScanStruct(values, userInfo)
		}

		return err
	})

	return err
}

func (db *database) queryUserBasicInfoMongoDB(key *proto.DB_UserKey, userInfo *pc.UserBasicInfo) error {
	err := db.mongo.Execute(func(sess *mgo.Session) error {
		index := mgo.Index{
			Key:        []string{"usn", "uid", "email"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}

		c := sess.DB("master").C(keyUser)
		if c == nil {
			return errMongoDBInvalid
		}
		err := c.EnsureIndex(index)
		if err != nil {
			log.Error(err)
		}
		query := &bson.M{}
		switch {
		case key.Usn != 0:
			query = &bson.M{"usn": key.Usn}
			break
		case key.Uid != 0:
			query = &bson.M{"uid": key.Uid}
			break
		case key.Email != "":
			query = &bson.M{"email": key.Email}
			break
		}
		err = c.Find(query).One(userInfo)
		if err != nil {
			// not found in mongodb
			return err
		}

		return nil
	})

	return err
}

// update user basic info in redis
func (db *database) updateUserInfoRedis(userInfo *pc.UserBasicInfo) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
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
	})

	return err
}

func (db *database) updateUserInfoMongoDB(userInfo *pc.UserBasicInfo) error {
	err := db.mongo.Execute(func(sess *mgo.Session) error {
		c := sess.DB("master").C(keyUser)
		if c == nil {
			return errMongoDBInvalid
		}
		_, err := c.Upsert(bson.M{"usn": userInfo.Usn}, userInfo)
		if err != nil {
			log.Errorf("error while upsert userinfo:%v error:%v", userInfo, err)
		}

		return err
	})
	return err
}

func (db *database) deleteUserInfoRedis(usn uint64) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
		_, err := conn.Do("DEL", redisKey(keyUser, usn))
		if err != nil {
			log.Errorf("error while removing userinfo:%v error%v", usn, err)
		}

		return err
	})
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Extra Info
////////////////////////////////////////////////////////////////////////////////

// query user extra info in both redis and mongodb
func (db *database) queryUserExtraInfo(usn uint64) (*proto.DB_UserExtraInfo, error) {
	var extraInfo proto.DB_UserExtraInfo
	var err error

	if usn != 0 {
		// a valid usn, query in redis first
		err = db.queryUserExtraRedis(usn, &extraInfo)
		if err == nil && len(extraInfo.Secret) != 0 {
			// found in redis
			return &extraInfo, err
		}
	}

	// query in mongodb
	err = db.queryUserExtraMongoDB(usn, &extraInfo)
	if err != nil {
		// found in mongodb, update to redis
		db.updateUserExtraRedis(&extraInfo)
	}

	return &extraInfo, err
}

func (db *database) queryUserExtraRedis(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
		values, err := redis.Values(conn.Do("HGETALL", redisKey(keyUserExtra, usn)))
		if err == nil && len(values) > 0 {
			err = redis.ScanStruct(values, extraInfo)
		}

		return err
	})
	return err
}

func (db *database) queryUserExtraMongoDB(usn uint64, extraInfo *proto.DB_UserExtraInfo) error {
	err := db.mongo.Execute(func(sess *mgo.Session) error {
		c := sess.DB("master").C(keyUserExtra)
		if c == nil {
			return errMongoDBInvalid
		}
		err := c.Find(bson.M{"usn": usn}).One(extraInfo)
		if err != nil {
			// not found in mongodb
			return err
		}

		return nil
	})
	return err
}

func (db *database) updateUserExtraRedis(extraInfo *proto.DB_UserExtraInfo) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
		// store result to redis
		key := redisKey(keyUserExtra, extraInfo.Usn)
		_, err := conn.Do("HMSET", redisFlat(key, extraInfo)...)
		if err != nil {
			log.Error(err)
		}
		_, err = conn.Do("EXPIRE", key, time.Hour.Seconds()*12)
		if err != nil {
			log.Error(err)
		}

		return err
	})
	return err
}

func (db *database) updateUserExtraMongoDB(extraInfo *proto.DB_UserExtraInfo) error {
	err := db.mongo.Execute(func(sess *mgo.Session) error {
		c := sess.DB("master").C(keyUserExtra)
		if c == nil {
			return errMongoDBInvalid
		}
		_, err := c.Upsert(bson.M{"usn": extraInfo.Usn}, extraInfo)
		if err != nil {
			log.Errorf("error while upsert userextra:%v error:%v", extraInfo, err)
		}

		return err
	})
	return err
}

func (db *database) deleteUserExtraRedis(usn uint64) error {
	err := db.redis.Execute(func(conn redis.Conn) error {
		_, err := conn.Do("DEL", redisKey(keyUserExtra, usn))
		if err != nil {
			log.Errorf("error while removing userextra:%v error%v", usn, err)
		}
		return err
	})
	return err
}
