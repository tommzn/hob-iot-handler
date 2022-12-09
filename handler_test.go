package main

import (
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/suite"
	log "github.com/tommzn/go-log"
	timetracker "github.com/tommzn/hob-timetracker"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestProcessRequests() {

	handler := handlerForTest()
	event := clickEventForTest()

	suite.Nil(handler.Process(event))

	event.DeviceEvent.ButtonClicked.ReportedTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	suite.Nil(handler.Process(event))
}

func (suite *HandlerTestSuite) TestConvertClickType() {

	suite.Equal(timetracker.WORKDAY, toTimeTrackingRecordType(SINGLE_CLICK))
	suite.Equal(timetracker.ILLNESS, toTimeTrackingRecordType(DOUBLE_CLICK))
	suite.Equal(timetracker.VACATION, toTimeTrackingRecordType(LONG_PRESS))
}

func clickEventForTest() events.IoTOneClickEvent {
	return events.IoTOneClickEvent{
		DeviceEvent: events.IoTOneClickDeviceEvent{
			ButtonClicked: events.IoTOneClickButtonClicked{
				ClickType: string(SINGLE_CLICK),
			},
		},
		DeviceInfo: events.IoTOneClickDeviceInfo{DeviceID: "Device01"},
	}
}

func handlerForTest() *IOTOneClickRequestHandler {
	return newRequestHandler(timetracker.NewLocaLRepository(), loggerForTest())
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}
