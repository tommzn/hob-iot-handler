package main

import (
	"github.com/aws/aws-lambda-go/events"
)

// Handler is used to process published from AWS IOT 1-Click.
type Handler interface {

	// Process will handle button click events
	Process(events.IoTOneClickEvent) error
}
