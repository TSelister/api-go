package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", switchHTTPMethod)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10,
		ReadHeaderTimeout: 10,
		WriteTimeout:      10,
	}

	srv.ListenAndServe()
}

func switchHTTPMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHello(w, r)
	case "POST":
		postHello(w, r)
	case "PUT":
		putHello(w, r)
	case "DELETE":
		deleteHello(w, r)
	default:
		return
	}
}

func getHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello Get"))
}

func postHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello post"))
}

func putHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello put"))
}

func deleteHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello delete"))
}
