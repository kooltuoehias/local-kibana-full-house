package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func Do(env string, service string, buffer *bytes.Buffer) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		buffer.WriteString(fmt.Sprintln("Error creating session:", err))
		return
	}
	cloudWatchLog := cloudwatchlogs.New(session)
	logGroupName := env + "-" + service
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Descending:   aws.Bool(true),
		LogGroupName: aws.String(logGroupName),
		Limit:        aws.Int64(5),
		OrderBy:      aws.String("LastEventTime"),
	}
	buffer.WriteString(fmt.Sprintln("Log Streams:", input))

	describeLogStreamsOutput, err := cloudWatchLog.DescribeLogStreams(input)
	if err != nil {
		buffer.WriteString(fmt.Sprintln("Error describing log streams:", err))
		return
	}

	for _, logStream := range describeLogStreamsOutput.LogStreams {
		saveLogStream(cloudWatchLog, logGroupName, *logStream.LogStreamName)
	}

}

func saveLogStream(cloudWatchLog *cloudwatchlogs.CloudWatchLogs, logGroupName string, logStreanName string) {
	log.Println("\tShall Save ", logStreanName)

	getLog := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreanName),
	}
	events, err := cloudWatchLog.GetLogEvents(getLog)
	if err != nil {
		panic("Error fetching log events: " + err.Error())
	}
	saveAsLogFile(logStreanName, events.Events)
}

func saveAsLogFile(logGroupName string, events []*cloudwatchlogs.OutputLogEvent) {
	path := "logs/" + createLogFileNameFromLogGroupName(logGroupName)
	content := ""
	for _, event := range events {
		content += *event.Message + "\n"
	}
	saveAsFile(path, content)
}

func createLogFileNameFromLogGroupName(logGroupName string) string {
	splits := strings.Split(logGroupName, "/")
	return splits[0] + "_" + splits[1] + "_" + splits[2] + ".log"
}

func saveAsFile(path string, logs string) {
	err := os.WriteFile(path, []byte(logs), 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("\tJust Saved ", path)
	}
}
