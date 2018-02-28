package twilio

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

const Error = "error"
const Warning = "warning"
const Notice = "notice"
const Debug = "debug"

type Interval struct {
	Start time.Time
	End   time.Time
}

type MonitorDate struct {
	time.Time
}

func (date *MonitorDate) UnmarshalJSON(data []byte) error {
	var iso8601 string
	json.Unmarshal(data, &iso8601)
	date.Time, _ = time.Parse(time.RFC3339, iso8601)
	return nil
}

type MonitorAlert struct {
	AlertText        string `json:"alert_text"`
	RequestURL       string `json:"request_url"`
	RequestMethod    string `json:"request_method"`
	RequestVariables string `json:"request_variables"`
	ResponseHeaders  string `json:"response_headers"`
	ResponseBody     string `json:"response_body"`
}

type MonitorTruncatedAlert struct {
	Sid           string      `json:"sid"`
	AccountSid    string      `json:"account_sid"`
	ServiceSid    string      `json:"service_sid"`
	LogLevel      string      `json:"log_level"`
	ErrorCode     string      `json:"error_code"`
	DateUpdate    MonitorDate `json:"date_updated"`
	DateGenerated MonitorDate `json:"date_generated"`
	URL           string      `json:"url"`
}

type MonitorAlertListMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type MonitorAlertList struct {
	Alerts []*MonitorTruncatedAlert `json:"alerts"`
	Meta   *MonitorAlertListMeta    `json:"meta"`
}

func (list *MonitorAlertList) Sort() *MonitorTruncatedAlert {
	sort.Slice(list.Alerts, func(i, j int) bool {
		return list.Alerts[i].DateGenerated.Unix() < list.Alerts[j].DateGenerated.Unix()
	})
	return nil
}

func (list *MonitorAlertList) Interval() *Interval {
	lastIndex := len(list.Alerts) - 1
	return &Interval{
		Start: list.Alerts[0].DateGenerated.Time,
		End:   list.Alerts[lastIndex].DateGenerated.Time,
	}
}

type MonitorClient struct {
	AccountSid string
	AuthToken  string
	BaseUrl    string
	HttpClient *http.Client
	// NOTE: Belongs on separate unit
	Ticker *time.Ticker
}

func NewTwilioMonitorClient(accountSid, authToken string) *MonitorClient {
	return &MonitorClient{
		AccountSid: accountSid,
		AuthToken:  authToken,
		BaseUrl:    "https://monitor.twilio.com/v1/Alerts/",
		HttpClient: &http.Client{},
	}
}

func (twilio *MonitorClient) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(twilio.AccountSid, twilio.AuthToken)
	return twilio.HttpClient.Do(req)
}

func (twilio *MonitorClient) List(startDate, endDate, logLevel string) (*MonitorAlertList, error) {
	req, _ := http.NewRequest("GET", twilio.BaseUrl, nil)

	q := req.URL.Query()
	q.Add("LogLevel", logLevel)
	q.Add("StartDate", startDate)
	q.Add("EndDate", endDate)

	res, _ := twilio.Do(req)

	decoder := json.NewDecoder(res.Body)
	var list MonitorAlertList
	decoder.Decode(&list)

	return &list, nil
}

// NOTE: Belongs on separate unit
func (twilio *MonitorClient) Poll() <-chan []*MonitorTruncatedAlert {
	twilio.Ticker = time.NewTicker(10 * time.Second)
	responses := make(chan []*MonitorTruncatedAlert)
	go func() {
		for _ = range twilio.Ticker.C {
			res, _ := twilio.List("2018-02-16T20:43:17Z", "2018-02-16T21:01:27Z", Error)
			responses <- res.Alerts
		}
	}()

	return responses
}
