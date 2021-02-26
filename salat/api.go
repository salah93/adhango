package salat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//APIURL -
const APIURL = "https://api.aladhan.com/v1/calendar"

var client *http.Client

func init() {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client = &http.Client{Timeout: 10 * time.Second, Transport: tr}
}

//GetPrayerTimes -
func GetPrayerTimes(coords Coordinates, date time.Time) (PrayerTimes, error) {
	req, err := http.NewRequest("GET", APIURL, nil)
	if err != nil {
		return PrayerTimes{}, err
	}

	q := req.URL.Query()
	q.Add("latitude", fmt.Sprintf("%f", coords.Latitude))
	q.Add("longitude", fmt.Sprintf("%f", coords.Longitude))
	q.Add("method", "02")

	q.Add("month", strconv.Itoa(int(date.Month())))
	q.Add("year", strconv.Itoa(int(date.Year())))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return PrayerTimes{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PrayerTimes{}, errors.New("api error")
	}
	var data prayerAPIResponse
	json.NewDecoder(resp.Body).Decode(&data)
	var timings = data.Data[date.Day()-1].Timings

	const layout = "15:04 (MST)"
	fajrTime, _ := time.Parse(layout, timings.Fajr)
	dhuhrTime, _ := time.Parse(layout, timings.Dhuhr)
	asrTime, _ := time.Parse(layout, timings.Asr)
	maghribTime, _ := time.Parse(layout, timings.Maghrib)
	ishaTime, _ := time.Parse(layout, timings.Isha)
	return PrayerTimes{
		Prayer{Hour: fajrTime.Hour(), Minute: fajrTime.Minute()},
		Prayer{Hour: dhuhrTime.Hour(), Minute: dhuhrTime.Minute()},
		Prayer{Hour: asrTime.Hour(), Minute: asrTime.Minute()},
		Prayer{Hour: maghribTime.Hour(), Minute: maghribTime.Minute()},
		Prayer{Hour: ishaTime.Hour(), Minute: ishaTime.Minute()},
	}, nil
}
