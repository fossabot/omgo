package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
)

const (
	profileAddress      = "0.0.0.0:6666"
	defaultETCD         = "http://127.0.0.1:2379"
	defaultRoot         = "/backends"
	defaultListen       = ":8080"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
)

var (
	defaultServices = []string{"db"}
)

func setRspHeader(rsp *pc.RspHeader) *pc.RspHeader {
	rsp.Timestamp = utils.Timestamp()
	rsp.Status = pc.ResultCode_RESULT_OK
	return rsp
}

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	// profiling
	go http.ListenAndServe(profileAddress, nil)

	// cli
	app := &cli.App{
		Name:    "auth",
		Usage:   "a auth service",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"l"},
				Name:    "listen",
				Usage:   "listening address:port",
				Value:   defaultListen,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"e"},
				Name:    "etcdhosts",
				Usage:   "etcd hosts",
				Value:   cli.NewStringSlice(defaultETCD),
			},
			&cli.StringFlag{
				Aliases: []string{"r"},
				Name:    "etcdroot",
				Usage:   "services root path on etcd",
				Value:   defaultRoot,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"s"},
				Name:    "services",
				Usage:   "service names",
				Value:   cli.NewStringSlice(defaultServices...),
			},
			&cli.DurationFlag{
				Name:  "readtimeout",
				Usage: "seconds before reads timeout",
				Value: defaultReadTimeout,
			},
			&cli.DurationFlag{
				Name:  "writetimeout",
				Usage: "seconds before writes timeout",
				Value: defaultWriteTimeout,
			},
		},
		Action: func(c *cli.Context) error {
			listen := c.String("listen")
			etcdHosts := c.StringSlice("etcdhosts")
			etcdRoot := c.String("etcdroot")
			serviceNames := c.StringSlice("services")
			rt := c.Duration("readtimeout")
			wt := c.Duration("writetimeout")
			log.Println("listen:", listen)
			log.Println("etcdhosts:", etcdHosts)
			log.Println("etcdroot:", etcdRoot)
			log.Println("services:", serviceNames)
			log.Println("read timeout:", rt)
			log.Println("write timeout:", wt)

			startHTTP(listen, rt, wt)

			return nil
		},
	}

	app.Run(os.Args)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var ret pc.S2CLoginRsp
	ret.Header = &pc.RspHeader{}
	setRspHeader(ret.Header)

	email := r.Header.Get("email")
	pass := r.Header.Get("pass")

	log.Info("email:", email)
	log.Info("pass:", pass)

	if email == "" || pass == "" {
		ret.Header.Status = pc.ResultCode_RESULT_INVALID
		ret.Header.Msg = "invalid parameter(s)"
	} else {
		ret.UserInfo = &pc.UserBasicInfo{
			Usn:       utils.Timestamp(),
			Uid:       1234,
			Birthday:  0,
			Gender:    pc.Gender_GENDER_FEMALE,
			Nickname:  "wow",
			Email:     email,
			Avatar:    "https://www.gravatar.com/avatar/" + utils.GetStringMD5Hash(email) + "?s=200&r=pg&d=404",
			Country:   "cn",
			LastLogin: utils.Timestamp(),
		}
	}

	js, err := json.Marshal(ret)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func startHTTP(addr string, rt, wt time.Duration) {
	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: wt,
		ReadTimeout:  rt,
	}
	log.Fatal(srv.ListenAndServe())
}
