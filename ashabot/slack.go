package ashabot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func getSlackTokens() (slackTokens, error) {
	appToken, found := os.LookupEnv("SLACK_APP_TOKEN")
	if !found {
		return slackTokens{}, fmt.Errorf("SLACK_APP_TOKEN not found in environment")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return slackTokens{}, fmt.Errorf("SLACK_APP_TOKEN is not a valid xapp token")
	}

	botToken, found := os.LookupEnv("SLACK_BOT_TOKEN")
	if !found {
		return slackTokens{}, fmt.Errorf("SLACK_BOT_TOKEN not found in environment")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return slackTokens{}, fmt.Errorf("SLACK_BOT_TOKEN is not a valid xoxb token")
	}

	return slackTokens{appToken: appToken, botToken: botToken}, nil
}

func handleSlackEvents(client *socketmode.Client) {
	for event := range client.Events {
		switch event.Type {
		case socketmode.EventTypeConnecting:
			log.Println("Connecting to Slack...")

		case socketmode.EventTypeConnectionError:
			log.Println("Connection failed. Retrying later...")

		case socketmode.EventTypeConnected:
			log.Println("Connected to Slack")

		case socketmode.EventTypeEventsAPI:
			eventsApiEvent, ok := event.Data.(slackevents.EventsAPIEvent)
			if !ok {
				log.Printf("Ignored %+v\n", event)
				continue
			}
			client.Debugf("Event recieved: %+v\n", eventsApiEvent)
			client.Ack(*event.Request)
			handleEventsApi(client, eventsApiEvent)

		case socketmode.EventTypeInteractive:
		case socketmode.EventTypeSlashCommand:
			client.Debugf("Received slash command: %s", event.Data)
			client.Ack(*event.Request)
			continue

		default:
			log.Printf("Unexpected event type received: %s\n", event.Type)
		}
	}
}

func handleEventsApi(client *socketmode.Client, event slackevents.EventsAPIEvent) {
	switch event.Type {
	case slackevents.MemberJoinedChannel:
		log.Printf("Member joined channel: %s\n", event.Data)
	default:
		client.Debugf("Unsupported Events API event recieved: %s\n", event.Type)
	}
}
