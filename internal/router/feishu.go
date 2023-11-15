package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/luomu/clean-code/pkg/logging/zerolog"
)

type AlertStatus string

const (
	AlertFiring        AlertStatus = "firing"
	AlertResolved      AlertStatus = "resolved"
	AlertNameLabel                 = "alertname"
	AlertSeverityLabel             = "severity"

	AlertSummaryAnnotation     = "summary"
	AlertDescriptionAnnotation = "description"

	FeishuWebhookAddr = "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx"
)

type AlertWebhookRequest struct {
	Data Data `json:"data"`
}

// Pair is a key/value string pair.
type Pair struct {
	Name, Value string
}

// Pairs is a list of key/value string pairs.
type Pairs []Pair

// Names returns a list of names of the pairs.
func (ps Pairs) Names() []string {
	ns := make([]string, 0, len(ps))
	for _, p := range ps {
		ns = append(ns, p.Name)
	}
	return ns
}

// Values returns a list of values of the pairs.
func (ps Pairs) Values() []string {
	vs := make([]string, 0, len(ps))
	for _, p := range ps {
		vs = append(vs, p.Value)
	}
	return vs
}

// KV is a set of key/value string pairs.
type KV map[string]string

// SortedPairs returns a sorted list of key/value pairs.
func (kv KV) SortedPairs() Pairs {
	var (
		pairs     = make([]Pair, 0, len(kv))
		keys      = make([]string, 0, len(kv))
		sortStart = 0
	)
	for k := range kv {
		if k == string(AlertNameLabel) {
			keys = append([]string{k}, keys...)
			sortStart = 1
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys[sortStart:])

	for _, k := range keys {
		pairs = append(pairs, Pair{k, kv[k]})
	}
	return pairs
}

// Remove returns a copy of the key/value set without the given keys.
func (kv KV) Remove(keys []string) KV {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	res := KV{}
	for k, v := range kv {
		if _, ok := keySet[k]; !ok {
			res[k] = v
		}
	}
	return res
}

// Names returns the names of the label names in the LabelSet.
func (kv KV) Names() []string {
	return kv.SortedPairs().Names()
}

// Values returns a list of the values in the LabelSet.
func (kv KV) Values() []string {
	return kv.SortedPairs().Values()
}

// Data is the data passed to notification templates and webhook pushes.
//
// End-users should not be exposed to Go's type system, as this will confuse them and prevent
// simple things like simple equality checks to fail. Map everything to float64/string.
type Data struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

// Alerts is a list of Alert objects.
type Alerts []Alert

// Firing returns the subset of alerts that are firing.
func (as Alerts) Firing() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(AlertFiring) {
			res = append(res, a)
		}
	}
	return res
}

// Resolved returns the subset of alerts that are resolved.
func (as Alerts) Resolved() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(AlertResolved) {
			res = append(res, a)
		}
	}
	return res
}

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type FeishuBotRequest struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

type CompanyMeta struct {
	companyName string
	serialNo    string
}

func registryWebhook(r *gin.RouterGroup) {
	r.POST("/feishu-bot", func(ctx *gin.Context) {
		var request AlertWebhookRequest
		ctx.ShouldBindJSON(&request.Data)
		meta := getCompanyMeta()

		marshal, _ := json.Marshal(request)
		log.Info("request: " + string(marshal))
		if len(request.Data.Alerts) > 0 {
			for _, alert := range request.Data.Alerts {
				message := getAlertMessage(alert, meta)
				feishuBotRequest := &FeishuBotRequest{
					MsgType: "text",
					Content: Content{
						Text: message,
					},
				}

				fbr, _ := json.Marshal(feishuBotRequest)
				log.Info("bot request: " + string(fbr))
				_, err := doRequest(FeishuWebhookAddr, "POST", fbr)
				if err != nil {
					log.Errorf("do request failed: %+v", err)
				}
			}
		}
	})
}

func getAlertMessage(alert Alert, meta CompanyMeta) string {
	company := fmt.Sprintf("%s | %s", meta.companyName, meta.serialNo)
	alertName := alert.Labels[AlertNameLabel]
	summary := alert.Annotations[AlertSummaryAnnotation]
	description := alert.Annotations[AlertDescriptionAnnotation]
	msg := fmt.Sprintf("[集群归属]：%s\n[AlertName]：%s\n[Status]：%s\n[Summary]：%s\n[Description]：%s\n", company, alertName, alert.Status, summary, description)
	return msg
}

func getCompanyMeta() CompanyMeta {
	companyName := os.Getenv("COMPANY_NAME")
	if len(companyName) == 0 {
		companyName = "Dev环境"
	}
	serialNo := os.Getenv("SERIAL_NO")
	if len(serialNo) == 0 {
		serialNo = ""
	}
	return CompanyMeta{
		companyName: companyName,
		serialNo:    serialNo,
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

	if res.StatusCode/100 == 2 {
		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)
		return buf, nil
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request failed with status code %v: %w", res.StatusCode, err)
	}
	return buf, err
}
