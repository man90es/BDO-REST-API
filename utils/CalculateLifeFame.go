package utils

import (
	"strconv"
	"strings"
)

func CalculateLifeFame(specs [11]string) (lifeFame uint16) {
	for _, spec := range specs {
		spaceI := strings.Index(spec, " ")
		text := spec[0:spaceI]
		number, _ := strconv.Atoi(spec[spaceI+1:])

		if text == "Professional" {
			lifeFame += 90 + uint16(number)*3
		} else if text == "Artisan" {
			lifeFame += 120 + uint16(number)*3
		} else if text == "Master" {
			lifeFame += 150 + uint16(number)*3
		} else if text == "Guru" {
			lifeFame += 253
		}
	}

	return lifeFame + 1
}
