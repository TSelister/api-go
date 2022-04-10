package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

var database = make(map[string]user)

func main() {
	r := mux.NewRouter()
	r.Path("/user").Methods("POST").HandlerFunc(createUser)
	r.Path("/user/{email}").Methods("GET").HandlerFunc(getUser)
	r.Path("/user").Methods("PUT").HandlerFunc(putUser)
	r.Path("/user/{email}").Methods("DELETE").HandlerFunc(deleteUser)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10,
		ReadHeaderTimeout: 10,
		WriteTimeout:      10,
	}

	srv.ListenAndServe()
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u user

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	database[u.Email] = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(body)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	q := mux.Vars(r)
	email := q["email"]

	user, ok := database[email]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("user not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(user)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	var u user

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	_, ok := database[u.Email]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("user not found"))
		return
	}

	database[u.Email] = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(body)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	q := mux.Vars(r)
	email := q["email"]

	_, ok := database[email]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("user not found"))
		return
	}

	delete(database, email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("user sucessfully deleted"))
}
