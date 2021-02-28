package salat

// 5 daily prayers
const (
	FAJR = iota
	DHUHR
	ASR
	MAGHRIB
	ISHA
)

//Prayer - prayer object
type Prayer struct {
	Hour   int
	Minute int
}

//Coordinates :
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

//PrayerTimes -
type PrayerTimes [5]Prayer

type prayerTimingsData struct {
	Fajr     string
	Dhuhr    string
	Asr      string
	Maghrib  string
	Isha     string
	Sunrise  string `json:"-"`
	Sunset   string `json:"-"`
	Imsak    string `json:"-"`
	Midnight string `json:"-"`
}

type prayerAPIData struct {
	Timings *prayerTimingsData
	Date    interface{} `json:"-"`
	Meta    interface{} `json:"-"`
}
type prayerAPIResponse struct {
	Code   string `json:"-"`
	Status string `json:"-"`
	Data   []prayerAPIData
}
