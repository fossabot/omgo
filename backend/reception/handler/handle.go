package handler

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

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

func Login(w http.ResponseWriter, r *http.Request) {
	var ret pc.S2CLoginRsp
	ret.Header = &pc.RspHeader{}
	setRspHeader(ret.Header)

	defer func() {
		js, err := json.Marshal(ret)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}()

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

		var userLoginReq pb.DB_UserLoginRequest
		userLoginReq.Info = &pc.UserBasicInfo{Email: email}
		userLoginReq.Secret = secretBytes

		cli := pb.NewDBServiceClient(conn)
		loginRsp, err := cli.UserLogin(context.Background(), &userLoginReq)
		if err != nil {
			ret.Header.Status = pc.ResultCode_RESULT_INVALID
			ret.Header.Msg = "login failed"
			log.Infof("login failed: %v", err)
			return
		}

		ret.UserInfo = loginRsp.GetInfo()
		ret.Token = loginRsp.GetToken()
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var ret pc.S2CLoginRsp
	ret.Header = &pc.RspHeader{}
	setRspHeader(ret.Header)

	defer func() {
		js, err := json.Marshal(ret)

		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}()

	//registerReq := &pb.DB_UserRegisterRequest{}
	//registerReq.
}
