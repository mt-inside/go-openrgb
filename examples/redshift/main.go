/*
* Heavily inspired by: https://github.com/lucasb-eyer/go-colorful/blob/master/doc/gradientgen/gradientgen.go
 */

package main

import (
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mt-inside/go-openrgb/pkg/logging"
	"github.com/mt-inside/go-openrgb/pkg/model"
)

const userAgent = "go-openrgb-redshift"

func mustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

func mustParseTime(s, fmt string) time.Time {
	t, err := time.Parse(s, fmt)
	if err != nil {
		panic("MustParseTime: " + err.Error())
	}
	return t
}
func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("MustParseDuration: " + err.Error())
	}
	return d
}

func sinceMidnight(t time.Time) time.Duration {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return t.Sub(midnight)
}

func main() {
	log := logging.GetLogger(false)

	m, err := model.NewModel(log, "localhost:6742", userAgent)
	if err != nil {
		log.Error(err, "Couldn't synchronise devices and colors from server")
		os.Exit(1)
	}

	// FIXME update as we cross midnight local time
	sunrise, noon, sunset, err := getSolarTimes(log)
	if err != nil {
		log.Error(err, "Couldn't get Solar times")
	}

	wake, sleep := getSchedule()

	keyPoints := gradientTable{
		{mustParseHex("#000000"), mustParseDuration("0h")},
		{mustParseHex("#000000"), wake},
		{mustParseHex("#00ffff"), sunrise},
		{mustParseHex("#ffff00"), noon},
		{mustParseHex("#0000ff"), sunset},
		{mustParseHex("#ff0000"), sleep},
		{mustParseHex("#000000"), mustParseDuration("24h")},
	}

	for {
		d := sinceMidnight(time.Now())
		c := getColor(keyPoints, d)
		m.SetColor(c)

		err = m.Thither()
		if err != nil {
			log.Error(err, "Couldn't synchronise colors to server")
			os.Exit(1)
		}

		log.Info("Updated color", "color", c.Hex())

		time.Sleep(time.Second)
	}
}

func getSchedule() (sunrise, sunset time.Duration) {
	return mustParseDuration("7h"), mustParseDuration("23h")
}

type gradientTable []struct {
	c    colorful.Color
	time time.Duration
}

func getColor(gt gradientTable, t time.Duration) colorful.Color {
	if t < gt[0].time {
		return gt[0].c
	}

	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.time <= t && t <= c2.time {
			t := float64(t-c1.time) / float64(c2.time-c1.time)
			return c1.c.BlendLab(c2.c, t).Clamped() //super important to blend in a decent colorspace
		}
	}

	return gt[len(gt)-1].c
}
