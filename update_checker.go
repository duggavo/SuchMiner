package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

func checkForUpdates() {
	res, err := http.Get("https://raw.githubusercontent.com/duggavo/SuchMiner/refs/heads/main/version.go")
	if err != nil || res.StatusCode != 200 {
		Log.Println(Yellow+"update checking failed:", err, "status code:", res.StatusCode, Reset)
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		Log.Println(Yellow+"downloading update check file failed:", err, Reset)
		return
	}

	s := strings.SplitN(string(body), "const VERSION = \"", 2)
	if len(s) != 2 {
		Log.Println(Yellow+"failed to parse version file", Reset)
		return
	}

	version := strings.SplitN(s[1], "\"", 2)[0]

	if len(version) > 10 {
		Log.Println(Yellow + Bright + "version checking failed, malformed version: " + strconv.Quote(version) + Reset)
		return
	}

	if debug {
		Log.Println("update checker found latest version:", version)
	}
	if version != VERSION {
		Log.Println(Yellow + Bright + "You're running an outdated version of SuchMiner! (" + VERSION + " -> " + version + ")" + Reset)
		Log.Println(Yellow + Bright + "Download updated version from https://github.com/duggavo/wowrig" + Reset)
	}
}
