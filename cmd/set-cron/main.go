package main

import (
	"fmt"
	"github.com/salah93/adhango/salat"
	"github.com/salah93/go-cron"
	"os/exec"
	"os/user"
	"strconv"
	"time"
)

const croncmd = "/usr/bin/crontab"

func main() {
	u, _ := user.Current()
	userID, _ := strconv.Atoi(u.Uid)

	const identifyingComment = "go-adhan"
	today := time.Now()
	coords := salat.Coordinates{Latitude: 40.730610, Longitude: -73.935242}
	prayertimes, _ := salat.GetPrayerTimes(coords, today)

	job := cron.NewJob()
	job.RemoveItemsByComment(identifyingComment)
	for index, prayer := range prayertimes {
		var adhanFile string
		if index == salat.FAJR {
			adhanFile = "/home/salah/Projects/adhan-pi/static/azan-fajr.mp3"
		} else {
			adhanFile = "/home/salah/Projects/adhan-pi/static/azan2.mp3"
		}
		cmd := exec.Command("/usr/bin/ffplay", "-nodisp", adhanFile)
		cmd.Env = append(cmd.Env, fmt.Sprintf("XDG_RUNTIME_DIR=/run/user/%d", userID))
		job.AddItem(
			&cron.Item{
				Command: cmd,
				Comment: identifyingComment,
				Time: &cron.ItemTime{
					Minute:     strconv.Itoa(prayer.Minute),
					Hour:       strconv.Itoa(prayer.Hour),
					DayOfMonth: "*",
					Month:      "*",
					WeekDay:    "*",
				},
			},
		)
	}
	job.Save()
}
