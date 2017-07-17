package driver

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type MongoDriver struct {
	session *mgo.Session
	latch   chan *mgo.Session
}

func (d *MongoDriver) Init(dialInfo *mgo.DialInfo, concurrent int, timeout time.Duration) {
	d.latch = make(chan *mgo.Session, concurrent)
	sess, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	sess.SetMode(mgo.Monotonic, true)
	sess.SetSocketTimeout(timeout)
	sess.SetCursorTimeout(0)
	d.session = sess

	for i := 0; i < cap(d.latch); i++ {
		d.latch <- sess.Copy()
	}
}

func (d *MongoDriver) Execute(f func(sess *mgo.Session) error) error {
	sess := <-d.latch
	defer func() {
		d.latch <- sess
	}()
	sess.Refresh()
	return f(sess)
}
