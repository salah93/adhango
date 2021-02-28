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
		req = nil
	} else {
		q := req.URL.Query()
		q.Add("latitude", fmt.Sprintf("%f", r.coords.Latitude))
		q.Add("longitude", fmt.Sprintf("%f", r.coords.Longitude))
		q.Add("method", "02")
		q.Add("month", strconv.Itoa(int(r.date.Month())))
		q.Add("year", strconv.Itoa(int(r.date.Year())))
		req.URL.RawQuery = q.Encode()
	}
	return req, err
}

func (r request) getTimings() (*prayerTimingsData, error) {
	var data prayerAPIResponse
	var timings *prayerTimingsData
	req, err := r.buildRequest()
	if err == nil {
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				err = errors.New("api error")
			} else {
				json.NewDecoder(resp.Body).Decode(&data)
				timings = data.Data[r.date.Day()-1].Timings
			}
		}
	}
	return timings, err
}

//GetPrayerTimes -
func GetPrayerTimes(coords Coordinates, date time.Time) (PrayerTimes, error) {
	var prayerTimes PrayerTimes
	r := request{coords: coords, date: date}
	timings, err := r.getTimings()
	if err == nil {
		const timeLayout = "15:04 (MST)"
		fajrTime, _ := time.Parse(timeLayout, timings.Fajr)
		dhuhrTime, _ := time.Parse(timeLayout, timings.Dhuhr)
		asrTime, _ := time.Parse(timeLayout, timings.Asr)
		maghribTime, _ := time.Parse(timeLayout, timings.Maghrib)
		ishaTime, _ := time.Parse(timeLayout, timings.Isha)

		prayerTimes = PrayerTimes{
			Prayer{Hour: fajrTime.Hour(), Minute: fajrTime.Minute()},
			Prayer{Hour: dhuhrTime.Hour(), Minute: dhuhrTime.Minute()},
			Prayer{Hour: asrTime.Hour(), Minute: asrTime.Minute()},
			Prayer{Hour: maghribTime.Hour(), Minute: maghribTime.Minute()},
			Prayer{Hour: ishaTime.Hour(), Minute: ishaTime.Minute()},
		}
	}
	return prayerTimes, err
}
