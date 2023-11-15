package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/feishu-bot", feishuBot)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal(err)
	}
}

func feishuBot(writer http.ResponseWriter, request *http.Request) {

}
