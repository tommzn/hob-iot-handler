package main

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

// NewRequestHandler create a handler to process API Gateway requests.
func newRequestHandler(timeTracker timetracker.TimeTracker, logger log.Logger) *IOTOneClickRequestHandler {
	return &IOTOneClickRequestHandler{
		logger:      logger,
		timeTracker: timeTracker,
	}
}

// Process will process time tracking request and persist it using time tracker repository.
func (handler *IOTOneClickRequestHandler) Process(event events.IoTOneClickEvent) error {

	defer handler.logger.Flush()

	handler.logger.Debugf("IoTOneClickDeviceEvent: %+v", event)

	timeTrackingRecord := toTimeTrackingRecord(event)
	handler.logger.Statusf("Receive capture request (%s) from %s at %s", timeTrackingRecord.ClickType, timeTrackingRecord.DeviceId, timeTrackingRecord.Timestamp)

	recordType := toTimeTrackingRecordType(timeTrackingRecord.ClickType)
	var err error
	if timeTrackingRecord.Timestamp == nil {
		err = handler.timeTracker.Capture(timeTrackingRecord.DeviceId, recordType)
	} else {
		err = handler.timeTracker.Captured(timeTrackingRecord.DeviceId, recordType, *timeTrackingRecord.Timestamp)
	}
	if err != nil {
		handler.logger.Error(err)
	}
	return err
}

// ToTimeTrackingRecordType converts a AWS IOT click type to a time tracking record type.
func toTimeTrackingRecordType(clickType IotClickType) timetracker.RecordType {
	switch clickType {
	case SINGLE_CLICK:
		return timetracker.WORKDAY
	case DOUBLE_CLICK:
		return timetracker.ILLNESS
	case LONG_PRESS:
		return timetracker.VACATION
	default:
		return timetracker.WORKDAY
	}
}

func toTimeTrackingRecord(event events.IoTOneClickEvent) TimeTrackingRecord {
	timeTrackingRecord := TimeTrackingRecord{
		DeviceId:  event.DeviceInfo.DeviceID,
		ClickType: IotClickType(event.DeviceEvent.ButtonClicked.ClickType),
	}
	if event.DeviceEvent.ButtonClicked.ReportedTime != "" {
		if timesstamp, err := time.Parse("2006-01-02T15:04:05.000Z", event.DeviceEvent.ButtonClicked.ReportedTime); err == nil {
			timeTrackingRecord.Timestamp = &timesstamp
		}
	}
	return timeTrackingRecord
}
