package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

func main() {
	// Coordinates for your location
	latitude := 34.0522
	longitude := -118.2437

	// Load local timezone
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	for {
		now := time.Now().In(loc)
		year, month, day := now.Date()

		sunriseTime, sunsetTime := sunrise.SunriseSunset(latitude, longitude, year, month, day)
		sunriseTime = sunriseTime.In(loc)
		sunsetTime = sunsetTime.In(loc)

		fmt.Printf("[INFO] Now: %s | Sunrise: %s | Sunset: %s\n", now.Format(time.RFC1123), sunriseTime.Format(time.Kitchen), sunsetTime.Format(time.Kitchen))

		// Check if next event is sunrise or sunset
		var nextEvent string
		// var eventTime time.Time

		if now.Before(sunriseTime) {
			nextEvent = "sunrise"
		} else if now.Before(sunsetTime) {
			nextEvent = "sunset"
		} else {
			// Past both events: get sunrise/sunset for next day
			tomorrow := now.Add(24 * time.Hour)
			y, m, d := tomorrow.Date()
			sunriseTime, sunsetTime = sunrise.SunriseSunset(latitude, longitude, y, m, d)
			sunriseTime = sunriseTime.In(loc)
			nextEvent = "sunrise"
			// eventTime = sunriseTime
		}

		fmt.Printf("[TRIGGER] Executing action for %s\n", nextEvent)
		runAction(nextEvent)
		// lastTriggered = nextEvent

		// Sleep for 30 seconds
		time.Sleep(30 * time.Second)
	}
}

func runAction(event string) {
	var cmd *exec.Cmd

	// Define actions based on event
	switch event {
	case "sunrise":
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark")
	case "sunset":
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
