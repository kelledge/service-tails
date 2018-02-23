package twilio

import "net/http"
import "net/url"
import "encoding/json"

type MonitorAlert struct {
	Sid string `json:"sid"`
	AccountSid string `json:"account_sid"`
	ServiceSid string `json:"service_sid"`
	LogLevel   string `json:"log_level"`
	ErrorCode string `json:"error_code"`
	DateUpdate string `json:"date_updated"`
	DateGenerated string `json:"date_generated"`
	Url string `json:"url"`
}

type MonitorAlertListMeta struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

type MonitorAlertList struct {
	Alerts []*MonitorAlert `json:"alerts"`
	Meta *MonitorAlertListMeta `json:"meta"`
}

type MonitorClient struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
	HttpClient *http.Client
}

func NewTwilioMonitorClient(accountSid, authToken string) *MonitorClient {
	return &MonitorClient{
		AccountSid:accountSid,
		AuthToken:authToken,
		BaseUrl: "https://monitor.twilio.com/v1/Alerts/",
		HttpClient: &http.Client{},
	}
}

func (twilio *MonitorClient) List(startDate, endDate string) (*MonitorAlertList, error) {
	url := url.UR
	req, _ := http.NewRequest("GET", twilio.BaseUrl, nil)
	res, _ := twilio.HttpClient.Do(req)

	decoder := json.NewDecoder(res.Body)
	var list MonitorAlertList
	decoder.Decode(&list)

	return &list, nil
}