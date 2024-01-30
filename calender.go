package main

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

//func addEvent(w http.ResponseWriter, r *http.Request) {

// code := r.URL.Query().Get("code")
// token, err := oauthConfig.Exchange(context.Background(), code)
// if err != nil {
// 	http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
// 	return
// }

// 	events := []struct {
// 		Summary   string
// 		StartTime time.Time
// 		EndTime   time.Time
// 	}{
// 		{"Event 1", time.Now().Add(24 * time.Hour), time.Now().Add(24*time.Hour + time.Hour)},
// 		{"Event 2", time.Now().Add(48 * time.Hour), time.Now().Add(48*time.Hour + time.Hour)},
// 	}

// 	// Create events for each event detail
// 	for _, event := range events {
// 		err = createEvent(token, event.Summary, event.StartTime, event.EndTime)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	fmt.Fprintf(w, "Event created successfully!")
// }

func createEvent(token *oauth2.Token, summary string, startTime, endTime time.Time) error {
	client := oauthConfig.Client(context.Background(), token)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: "prayer timings",
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
		},
	}

	_, err = srv.Events.Insert("primary", event).Do()
	if err != nil {
		return err
	}
	return nil
}

// func(*calendar.EventDateTime).MarshalJSON(){

// }
