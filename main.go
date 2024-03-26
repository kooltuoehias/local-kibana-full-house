package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	session, err := session.NewSession()
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	cloudWatchLog := cloudwatchlogs.New(session)

	// List log groups
	fmt.Println("Listing Log Groups:")
	params := &cloudwatchlogs.DescribeLogGroupsInput{}

	describeLogGroupsOutput, err := cloudWatchLog.DescribeLogGroups(params)
	if err != nil {
		fmt.Println("Error describing log groups:", err)
		return
	}

	for _, group := range describeLogGroupsOutput.LogGroups {
		fmt.Println("\t", *group.LogGroupName)
	}

	// List log streams for a specific group (optional)
	fmt.Println("\nListing Log Streams (for a specific group):")
	logGroupName := "your-log-group-name" // Replace with your actual group name (optional)
	params.SetLogGroupNamePattern(logGroupName)
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	}

	describeLogStreamsOutput, err := cloudWatchLog.DescribeLogStreams(input)
	if err != nil {
		fmt.Println("Error describing log streams:", err)
		return
	}

	if len(describeLogStreamsOutput.LogStreams) == 0 {
		fmt.Println("\t No log streams found for this group.")
	} else {
		fmt.Println("\t Log Streams:")
		for _, stream := range describeLogStreamsOutput.LogStreams {
			fmt.Println("\t\t", *stream.LogStreamName)
		}
	}
}
