package main

import (
	"fmt"
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
	clientTest()
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
