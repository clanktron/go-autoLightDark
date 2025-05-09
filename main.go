package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

func estimateCoordinatesFromTimezone(tz string) (float64, float64) {
	tzCoords := map[string][2]float64{
		"America/Los_Angeles": {34.0522, -118.2437},
		"America/New_York":    {40.7128, -74.0060},
		"Europe/London":       {51.5074, -0.1278},
		"Asia/Tokyo":          {35.6895, 139.6917},
		"Australia/Sydney":    {-33.8688, 151.2093},
		// Add more as needed
	}

	if coords, ok := tzCoords[tz]; ok {
		return coords[0], coords[1]
	}

	fmt.Println("Unknown timezone, defaulting to lat=0, lon=0.")
	return 0.0, 0.0
}

func main() {
	link, err := os.Readlink("/etc/localtime")
	if err != nil {
		panic(err)
	}
	zone := filepath.Join("/", link)
	parts := strings.Split(filepath.Clean(link), string(os.PathSeparator))
	n := len(parts)
	if n < 2 {
		panic("unexpected timezone path format")
	}

	ianaZone := parts[n-2] + "/" + parts[n-1]
	location := time.Now().Location()
	latitude, longitude := estimateCoordinatesFromTimezone(ianaZone)
	fmt.Printf("tz %s \n", zone)
	fmt.Printf("Estimated location: Lat %.4f, Lon %.4f\n", latitude, longitude)

	for {
		now := time.Now().In(location)
		year, month, day := now.Date()

		sunriseTime, sunsetTime := sunrise.SunriseSunset(latitude, longitude, year, month, day)

		fmt.Printf("[INFO] Now: %s | Sunrise: %s | Sunset: %s\n", now.Format(time.RFC1123), sunriseTime.Format(time.Kitchen), sunsetTime.Format(time.Kitchen))

		var timeOfDay string
		if now.Before(sunriseTime) {
			timeOfDay = "Night"
		} else if now.Before(sunsetTime) {
			timeOfDay = "Day"
		} else {
			// Past both events: get sunrise/sunset for next day
			tomorrow := now.Add(24 * time.Hour)
			y, m, d := tomorrow.Date()
			sunriseTime, sunsetTime = sunrise.SunriseSunset(latitude, longitude, y, m, d)
			sunriseTime = sunriseTime.In(location)
			timeOfDay = "Night"
		}

		fmt.Printf("[TRIGGER] Executing action for %s\n", timeOfDay)
		runAction(timeOfDay)
		time.Sleep(30 * time.Second)
	}
}

func runAction(event string) {
	var cmd *exec.Cmd

	switch event {
	case "Night":
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark")
	case "Day":
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-light")
	default:
		fmt.Println("Unknown event")
		return
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
