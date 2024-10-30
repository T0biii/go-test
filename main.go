package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// FirmwareDirectory enthält die Verzeichnisse für Firmwareversionen
var FirmwareDirectory = map[string]string{
	"unknown":          "latest_v2022_5",
	"201[0-9]":         "latest_v2022_5",
	"202[0-1]":         "latest_v2022_5",
	"2022\\.[0-4]":     "latest_v2022_5",
	"2022\\.5\\.[0-7]": "latest_v2022_5",
	"202[2-3]":         "latest_v2024_6",
	"2024\\.[0-5]":     "latest_v2024_6",
}

func main() {
	// Start HTTP Server
	http.HandleFunc("/firmware", firmwareHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// firmwareHandler prüft Header und leitet den Benutzer zur passenden Firmware um
func firmwareHandler(w http.ResponseWriter, r *http.Request) {
	// Prüfe User-Agent
	userAgent := r.Header.Get("User-Agent")
	autoupdater := userAgent == "Gluon Autoupdater (using libuclient)"

	// Hole Firmware-Version aus Header
	firmwareVersion := r.Header.Get("X-Firmware-Version")
	version := parseFirmwareVersion(firmwareVersion)

	// Bestimme das Zielverzeichnis basierend auf der Firmware-Version
	targetDirectory := determineTargetDirectory(version)

	log.Println(strconv.FormatBool(autoupdater) + " " + userAgent + " " + firmwareVersion + " " + version + " " + targetDirectory)

	// Leite entsprechend weiter oder antworte mit einer Nachricht
	if autoupdater {
		redirectURL := fmt.Sprintf("/%s/firmware.bin", targetDirectory)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

// parseFirmwareVersion extrahiert die Versionsnummer aus dem Header mithilfe von Regex
func parseFirmwareVersion(version string) string {
	versionRegex := regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)`)
	if versionRegex.MatchString(version) {
		return versionRegex.ReplaceAllString(version, "$1.$2.$3")
	}
	return "unknown"
}

// determineTargetDirectory wählt das Zielverzeichnis basierend auf der Firmware-Version
func determineTargetDirectory(version string) string {
	for pattern, dir := range FirmwareDirectory {
		match, _ := regexp.MatchString(pattern, version)
		if match {
			return dir
		}
	}
	return "keep"
}
