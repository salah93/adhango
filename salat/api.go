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

type request struct {
	coords Coordinates
	date   time.Time
}

func (r request) buildRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", APIURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("latitude", fmt.Sprintf("%f", r.coords.Latitude))
	q.Add("longitude", fmt.Sprintf("%f", r.coords.Longitude))
	q.Add("method", "02")
	q.Add("month", strconv.Itoa(int(r.date.Month())))
	q.Add("year", strconv.Itoa(int(r.date.Year())))
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (r request) getTimings() (*prayerTimingsData, error) {
	req, _ := r.buildRequest()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("api error")
	}
	var data prayerAPIResponse
	json.NewDecoder(resp.Body).Decode(&data)
	return &data.Data[r.date.Day()-1].Timings, nil
}

//GetPrayerTimes -
func GetPrayerTimes(coords Coordinates, date time.Time) (PrayerTimes, error) {
	r := request{coords: coords, date: date}
	timings, _ := r.getTimings()

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
