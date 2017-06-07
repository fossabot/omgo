package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/abiosoft/ishell"
	"github.com/master-g/omgo/client/cli/session"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
)

const salt = "japari"

var (
	address    string
	sess       *session.Session
	httpclient *http.Client
)

func init() {
	sess = session.NewSession("")
	httpclient = &http.Client{
		Timeout: time.Second * 3,
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	defer sess.Close()
	defer utils.PrintPanicStack()

	app := &cli.App{
		Name:    "client",
		Usage:   "a cli-client for testing omgo",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"a"},
				Name:    "address",
				Usage:   "connect to address:port",
				Value:   ":8888",
			},
		},
		Action: func(c *cli.Context) error {
			address = c.String("address")
			log.Info(address)
			return nil
		},
	}
	app.Run(os.Args)

	shell := ishell.New()
	shell.AddCmd(&ishell.Cmd{
		Name: "conn",
		Help: "conn address:port",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				address = c.Args[0]
			}
			sess.Close()
			sess.Connect(address)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "disconn",
		Help: "disconnect from server",
		Func: func(c *ishell.Context) {
			if !sess.IsConnected {
				c.Println("no connection")
			} else {
				sess.Close()
				c.Println("disconnected from server")
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "heartbeat",
		Help: "sending heartbeat to server",
		Func: func(c *ishell.Context) {
			if !sess.IsConnected {
				c.Println("no connection")
				return
			}
			sess.Heartbeat()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "go",
		Help: "go through all tests",
		Func: func(c *ishell.Context) {
			sess.Close()
			sess.Connect(address)
			sess.Heartbeat()
			sess.ExchangeKey()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exchangekey",
		Help: "exchange public key with server",
		Func: func(c *ishell.Context) {
			if !sess.IsConnected {
				c.Println("no connection")
				return
			}
			sess.ExchangeKey()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "send login request to reception server",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			// http address
			c.Print("API host (http://localhost:8080):")
			apiHost := c.ReadLine()
			if apiHost == "" {
				apiHost = "http://localhost:8080"
			}
			// email
			c.Print("Email:")
			email := c.ReadLine()
			// pass
			c.Print("Password:")
			pass := c.ReadPassword()

			afterPass := strings.TrimSpace(pass)
			if afterPass == "" {
				log.Error("password invalid")
				return
			}

			secret := utils.GetStringSHA1Hash(afterPass + salt)

			// send request
			req, err := http.NewRequest("GET", apiHost+"/login", nil)
			if err != nil {
				log.Errorf("error while create http request:%v", err)
			}
			req.Header.Add("email", email)
			req.Header.Add("secret", secret)
			resp, err := httpclient.Do(req)
			if err != nil {
				log.Errorf("error while sending request:%v", err)
			}

			var rsp pc.S2CLoginRsp
			json.NewDecoder(resp.Body).Decode(&rsp)

			log.Info(rsp)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "bye",
		Help: "offline",
		Func: func(c *ishell.Context) {
			if !sess.IsConnected {
				c.Println("no connection")
				return
			}
			sess.Bye()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "japari",
		Help: "welcome to japari park",
		Func: func(c *ishell.Context) {
			serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6666")
			if err != nil {
				log.Error(err)
			}
			localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
			if err != nil {
				log.Error(err)
			}
			udpConn, err := net.DialUDP("udp", localAddr, serverAddr)
			if err != nil {
				log.Error(err)
			}
			defer udpConn.Close()
			buf := []byte("hello")
			_, err = udpConn.Write(buf)
			if err != nil {
				log.Error(err)
			}
		},
	})

	shell.Start()
}
