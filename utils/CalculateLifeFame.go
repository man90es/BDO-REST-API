package utils

import (
	"bdo-rest-api/models"
	"reflect"
	"strconv"
	"strings"
)

func CalculateLifeFame(specs *models.Specs) (lifeFame uint16) {
	v := reflect.ValueOf(*specs)

	for i := 0; i < v.NumField(); i++ {
		spec := v.Field(i).String()
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
			lifeFame += 240 + uint16(number)*3
		}
	}

	return lifeFame + 1
}
