package main

import (
	"fmt"
	"github.com/salah93/adhango/salat"
	"time"
)

func main() {
	today := time.Now()
	coords := salat.Coordinates{Latitude: 40.730610, Longitude: -73.935242}
	prayertimes, _ := salat.GetPrayerTimes(coords, today)
	fmt.Println(prayertimes)
}
