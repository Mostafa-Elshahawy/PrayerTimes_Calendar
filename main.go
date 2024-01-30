package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

	//accessToken := fmt.Sprint("Access token:", token.AccessToken)
	//refreshToken := fmt.Sprint("Refresh token:", token.RefreshToken)

	events := []struct {
		Summary   string
		StartTime time.Time
		EndTime   time.Time
	}{
		{"Event 1", time.Now().Add(24 * time.Hour), time.Now().Add(24*time.Hour + time.Hour)},
		{"Event 2", time.Now().Add(48 * time.Hour), time.Now().Add(48*time.Hour + time.Hour)},
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

	fmt.Fprintf(w, "Authorization successful! You can close this window.")
}

func GetPrayersTimingsApi() {
	res, err := http.Get("http://api.aladhan.com/v1/calendarByAddress/2024/1?address=Cairo&method=5")
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

	prayers, date := PrayingTimes.Data[1].Timings, PrayingTimes.Data[1].Date.Gregorian.Date

	//Get the reflect.Value of the struct
	val := reflect.ValueOf(prayers)

	//Loop over each field
	for i := 0; i < val.NumField(); i++ {
		// Get the field name
		fieldName := val.Type().Field(i).Name

		// Get the field value
		fieldValue := val.Field(i).Interface()

		// Print the field name and value
		fmt.Printf("%s: %v\n", fieldName, fieldValue)

	}

	fmt.Printf("date: %s\n",
		date,
	)
}
