package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/master-g/omgo/backend/reception/handler"
	"github.com/master-g/omgo/services"
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

			services.Init(etcdRoot, etcdHosts, serviceNames)

			startHTTP(listen, rt, wt)

			return nil
		},
	}

	app.Run(os.Args)
}

func startHTTP(addr string, rt, wt time.Duration) {
	router := mux.NewRouter()
	router.HandleFunc("/login", handler.Login).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: wt,
		ReadTimeout:  rt,
	}
	log.Fatal(srv.ListenAndServe())
}
