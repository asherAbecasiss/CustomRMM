package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	serverData := GetInfoForPc()

	// fileBytes, err := ioutil.ReadFile("info/info.txt")
	// if err != nil {
	// 	panic(err)
	// }

	url := "http://localhost:3000/api/postInfo"
	fmt.Println("URL:>", url)

	marshalled, err := json.Marshal(serverData)
	if err != nil {
		log.Fatalf("impossible to marshall teacher: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalled))
	if err != nil {
		log.Print(err)
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
