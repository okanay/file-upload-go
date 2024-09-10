package utils

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s ~TOOK~ %s", name, elapsed.Round(time.Millisecond))
}
