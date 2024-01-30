package main

type Timings struct {
	Fajr    string `json:"fajr"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

type Gregorian struct {
	Date string `json:"date"`
}

type PrayerTimes struct {
	Data struct {
		Timings Timings `json:"timings"`
		Date    Date    `json:"date"`
	} `json:"data"`
}

type Date struct {
	Gregorian Gregorian `json:"gregorian"`
}
