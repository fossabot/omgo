package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/kit/utils"
	pc "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/abiosoft/ishell.v2"
)

type UserInfoExt struct {
	Usn        uint64        `json:"usn,omitempty"`
	Uid        uint64        `json:"uid"`
	Avatar     string        `json:"avatar,omitempty"`
	Birthday   uint64        `json:"birthday,omitempty"`
	Country    string        `json:"country,omitempty"`
	Email      string        `json:"email,omitempty"`
	Gender     pc.Gender     `json:"gender"`
	LastLogin  uint64        `json:"last_login,omitempty"`
	LoginCount int32         `json:"login_count"`
	Nickname   string        `json:"nickname,omitempty"`
	Since      uint64        `json:"since,omitempty"`
	Status     pc.UserStatus `json:"status"`
	Token      string        `json:"token"`
}

type HttpLoginRsp struct {
	UserInfo  UserInfoExt `json:"user_info,omitempty"`
	Timestamp uint64      `json:"timestamp,omitempty"`
}

var (
	address      string
	sess         *Session
	httpclient   *http.Client
	apiHost      string
	shell        *ishell.Shell
	httpLoginRsp HttpLoginRsp
)

func init() {
	sess = &Session{}
	httpclient = &http.Client{
		Timeout: time.Second * 3,
	}
	apiHost = "http://127.0.0.1:8080/api"
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
				address = ":8888"
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
			sess.Handshake()
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
			if email == "" {
				email = "tester@acme.com"
			}
			// pass
			c.Print("Password:")
			pass := strings.TrimSpace(c.ReadPassword())
			if pass == "" {
				pass = "123456"
			}
			// send request
			req, err := http.NewRequest("GET", apiHost+"/login", nil)
			if err != nil {
				log.Errorf("error while create http request:%v", err)
				return
			}
			req.Header.Add("email", email)
			req.Header.Add("password", pass)
			req.Header.Add("nonce", strconv.FormatUint(utils.Timestamp(), 10))
			req.Header.Add("Content-Type", "application/json")
			//req.Header.Add("Accept-Encoding", "application/json")
			resp, err := httpclient.Do(req)
			defer resp.Body.Close()
			if err != nil {
				log.Errorf("error while sending request:%v", err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("error while parsing response:%v", err)
				return
			}

			err = json.Unmarshal(body, &httpLoginRsp)
			if err != nil {
				log.Errorf("error while parsing http login response:%v", err)
				return
			}

			log.Infof("login success:%v", httpLoginRsp)
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
			email = strings.TrimSpace(email)
			email = strings.ToLower(email)
			if email == "" {
				log.Errorf("email invalid")
				return
			}
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
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("nonce", strconv.FormatUint(utils.Timestamp(), 10))
			req.Header.Add("app_language", "zh-rCN")
			req.Header.Add("app_version", "0.0.1")
			req.Header.Add("avatar", "http://gravatar.com/avatar/"+utils.GetStringMD5Hash(email))
			req.Header.Add("birthday", "531262800000")
			req.Header.Add("country", country)
			req.Header.Add("device_type", "1")
			req.Header.Add("email", email)
			req.Header.Add("gender", strconv.FormatInt(int64(gender), 10))
			req.Header.Add("mcc", "460")
			req.Header.Add("nickname", nick)
			req.Header.Add("os", "macOS High Sierra v10.13")
			req.Header.Add("os_locale", "zh-rCN")
			req.Header.Add("phone", "1234567890")
			req.Header.Add("secret", pass)
			req.Header.Add("timezone", "8")

			resp, err := httpclient.Do(req)
			defer resp.Body.Close()
			if err != nil {
				log.Errorf("error while sending request:%v", err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("error while parsing response:%v", err)
				return
			}

			err = json.Unmarshal(body, &httpLoginRsp)
			if err != nil {
				log.Errorf("error while parsing http register response:%v", err)
				return
			}

			log.Infof("register success:%v", httpLoginRsp)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "handshake",
		Help: "handshake to agent server",
		Func: func(c *ishell.Context) {
			if !sess.IsFlagConnectedSet() {
				log.Error("no connection")
				return
			}
			if httpLoginRsp.UserInfo.Usn == 0 {
				log.Error("do http login first")
				return
			}
			sess.Usn = httpLoginRsp.UserInfo.Usn
			sess.Token = httpLoginRsp.UserInfo.Token
			sess.Handshake()
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
