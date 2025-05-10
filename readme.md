# go-autoLightDark

Daemon that auto sets system-wide light-dark mode on gtk systems based on time of day.

---

## Dependencies
* `gsettings`

## Build
```sh
go build
```

* estimates latitude/longitude based on timezone
* uses `go-sunrise` to calculate sunrise and sunset
* recalculates every 30s
* Uses `gsettings` cli to toggle gtk color-scheme
