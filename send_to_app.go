package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io"
	"crypto/tls"
)

func sendToApp(data []byte) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := "https://spmonitor.soxprox.lan:3443/client/data"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println(err)
		}
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}