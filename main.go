package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var port string = "8000"

type requestLogger struct{}

func (rl requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	APIurl := "http://localhost:8000"
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
	}

	fmt.Printf("Headers: %+v\n", r.Header)

	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
			fmt.Printf("JSON parse error: %v", err)
			return
		}
		var text string = string(prettyJSON.Bytes())
		var oldText string
		if text == oldText && text != " " {
			text = " "
			return
		}
		dataSet := []string{"trigger", "TRIGGER"}
		for _, data := range dataSet {
			if strings.Contains(text, data) == true {
				fmt.Println("TRIGGER DETECTED! Sending message to proxy")
				//POSTs in API URL
				byteSlice := []byte(text)
				f, err := os.Create("data.txt")
				_, err2 := f.WriteString(text)
				req, err := http.NewRequest("POST", APIurl, bytes.NewBuffer(byteSlice))
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				if err != nil {
					panic(err)
				}
				if err2 != nil {
					log.Fatal(err2)
				}
				defer resp.Body.Close()
				oldText = text
			}
		}
		fmt.Println(text)
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}

	// CORS headers
	//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	//w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Headers", "Accept-Encoding,Authorization,X-Forwarded-For,Content-Type,Origin,Server")
}

func main() {
	fmt.Printf("Starting request echo server on port %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), requestLogger{})
	fmt.Printf("Server error: %v\n", err)
}
