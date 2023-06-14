package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/carlmjohnson/requests"
)

func getRequest(ctx context.Context) {
	var s string
	err := requests.
		URL("http://example.com").
		ToString(&s).
		Fetch(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func postRequest(ctx context.Context) {
	err := requests.
		URL("https://postman-echo.com/post").
		BodyBytes([]byte(`hello, world`)).
		ContentType("text/plain").
		ToWriter(os.Stdout).
		Fetch(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewGetRequest creates a new GET request to server
func NewGetRequest(addr, path string) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", addr, path)

	client := http.Client{}
	res, err := client.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, data, nil
}

// NewPostRequest creates a new POST request to server
func NewPostRequest(addr, path, contentType string, data interface{}) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", addr, path)

	buf, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}

	client := http.Client{}
	res, err := client.Post(url, contentType, bytes.NewReader(buf))
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, body, nil
}
