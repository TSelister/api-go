package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type user struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

var database = make(map[string]user)

func main() {
	r := mux.NewRouter()
	r.Path("/user").Methods("POST").HandlerFunc(createUser)
	r.Path("/user/{id}").Methods("GET").HandlerFunc(getUser)
	r.Path("/user").Methods("PUT").HandlerFunc(putUser)
	r.Path("/user/{id}").Methods("DELETE").HandlerFunc(deleteUser)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10,
		ReadHeaderTimeout: 10,
		WriteTimeout:      10,
	}

	srv.ListenAndServe()
}

func validateUser(u *user) error {
	if u.Name == "" {
		return errors.New("o nome não pode ser vazio")
	}
	if u.Email == "" {
		return errors.New("o email não pode ser vazio")
	}
	if len(u.Password) < 6 {
		return errors.New("a senha deve possuir mais que 6 caracteres")
	}
	if len(u.Username) < 3 {
		return errors.New("o nome de usuário deve possuir mais que 3 caracteres")
	}
	return nil
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

	err = validateUser(&u)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	u.ID = uuid.NewString()

	database[u.ID] = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(u)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	q := mux.Vars(r)
	UserID := q["id"]

	user, ok := database[UserID]
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

	_, ok := database[u.ID]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("user not found"))
		return
	}

	if u.ID == "" {
		w.WriteHeader(400)
		w.Write([]byte("o id não pode ser vazio"))
		return
	}

	err = validateUser(&u)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	database[u.ID] = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(body)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	q := mux.Vars(r)
	UserID := q["id"]

	_, ok := database[UserID]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("user not found"))
		return
	}

	delete(database, UserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("user sucessfully deleted"))
}
