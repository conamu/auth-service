package handler

import (
	"auth-service/auth"
	"auth-service/sender"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func SignUpHandlerFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Signup endpoint hit!")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if r.ContentLength == 0 {
			w.WriteHeader(400)
			return
		}
		userRequest := &auth.UserRequest{}
		err = json.Unmarshal(body, userRequest)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if userRequest == nil {
			w.WriteHeader(400)
			return
		}
		log.Println(userRequest)

		err = auth.RegisterUser(userRequest, db, sender)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
		}
		w.WriteHeader(201)
	}
}

func LogInHandlerFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Login endpoint hit!")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if r.ContentLength == 0 {
			w.WriteHeader(400)
			return
		}
		userRequest := &auth.UserRequest{}
		err = json.Unmarshal(body, userRequest)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if userRequest == nil {
			w.WriteHeader(400)
			return
		}
		log.Println(userRequest)
		err = auth.LoginUser(userRequest, db, sender)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
		}

		w.WriteHeader(200)
	}
}

func EditUserHandlerFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Edit user endpoint hit!")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if r.ContentLength == 0 {
			w.WriteHeader(400)
			return
		}
		userRegister := &auth.UserRequest{}
		err = json.Unmarshal(body, userRegister)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if userRegister == nil {
			w.WriteHeader(400)
			return
		}
		log.Println(userRegister)

		w.WriteHeader(202)
	}
}

func ResetPasswordFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Edit user endpoint hit!")
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}

		w.WriteHeader(202)
	}
}

func PingHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Ping endpoint hit!")
		if r.Method != "GET" {
			w.WriteHeader(405)
			return
		}
		w.WriteHeader(200)
	}
}
