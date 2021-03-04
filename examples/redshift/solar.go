package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-logr/logr"
)

type SunriseSunset struct {
	Status  string     `json:"status"`
	Results SolarTimes `json:"results"`
}
type SolarTimes struct {
	Sunrise               string `json:"sunrise"`
	Sunset                string `json:"sunset"`
	SolarNoon             string `json:"solar_noon"`
	DayLength             int    `json:"day_length"`
	CivilTwilightStart    string `json:"civil_twilight_begin"`
	CivilTwilightEnd      string `json:"civil_twilight_end"`
	NauticalTwilightStart string `json:"nautical_twilight_begin"`
	NauticalTwilightEnd   string `json:"nautical_twilight_end"`
	AstroTwilightStart    string `json:"astronomical_twilight_begin"`
	AstroTwilightEnd      string `json:"astronomical_twilight_end"`
}

func getSolarTimes(log logr.Logger) (sunrise, noon, sunset time.Duration, err error) {
	zeroD := mustParseDuration("0h")

	client := http.Client{}

	query := url.Values{}
	query.Set("lat", "52")
	query.Set("lng", "0.0000001") // Returns 400 INVALID_REQUEST if lng is any rendering of 0.0
	query.Set("formatted", "0")
	url := url.URL{
		Scheme:   "https",
		Host:     "api.sunrise-sunset.org",
		Path:     "json",
		RawQuery: query.Encode(),
	}
	log.V(1).Info("Querying sunrise-sunset.org", "url", url.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		log.Error(err, "Can't make http request?")
		return zeroD, zeroD, zeroD, err
	}

	req.Header.Set("user-agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		log.Error(err, "Can't get solar times info")
		return zeroD, zeroD, zeroD, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err, "Can't get solar times info")
		return zeroD, zeroD, zeroD, err
	}

	times := SunriseSunset{}
	err = json.Unmarshal(body, &times)
	if err != nil {
		// TODO: if eg request is invalid, we get status:INVALID_REQUEST,results:"". This breaks Unmarshal, hence status currently isn't seen
		log.Error(err, "Can't get solar times info")
		return zeroD, zeroD, zeroD, err
	}

	if times.Status != "OK" {
		log.Error(err, "Can't get solar times info", "status", times.Status)
		return zeroD, zeroD, zeroD, err
	}

	sunrise = extractDuration(times.Results.Sunrise)
	noon = extractDuration(times.Results.SolarNoon)
	sunset = extractDuration(times.Results.Sunset)
	err = nil

	log.V(2).Info("SolarTimes", "sunrise", sunrise, "noon", noon, "sunset", sunset)

	return
}

func extractDuration(str string) time.Duration {
	t := mustParseTime(time.RFC3339, str)
	return sinceMidnight(t)
}
