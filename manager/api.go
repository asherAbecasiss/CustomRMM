package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/swarm"
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
	router.HandleFunc("/api/postInfo", s.postInfo).Methods("POST")

	log.Printf("server Runing on port %s", s.portNumber)

	http.ListenAndServe(s.portNumber, router)
}

type DirFiles struct {
	Path  string
	Count int
}
type ServerInfo struct {
	Date           string          `json:"date"`
	ServerIp       string          `json:"serverIp"`
	FreeDiskSpace  uint64          `json:"freediskspace"`
	MemoryPercent  int             `json:"memorypercent"`
	SwarmNode      []swarm.Node    `json:"swarmnode"`
	DockerServices []swarm.Service `json:"dockerservices"`
	CountFils      []DirFiles      `json:"countfils"`
}

func (s *ApiServer) postInfo(w http.ResponseWriter, r *http.Request) {

	se := ServerInfo{}
	err := json.NewDecoder(r.Body).Decode(&se)
	if err != nil {
		panic(err)
	}
	var toSendMail int = 0
	// log.Print(se.CountFils)
	toSendMail = CheckForNode(se.SwarmNode)
	MailType(toSendMail, se)
	toSendMail = 0

	toSendMail = CheckForServices(se.DockerServices)
	MailType(toSendMail, se)
	toSendMail = 0

	toSendMail = CheckForNumberOfFilesInDir(se.CountFils)
	MailType(toSendMail, se)
	toSendMail = 0

	toSendMail = CheckForDisk(se.FreeDiskSpace)
	MailType(toSendMail, se)
	toSendMail = 0

	toSendMail = CheckForMemPercent(se.MemoryPercent)
	MailType(toSendMail, se)
	toSendMail = 0
}
