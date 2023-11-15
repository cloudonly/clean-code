package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/configor"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	internalhttp "github.com/luomu/clean-code/pkg/http"
)

type Config struct {
	APPName string `default:"app name"`
	DB      struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DB_PASSWORD"`
		Port     uint   `default:"3306"`
	}
	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}

func main() {
	//requestMain()
	//configMain()
	//fmt.Println(util.Hash("2734"))
	//getAllResource()
	//clientTest()
	test()
}

type FeishuBotRequest struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

func test() {
	msg := fmt.Sprintf("%s | %s | %s | %s", "firing", "KubePodNotReady", "Pod has been in a non-ready state for more than 15 minutes.", "Pod loki-monitoring/event-exporter-bb4557cc5-zsrb7 has been in a non-ready state for longer than 15 minutes.\\n")
	feishuBotRequest := &FeishuBotRequest{
		MsgType: "text",
		Content: Content{
			Text: msg,
		},
	}

	fbr, _ := json.Marshal(feishuBotRequest)
	log.Println("bot request: " + string(fbr))
	addr := "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx"
	_, err := doRequest(addr, "POST", fbr)
	if err != nil {
		log.Println("do request failed: ", err)
	}
}

func doRequest(apiEndpoint, method string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, apiEndpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute HTTP request
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println("#####1")

	if res.StatusCode/100 == 2 {
		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)
		fmt.Println("#####2")
		fmt.Println(string(buf))
		return buf, nil
	}

	fmt.Println("#####3")
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request failed with status code %v: %w", res.StatusCode, err)
	}

	return buf, nil
}

func requestMain() {
	code, r, err := internalhttp.NewGetRequest("http://example.com", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("code: %v\n", code)
	fmt.Printf("result: %s\n", r)
}

func configMain() {
	var config Config
	configor.Load(&config, "config.yaml")
	fmt.Printf("Config: %v\n", config)
}

func group() {
	var g errgroup.Group
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("OK"))
	})
	g.Go(func() error {
		return http.ListenAndServe(":8080", mux)
	})
	g.Wait()
}

func clientTest() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})
	serveFunc := func() {
		fmt.Println("Just do it...")
		err := http.ListenAndServe(net.JoinHostPort("localhost", strconv.Itoa(8080)), mux)
		if err != nil {
			klog.ErrorS(err, "Failed to start healthz server")
		}
	}
	wait.Until(serveFunc, 5*time.Second, wait.NeverStop)
	select {}
}
