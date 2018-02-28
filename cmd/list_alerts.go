package main

import (
	"fmt"
	"os"

	"github.com/kelledge/serviceTails/twilio"
)

func main() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	monitor := twilio.NewTwilioMonitorClient(accountSid, authToken)

	poller := monitor.Poll()

	for result := range poller {
		for _, alert := range result {
			fmt.Printf("%v\n", alert)
		}
	}
}
