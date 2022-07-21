package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/go-toast/toast"
)

func sendNotification(buttonURL *url.URL, messageFormat string, messageArgs ...interface{}) {
	notification := toast.Notification{
		AppID:   "Global Entry Schedule Watcher",
		Title:   "Appointment Found!",
		Message: fmt.Sprintf(messageFormat, messageArgs...),
		Audio:   toast.LoopingAlarm,

		// I can't figure out how to get the URL to work...
		// it makes the notification not appear at all.  Oh well.
	}
	if err := notification.Push(); err != nil {
		log.Printf("Failed to push notification: %s", err)
	}
}
