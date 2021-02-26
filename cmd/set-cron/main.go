package main

import (
	"fmt"
	"github.com/salah93/adhango/salat"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"
)

const croncmd = "/usr/bin/crontab"

func main() {
	u, _ := user.Current()
	userID, _ := strconv.Atoi(u.Uid)

	today := time.Now()
	coords := salat.Coordinates{Latitude: 40.730610, Longitude: -73.935242}
	prayertimes, _ := salat.GetPrayerTimes(coords, today)

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	oldCronJobs, err := exec.Command(croncmd, "-l").Output()
	if err != nil {
		panic(err)
	}
	f.Write(oldCronJobs)

	for index, prayer := range prayertimes {
		var adhanFile string
		if index == salat.FAJR {
			adhanFile = "/home/salah/Projects/adhan-pi/static/azan-fajr.mp3"
		} else {
			adhanFile = "/home/salah/Projects/adhan-pi/static/azan2.mp3"
		}

		f.WriteString(fmt.Sprintf("\n%d %d * * * XDG_RUNTIME_DIR=/run/user/%d /usr/bin/ffplay -nodisp %s > /dev/null 2>&1", prayer.Minute, prayer.Hour, userID, adhanFile))
	}
	f.Sync()

	cmd := exec.Command(croncmd, f.Name())
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
