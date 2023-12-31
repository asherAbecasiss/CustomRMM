package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
)

var serverName = ""

var (
	host       = "smtp.gmail.com"
	username   = ""
	password   = ""
	portNumber = "587"
)

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func New() *Sender {
	auth := smtp.PlainAuth("", username, password, host)
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.auth, username, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func SendMail(MessageTitle string, Message string) {

	sender := New()
	m := NewMessage(serverName+" Server ERROR! "+MessageTitle, Message)
	m.To = []string{"", ""}
	m.CC = []string{""}
	m.BCC = []string{""}
	// m.AttachFile("info.txt")
	fmt.Println(sender.Send(m))

}

func MailType(typeOfMessage int, se ServerInfo) {

	switch typeOfMessage {
	case 1:
		log.Print("Swarm Node ERROR")
		SendMail("Swarm Node ERROR. Server IP : "+se.ServerIp, fmt.Sprint(se.SwarmNode))
	case 2:
		log.Print("Swarm Services ERROR")
		// services := fmt.Sprint(se.DockerServices)
		var temp string
		for _, v := range se.DockerServices {
			//temp := v.ServiceStatus.DesiredTasks
			needToBy := fmt.Sprint(v.ServiceStatus.RunningTasks)
			repl := fmt.Sprint(v.ServiceStatus.DesiredTasks)
			// fmt.Println(nn)
			// fmt.Println(b)
			i, _ := strconv.Atoi(needToBy)
			j, _ := strconv.Atoi(repl)
			var point string = ""
			if i != j {
				point = "--> "
			} else {
				point = ""
			}

			temp += point + v.Spec.Name + " " + needToBy + "/" + repl + "\n"

		}
		SendMail("Swarm Services ERROR. Server IP : "+se.ServerIp, temp)
	case 3:
		log.Print("Multiple files in queue-main.")
		SendMail("Multiple files in queue-main. Server IP : "+se.ServerIp, "if number = -100  no such file or directory\n"+fmt.Sprint(se.CountFils))

	case 4:
		log.Print("Disk space less then 50 Gb. ")
		mm := fmt.Sprint(se.FreeDiskSpace)
		fmt.Println(mm)
		SendMail("Disk space less then 50 Gb. Server IP : "+se.ServerIp, mm)

	case 5:
		log.Print("Mem Percent high")
		m := fmt.Sprint(se.MemoryPercent)

		SendMail("Memory Percent High. Server IP : "+se.ServerIp, m)
	}

}
