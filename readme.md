# autoLightDark

Daemon that auto sets system-wide light-dark mode based on time of day.

---

## Dependencies
* `gsettings`

## Build
```sh
go build
```

* detects timezone via `/etc/localtime`
* estimates latitude/longitude based on timezone
* uses `go-sunrise` to calculate sunrise and sunset
* periodically recalculates
* Uses `gsettings` to toggle GNOME color-scheme
