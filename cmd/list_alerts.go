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
	deduper := twilio.NewMonitorAlertDeduplicator()
	poller := monitor.Poll()

	for result := range poller {
		dedupedResult := deduper.Update(result)
		for _, alert := range dedupedResult {
			fmt.Printf("%v\n", alert)
		}
	}
}
