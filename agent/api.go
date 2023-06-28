package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	portNumber string
}

func NewApiServer(portNumber string) *ApiServer {
	return &ApiServer{
		portNumber: portNumber,
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Contanet-type", "application/json")
	json.NewEncoder(w).Encode(v)

}

func (s *ApiServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/api/GetInfo", s.GetInfo).Methods("GET")

	log.Printf("server Runing on port %s", s.portNumber)

	http.ListenAndServe(s.portNumber, router)
}

func (s *ApiServer) GetInfo(w http.ResponseWriter, r *http.Request) {

	GetInfoForPc()

	fileBytes, err := ioutil.ReadFile("info/info.txt")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)

}
