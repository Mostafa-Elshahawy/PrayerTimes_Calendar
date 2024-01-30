package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	oauthConfig = &oauth2.Config{
		Scopes: []string{"https://www.googleapis.com/auth/calendar"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
)

// loading .env vars that contains the Oauth client credentials
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
	// prompt the user with the redirection link
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Please open the following URL in your browser to authorize the application:")
	fmt.Println(authURL)

	http.HandleFunc("/auth/google/callback", handleOAuth2Callback)
	fmt.Println("Waiting for authorization...")
	http.ListenAndServe(":3000", nil)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// get the user token
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// get prayer timings from the api endpoint
	prayerTimes := GetPrayersTimingsApi()
	for _, timing := range prayerTimes {

		events := []struct {
			Summary   string
			StartTime string
			EndTime   string
		}{
			{"فجر", timing.Fajr, timing.Fajr},
			{"ظهر", timing.Dhuhr, timing.Dhuhr},
			{"عصر", timing.Asr, timing.Asr},
			{"مغرب", timing.Maghrib, timing.Maghrib},
			{"عشاء", timing.Isha, timing.Isha},
		}

		// Create events for each event detail
		for _, event := range events {
			err = createEvent(token, event.Summary, event.StartTime, event.EndTime)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	// print the success message and exit
	fmt.Println("Events Created")
	os.Exit(1)

}

func GetPrayersTimingsApi() []Timings {
	// getting the required data for the api (year, month, city, country)
	var year = fmt.Sprint(time.Now().Year())
	var month = fmt.Sprint(int(time.Now().Month()))
	var city string = "cairo"
	var country string = "egypt"

	// taking args from the user for their location for specific timings
	if len(os.Args) >= 3 {
		city = os.Args[1]
		country = os.Args[2]
	}

	res, err := http.Get("http://api.aladhan.com/v1/calendarByCity/" + year + "/" + month + "?city=" + city + "&country=" + country + "&method=5&iso8601=true")
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

	//converting the json response from the api to go structs to deal with
	var PrayingTimes PrayerTimes
	err = json.Unmarshal(body, &PrayingTimes)
	if err != nil {
		panic(err)
	}

	// get timings only from the json response
	timings := PrayingTimes.GetTimings()

	return timings

}

// createEvent takes the user token and the event data to make the actual event in the calendar
func createEvent(token *oauth2.Token, summary string, startTime, endTime string) error {
	client := oauthConfig.Client(context.Background(), token)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: "Prayer Timing",
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

// parsing the json data from the response to get the required form for the calendar
func (pt *PrayerTimes) GetTimings() []Timings {
	timings := make([]Timings, len(pt.Data))
	for i, data := range pt.Data {
		timings[i] = data.Timings
	}
	return timings
}
