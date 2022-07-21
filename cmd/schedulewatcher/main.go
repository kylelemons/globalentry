// schedulewatcher watches for openings for Global Entry (and probably other TTP) schedules.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	location = flag.Int("location", -1, "Location codes (https://github.com/Drewster727/goes-notify#goes-center-codes).")
	remote   = flag.Bool("remote", false, "If true, also watch for remote openings.")
	people   = flag.Int("people", 1, "Number of people to schedule for")
	every    = flag.Duration("every", 5*time.Minute, "Frequency to poll the API")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *location == 0 && !*remote {
		log.Fatalf("Specify --remote and/or a --location from e.g. https://github.com/Drewster727/goes-notify#goes-center-codes")
	}

	if *location >= 0 {
		go watchLocation(ctx, *every, locationURL(*location, *people))
	}
	if *remote {
		go watchLocation(ctx, *every, remoteURL(*people))
	}

	select {}
}

func locationURL(location, people int) *url.URL {
	// e.g. https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=1&locationId=5446&minimum=1
	return &url.URL{
		Scheme: "https",
		Host:   "ttp.cbp.dhs.gov",
		Path:   "/schedulerapi/slots",
		RawQuery: url.Values{
			"orderBy":    {"soonest"},
			"limit":      {"1"},
			"locationId": {fmt.Sprint(location)},
			"minimum":    {fmt.Sprint(people)},
		}.Encode(),
	}
}

func remoteURL(people int) *url.URL {
	// e.g. https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=1&remote=true&minimum=1
	// Example (successful) response:
	// [
	//   {
	//     "locationId": 16496,
	//     "startTimestamp": "2022-08-23T17:30",
	//     "endTimestamp": "2022-08-23T17:45",
	//     "active": true,
	//     "duration": 15,
	//     "remoteInd": true
	//   }
	// ]
	return &url.URL{
		Scheme: "https",
		Host:   "ttp.cbp.dhs.gov",
		Path:   "/schedulerapi/slots",
		RawQuery: url.Values{
			"orderBy": {"soonest"},
			"limit":   {"1"},
			"remote":  {"true"},
			"minimum": {fmt.Sprint(people)},
		}.Encode(),
	}
}

func watchLocation(ctx context.Context, every time.Duration, slotsURL *url.URL) {
	tick := time.NewTicker(every)
	defer tick.Stop()

	scrape(ctx, slotsURL)
	for range tick.C {
		scrape(ctx, slotsURL)
	}
}

func scrape(ctx context.Context, slotsURL *url.URL) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    slotsURL,
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Printf("Failed to scrape: %s", err)
		return
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Raading body from %q: %s", slotsURL, err)
		return
	}

	var appointments []struct {
		LocationID     int    `json:"locationId"`
		StartTimestamp string `json:"startTimestamp"`
		EndTimestamp   string `json:"endTimestamp"`
		Active         bool   `json:"active"`
		Duration       int    `json:"duration"`
		Remote         bool   `json:"remoteInd"`
	}
	if err := json.Unmarshal(content, &appointments); err != nil {
		log.Printf("Decoding response body (%q): %s", content, err)
		return
	}

	if len(appointments) == 0 {
		log.Printf("No appointment at %s", slotsURL)
		return
	}

	log.Printf("Appointment found!")
	log.Printf("  %s", slotsURL)
	for _, appt := range appointments {
		if appt.Remote {
			log.Printf("  - Remote appointment at %s)!", appt.StartTimestamp)
			sendNotification(slotsURL, "Remote appointment found!")
		} else {
			log.Printf("  - Appointment found at %v at %s!", appt.LocationID, appt.StartTimestamp)
			sendNotification(slotsURL, "Onsite appointment (%v) found!", appt.LocationID)
		}
	}

}
