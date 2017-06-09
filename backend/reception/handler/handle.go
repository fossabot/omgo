package handler

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/proto/grpc/db"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
	"golang.org/x/net/context"
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

// Login handles user login request
func Login(w http.ResponseWriter, r *http.Request) {
	ret := &pc.S2CLoginRsp{}
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
	}
}

// Register handles user register request
func Register(w http.ResponseWriter, r *http.Request) {
	ret := &pc.S2CLoginRsp{}
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
	}
}
