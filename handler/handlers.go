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
			w.WriteHeader(400)
		}
		w.WriteHeader(201)
	}
}

func LogInHandlerFunc(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
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
		token, err := auth.LoginUser(userRequest, db)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
		}
		data, err := json.MarshalIndent(&token, "", " ")
		w.WriteHeader(200)
		w.Write(data)
	}
}

func ValidateHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Validate endpoint hit!")
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
		validationRequest := &auth.ValidateRequest{}
		err = json.Unmarshal(body, validationRequest)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if validationRequest == nil {
			w.WriteHeader(400)
			return
		}
		role, err := auth.ValidateToken(validationRequest.Token)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
		}
		res := &auth.ValidationResponse{Role: role}
		data, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	}
}

func ResetPasswordFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Edit user endpoint hit!")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
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
		passwordResetRequest := &auth.PasswordResetRequest{}
		err = json.Unmarshal(body, passwordResetRequest)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		if passwordResetRequest == nil {
			w.WriteHeader(400)
			return
		}
		resetId, err := auth.ResetPasswordRequest(passwordResetRequest, db, sender)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(401)
		}
		response := &auth.PasswordResetResponse{ResetId: resetId}
		data, err := json.MarshalIndent(response, "", " ")
		w.WriteHeader(200)
		w.Write(data)
	}
}

func PerformPasswordResetFunc(db *sql.DB, sender sender.ISender) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("PerformPasswordReset endpoint hit!")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
		}
		keys, ok := r.URL.Query()["resetId"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'resetId is missing'")
			w.WriteHeader(400)
			return
		}
		resetId := keys[0]
		pwResetRequest := &auth.PasswordReset{}
		err = json.Unmarshal(body, pwResetRequest)
		if err != nil {
			w.WriteHeader(500)
		}
		log.Println("Using rest id: " + resetId)
		pwResetRequest.ResetId = resetId
		if err := auth.PerformPasswordReset(pwResetRequest, db, sender); err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
			return
		}
		w.WriteHeader(200)
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
