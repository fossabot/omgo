package handler

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	pb "github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
	"golang.org/x/net/context"
)

const (
	pathSep = string(os.PathSeparator)
)

var (
	etcdClient   etcd.Client
	keyAgentETCD string
)

func setRspHeader(rsp *pc.RspHeader) *pc.RspHeader {
	rsp.Timestamp = utils.Timestamp()
	rsp.Status = pc.ResultCode_RESULT_OK
	return rsp
}

func responseFunc(w http.ResponseWriter, ret interface{}) {
	js, err := json.Marshal(ret)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getAppConfig() *pc.AppConfig {
	keysAPI := etcd.NewKeysAPI(etcdClient)
	resp, err := keysAPI.Get(context.Background(), keyAgentETCD, &etcd.GetOptions{Recursive: true})
	log.Infof("reading agent etcd config from:%v", keyAgentETCD)
	if err != nil {
		log.Error(err)
		return nil
	}

	if !resp.Node.Dir || len(resp.Node.Nodes) == 0 {
		log.Error("not a directory")
		return nil
	}

	ret := &pc.AppConfig{}

	for i, agent := range resp.Node.Nodes {
		ip, strPort, err := net.SplitHostPort(agent.Value)
		if err != nil {
			log.Errorf("error while parsing agent host:%v", agent.Value)
			break
		}
		port, err := strconv.ParseInt(strPort, 10, 32)
		if err != nil {
			log.Errorf("error while parsing agent port:%v", strPort)
		}
		ret.NetworkCfg = append(ret.NetworkCfg, &pc.NetworkConfig{
			Id:   int32(i),
			Desc: agent.Key,
			Ip:   ip,
			Port: int32(port),
		})
	}

	return ret
}

// Init handle with agent etcd config information
func Init(root string, endpoints []string, agent string) {
	cfg := etcd.Config{
		Endpoints: endpoints,
		Transport: etcd.DefaultTransport,
	}
	var err error
	etcdClient, err = etcd.New(cfg)
	if err != nil {
		log.Errorf("error while creating etcd client:%v", err)
	}
	keyAgentETCD = pathSep + root + pathSep + agent
}

// Login handles user login request
func Login(w http.ResponseWriter, r *http.Request) {
	ret := &pc.LoginRsp{}
	ret.Header = &pc.RspHeader{}
	setRspHeader(ret.Header)

	defer responseFunc(w, ret)

	log.Infof("processing login request from %v", r.RemoteAddr)

	email := r.Header.Get("email")
	password := r.Header.Get("password")

	log.Info("email:", email)
	log.Info("password:", password)

	secret := utils.GetStringSHA1Hash(password + defaultSalt)
	secretBytes, err := hex.DecodeString(secret)

	if email == "" || secret == "" || secretBytes == nil || err != nil {
		ret.Header.Status = pc.ResultCode_RESULT_INVALID
		ret.Header.Msg = "invalid parameter(s)"
	} else {
		conn := services.GetServiceWithID("dbservice", defaultDBSID)
		if conn == nil {
			ret.Header.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
			ret.Header.Msg = "interal error"
			log.Error("cannot get db service:", defaultDBSID)
			return
		}

		userLoginReq := &pb.DB_UserLoginRequest{}
		userLoginReq.Info = &pc.UserBasicInfo{Email: email}
		userLoginReq.Secret = secretBytes

		cli := pb.NewDBServiceClient(conn)
		loginRsp, err := cli.UserLogin(context.Background(), userLoginReq)
		if err != nil || loginRsp.Result.Status != pc.ResultCode_RESULT_OK {
			ret.Header.Status = pc.ResultCode_RESULT_INVALID
			ret.Header.Msg = "login failed"
			log.Infof("login failed: %v", err)
			return
		}

		ret.UserInfo = loginRsp.GetInfo()
		ret.Token = loginRsp.GetToken()
		ret.Config = getAppConfig()
	}
}

// Register handles user register request
func Register(w http.ResponseWriter, r *http.Request) {
	ret := &pc.LoginRsp{}
	ret.Header = &pc.RspHeader{}
	setRspHeader(ret.Header)

	defer responseFunc(w, ret)

	log.Infof("processing register request from %v", r.RemoteAddr)

	email := r.Header.Get("email")
	nick := r.Header.Get("nickname")
	password := r.Header.Get("password")
	birthday, err := strconv.ParseUint(r.Header.Get("birthday"), 10, 64)
	if err != nil {
		birthday = 0
	}
	gender := pc.Gender_GENDER_UNKNOWN
	genderValue, err := strconv.ParseInt(r.Header.Get("gender"), 10, 32)
	if err == nil {
		gender = pc.Gender(genderValue)
	}

	country := r.Header.Get("country")

	email = strings.TrimSpace(email)
	nick = strings.TrimSpace(nick)
	password = strings.TrimSpace(password)
	country = strings.TrimSpace(country)

	registerReq := &pb.DB_UserRegisterRequest{}
	registerReq.Info = &pc.UserBasicInfo{
		Email:    email,
		Birthday: birthday,
		Avatar:   r.Header.Get("avatar"),
		Nickname: nick,
		Gender:   gender,
		Country:  country,
	}

	secret := utils.GetStringSHA1Hash(password + defaultSalt)
	secretBytes, err := hex.DecodeString(secret)
	registerReq.Secret = secretBytes

	if email == "" || nick == "" || password == "" || country == "" {
		ret.Header.Status = pc.ResultCode_RESULT_INVALID
		ret.Header.Msg = "invalid parameter(s)"
		if err != nil {
			log.Errorf("error while register user:%v", err)
		}
	} else {
		conn := services.GetServiceWithID("dbservice", defaultDBSID)
		if conn == nil {
			ret.Header.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
			ret.Header.Msg = "interal error"
			log.Error("cannot get db service:", defaultDBSID)
			return
		}

		log.Infof("sending register request to db:%v", registerReq)

		cli := pb.NewDBServiceClient(conn)
		registerRsp, err := cli.UserRegister(context.Background(), registerReq)
		if err != nil || registerRsp.Result.Status != pc.ResultCode_RESULT_OK {
			ret.Header.Status = pc.ResultCode_RESULT_INVALID
			ret.Header.Msg = "register failed"
			log.Infof("register failed: %v", err)
			return
		}

		ret.UserInfo = registerRsp.Info
		ret.Token = registerRsp.Token
		ret.Config = getAppConfig()
	}
}
