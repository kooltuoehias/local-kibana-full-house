package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	cloudWatchLog := cloudwatchlogs.New(session)
	// List log streams for a specific group (optional)
	logGroupName := "test-pnc-contract-service"
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Descending:   aws.Bool(true),
		LogGroupName: aws.String(logGroupName),
		Limit:        aws.Int64(5),
		OrderBy:      aws.String("LastEventTime"),
	}

	describeLogStreamsOutput, err := cloudWatchLog.DescribeLogStreams(input)
	if err != nil {
		fmt.Println("Error describing log streams:", err)
		return
	}

	if len(describeLogStreamsOutput.LogStreams) == 0 {
		log.Println("\t No log streams found for this group.")
	} else {
		log.Println("\t Log Streams:")
		for _, logStream := range describeLogStreamsOutput.LogStreams {
			log.Println("\t\t", *logStream.LogStreamName)
			getLog := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(logGroupName),
				LogStreamName: aws.String(*logStream.LogStreamName),
			}
			events, err := cloudWatchLog.GetLogEvents(getLog)
			if err != nil {
				panic("Error fetching log events: " + err.Error())
			}
			saveAsLogFile(*logStream.LogStreamName, events.Events)
		}
	}
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
	fmt.Println(splits)
	return splits[0] + "_" + splits[1] + "_" + splits[2] + ".log"
}

func saveAsFile(path string, logs string) {
	err := os.WriteFile(path, []byte(logs), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
