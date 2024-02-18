package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestGetPrayersTimingsApi(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"data":[{"timings":{"Fajr":"05:00","Dhuhr":"12:00","Asr":"15:00","Maghrib":"18:00","Isha":"20:00"}}]}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	apiURL := "http://" + ts.Listener.Addr().String()
	os.Args = []string{"cmd", "city", "country", apiURL}

	timings := GetPrayersTimingsApi()
	// assuming the month is 30 days
	if len(timings) != 30 {
		t.Errorf("Expected many timings entry, got %d", len(timings))
	}
}

func TestHandleOAuth2Callback(t *testing.T) {
	req := httptest.NewRequest("GET", "/auth/google/callback?code=test", nil)
	w := httptest.NewRecorder()

	handleOAuth2Callback(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected internal server error, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	expectedErrorMsg := "Failed to exchange token"
	if !strings.Contains(body, expectedErrorMsg) {
		t.Errorf("Expected error message '%s' in response, got '%s'", expectedErrorMsg, body)
	}
}

func TestCreateEvent(t *testing.T) {

	token := &oauth2.Token{
		AccessToken: "mockAccessToken",
		TokenType:   "Bearer",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		expectedAuthHeader := "Bearer mockAccessToken"
		if authHeader != expectedAuthHeader {
			t.Errorf("Expected Authorization header: %s, got: %s", expectedAuthHeader, authHeader)
		}

		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	oauthConfig.Endpoint = oauth2.Endpoint{TokenURL: ts.URL}
	client := oauthConfig.Client(context.Background(), token)
	service, _ := calendar.NewService(context.Background(), option.WithHTTPClient(client))

	err := createEvent(token, "Test Event", "2024-02-18T12:00:00Z", "2024-02-18T13:00:00Z")

	if err != nil {
		t.Errorf("Expected no error during event creation, got %v", err)
	}

	events, _ := service.Events.List("primary").Do()
	if len(events.Items) == 0 {
		t.Error("Expected events to be created")
	}
}
