package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	log "github.com/Sirupsen/logrus"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"gopkg.in/abiosoft/ishell.v2"
)

var (
	address    string
	sess       *Session
	httpclient *http.Client
	apiHost    string
	loginRsp   pc.LoginRsp
	shell      *ishell.Shell
)

func init() {
	sess = &Session{}
	httpclient = &http.Client{
		Timeout: time.Second * 3,
	}
	apiHost = "http://localhost:8080"
}

func getAddressFromLoginRsp(rsp *pc.LoginRsp) string {
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

	shell = ishell.New()
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
			if sess.IsFlagConnectedSet() {
				sess.Close()
			}

			if len(c.Args) > 0 {
				if strings.Compare(c.Args[0], "default") == 0 {
					address = utils.GetLocalIP() + ":8888"
				} else {
					address = c.Args[0]
				}
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
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
			} else {
				sess.Close()
				log.Info("disconnected from server")
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "heartbeat",
		Help: "sending heartbeat to server",
		Func: func(c *ishell.Context) {
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
				return
			}
			sess.Heartbeat()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "go",
		Help: "go through all tests",
		Func: func(c *ishell.Context) {
			if sess.IsFlagConnectedSet() {
				log.Error("already connected to server, disconnect first")
				return
			}
			sess.Connect(address)
			sess.Heartbeat()
			sess.ExchangeKey()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exchangekey",
		Help: "exchange public key with server",
		Func: func(c *ishell.Context) {
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
				return
			}
			sess.ExchangeKey()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "httplogin",
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

			if loginRsp.Header.Status != pc.ResultCode_RESULT_OK {
				log.Errorf("error while login:%v", loginRsp.Header.Msg)
				return
			}
			log.Infof("login success, token:%v", loginRsp.GetToken())
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

			if loginRsp.Header.Status != pc.ResultCode_RESULT_OK {
				log.Errorf("error while login:%v", loginRsp.Header.Msg)
				return
			}
			log.Infof("login success, token:%v", loginRsp.GetToken())
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "login to agent server",
		Func: func(c *ishell.Context) {
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
				return
			}
			if !sess.IsFlagEncryptedSet() {
				log.Error("need to exchange key first")
				return
			}
			if loginRsp.UserInfo != nil {
				sess.Usn = loginRsp.UserInfo.Usn
				sess.Token = loginRsp.Token
				sess.Login()
			} else {
				log.Error("error while try to login, user info is nil")
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "bye",
		Help: "offline",
		Func: func(c *ishell.Context) {
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
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
