package main

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func createEvent(token *oauth2.Token, summary string, startTime, endTime string) error {
	client := oauthConfig.Client(context.Background(), token)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: "prayer timing",
		Start: &calendar.EventDateTime{
			DateTime: startTime,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
		},
	}

	_, err = srv.Events.Insert("primary", event).Do()
	if err != nil {
		return err
	}
	return nil
}
