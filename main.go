package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/calendar"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	oauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	oauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	oauthConfig.RedirectURL = os.Getenv("REDIRECT_URL")
}

func main() {

	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Please open the following URL in your browser to authorize the application:")
	fmt.Println(authURL)

	http.HandleFunc("/auth/google/callback", handleOAuth2Callback)
	fmt.Println("Waiting for authorization...")
	http.ListenAndServe(":3000", nil)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	prayerTimes := GetPrayersTimingsApi()

	events := []struct {
		Summary   string
		StartTime string
		EndTime   string
	}{
		{"فجر", prayerTimes.Fajr, prayerTimes.Fajr},
		{"ظهر", prayerTimes.Dhuhr, prayerTimes.Dhuhr},
		{"عصر", prayerTimes.Asr, prayerTimes.Asr},
		{"مغرب", prayerTimes.Maghrib, prayerTimes.Maghrib},
		{"عشاء", prayerTimes.Isha, prayerTimes.Isha},
	}

	// Create events for each event detail
	for _, event := range events {
		err = createEvent(token, event.Summary, event.StartTime, event.EndTime)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Event created successfully!")
	fmt.Println("Events Created")
	os.Exit(1)

}

func GetPrayersTimingsApi() Timings {
	var city string = "cairo"
	var country string = "egypt"
	if len(os.Args) >= 3 {
		city = os.Args[1]
		country = os.Args[2]
	}

	res, err := http.Get("http://api.aladhan.com/v1/timingsByCity?city=" + city + "&country=" + country + "&method=5&iso8601=true")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("could not get the prayers api")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var PrayingTimes PrayerTimes
	err = json.Unmarshal(body, &PrayingTimes)
	if err != nil {
		panic(err)
	}

	prayersTime := PrayingTimes.Data.Timings
	return prayersTime

}

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
