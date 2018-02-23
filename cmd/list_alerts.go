package main

import (
	"github.com/kelledge/serviceTails/twilio"
	"fmt"
)

func main() {
	monitor := twilio.NewTwilioMonitorClient("ACabfe3927adeb3ccdde96c4ff19169620", "2c3d1dd51fa054bae8a6e6c64f0a4b94")
	list, _ := monitor.List("2018-02-16T20:43:17Z", "2018-02-16T21:01:27Z")
	fmt.Printf("%s", list)
}
