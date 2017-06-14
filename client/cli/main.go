package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/Pallinder/go-randomdata"
	log "github.com/Sirupsen/logrus"
	"github.com/abiosoft/ishell"
	"github.com/master-g/omgo/client/cli/session"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
)

var (
	address    string
	sess       *session.Session
	httpclient *http.Client
	apiHost    string
	loginRsp   pc.S2CLoginRsp
)

func init() {
	sess = session.NewSession("")
	httpclient = &http.Client{
		Timeout: time.Second * 3,
	}
	apiHost = "http://localhost:8080"
}

func getAddressFromLoginRsp(rsp *pc.S2CLoginRsp) string {
	if rsp != nil && rsp.Config != nil && len(rsp.Config.NetworkCfg) != 0 {
		ip := rsp.Config.NetworkCfg[0].Ip
		port := rsp.Config.NetworkCfg[0].Port
		addr := fmt.Sprintf("%v:%v", ip, port)
		log.Info(addr)
		return addr
	}
	return ""
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
		Name: "apihost",
		Help: "set api host address",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			// http address
			c.Printf("API host (%v):", apiHost)
			_apiHost := c.ReadLine()
			if _apiHost != "" {
				apiHost = _apiHost
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "conn",
		Help: "conn address:port",
		Func: func(c *ishell.Context) {
			sess.Close()
			if len(c.Args) > 0 {
				address = c.Args[0]
			} else {
				address = getAddressFromLoginRsp(&loginRsp)
			}
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
			// email
			c.Print("Email:")
			email := c.ReadLine()
			// pass
			c.Print("Password:")
			pass := strings.TrimSpace(c.ReadPassword())
			if pass == "" {
				log.Error("password invalid")
				return
			}

			// send request
			req, err := http.NewRequest("GET", apiHost+"/login", nil)
			if err != nil {
				log.Errorf("error while create http request:%v", err)
			}
			req.Header.Add("email", email)
			req.Header.Add("password", pass)
			resp, err := httpclient.Do(req)
			if err != nil {
				log.Errorf("error while sending request:%v", err)
			}

			json.NewDecoder(resp.Body).Decode(&loginRsp)

			log.Info(loginRsp)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "register",
		Help: "send register request to reception server",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			// email
			c.Print("Email:")
			email := c.ReadLine()
			// pass
			c.Print("Password:")
			pass := strings.TrimSpace(c.ReadPassword())
			if pass == "" {
				log.Error("password invalid")
				return
			}
			c.Print("Confirm :")
			confirm := strings.TrimSpace(c.ReadPassword())
			if strings.Compare(pass, confirm) != 0 {
				log.Error("confirm password is different from the first time")
				return
			}
			// Gender
			c.Print("Gender (U)nknow (F)emale (M)ale:")
			g := c.ReadLine()
			gender := pc.Gender_GENDER_UNKNOWN
			if g != "" {
				g = strings.ToLower(g)
				if g[0] == 'f' {
					gender = pc.Gender_GENDER_FEMALE
				} else {
					gender = pc.Gender_GENDER_MALE
				}
			}
			// nick
			nick := ""
			if gender == pc.Gender_GENDER_FEMALE {
				nick = randomdata.FirstName(randomdata.Female)
			} else {
				nick = randomdata.FirstName(randomdata.Male)
			}
			c.Printf("Nickname (%v):", nick)
			_nick := c.ReadLine()
			if _nick != "" {
				nick = _nick
			}
			// country
			country := randomdata.Country(randomdata.ThreeCharCountry)
			c.Printf("Country (%v):", country)
			_country := c.ReadLine()
			if _country != "" && len(_country) >= 3 {
				country = strings.ToUpper(_country)
				country = country[:3]
			}

			// send request
			req, err := http.NewRequest("GET", apiHost+"/register", nil)
			if err != nil {
				log.Errorf("error while create http request:%v", err)
			}
			req.Header.Add("email", email)
			req.Header.Add("password", pass)
			req.Header.Add("gender", strconv.FormatInt(int64(gender), 10))
			req.Header.Add("nickname", nick)
			req.Header.Add("country", country)

			resp, err := httpclient.Do(req)
			if err != nil {
				log.Errorf("error while sending request:%v", err)
			}

			json.NewDecoder(resp.Body).Decode(&loginRsp)

			log.Info(loginRsp)
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
