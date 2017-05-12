package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/master-g/omgo/proto/pb"
	"github.com/master-g/omgo/utils"
	"net/http"
	"time"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	pass := r.Header.Get("pass")

	log.Info("email:", email)
	log.Info("pass:", pass)

	if email == "" || pass == "" {
		http.Error(w, "invalid parameter(s)", http.StatusBadRequest)
		return
	}

	profile := &proto_common.UserBasicInfo{
		Usn:      uint64(time.Now().Unix()),
		Uid:      1234,
		Birthday: 0,
		Gender:   proto_common.Gender_GENDER_FEMALE,
		Nickname: "wow",
		Email:    email,
		Avatar:   "http://www.gravatar.com/" + utils.GetStringMD5Hash(email),
		Country:  "cn",
	}

	js, err := json.Marshal(profile)

	log.Debug(js)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func startHTTP(addr string, rt, wt time.Duration) {
	router.HandleFunc("/login", loginHandler).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: wt,
		ReadTimeout:  rt,
	}
	log.Fatal(srv.ListenAndServe())
}