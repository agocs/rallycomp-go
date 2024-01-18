package instructions

import (
	"log"
	"rallycomp-go/util"
	"strconv"
	"time"
)

type rawInstruction struct {
	number                    int
	carZeroBeginTime          *time.Time
	cast_mph                  *float64
	cast_kph                  *float64
	begin_distance_miles      float64
	begin_distance_kilometers float64
	instructionText           string
	pauseDuration             time.Duration
	cautions                  int
}

func parseCSVStringSlices(csvRaw [][]string) []rawInstruction {
	var instructions []rawInstruction
	today := time.Now()
	//TODO: set the timezone from config
	for i := 1; i < len(csvRaw); i++ {
		instruction := parseSingleStringSlice(csvRaw[i], today)
		instructions = append(instructions, instruction)
	}
	return instructions
}

func parseSingleStringSlice(rawSlice []string, today time.Time) rawInstruction {
	// csv header is
	//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
	instructionNumber, _ := strconv.Atoi(rawSlice[0])

	var carZeroBeginTime *time.Time
	if rawSlice[1] != "" {
		parsedTime, err := time.Parse("15:04:05", rawSlice[1])
		if err != nil {
			log.Fatal(err)
		}
		czbt := time.Date(today.Year(), today.Month(), today.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, today.Location())
		carZeroBeginTime = &czbt
	}

	cast_mph, cast_kph := getCastsFromRawSlice(rawSlice)
	begin_distance_miles, begin_distance_kilometers := getDistancesFromRawSlice(rawSlice)

	var pauseSeconds float64 = 0
	if rawSlice[7] != "" {
		pauseSeconds, _ = strconv.ParseFloat(rawSlice[7], 64)
	}
	cautions := 0
	if rawSlice[8] != "" {
		cautions, _ = strconv.Atoi(rawSlice[8])
	}
	return rawInstruction{
		number:                    instructionNumber,
		carZeroBeginTime:          carZeroBeginTime,
		cast_mph:                  cast_mph,
		cast_kph:                  cast_kph,
		begin_distance_miles:      begin_distance_miles,
		begin_distance_kilometers: begin_distance_kilometers,
		instructionText:           rawSlice[6],
		pauseDuration:             time.Second * time.Duration(pauseSeconds),
		cautions:                  cautions,
	}
}

func getCastsFromRawSlice(rawSlice []string) (*float64, *float64) {
	mph_string := rawSlice[2]
	kph_string := rawSlice[3]

	if mph_string == "" && kph_string == "" {
		return nil, nil
	}

	var cast_mph float64
	var cast_kph float64
	var err error
	if kph_string == "" {
		cast_mph, err = strconv.ParseFloat(mph_string, 32)
		if err != nil {
			//TODO: better error handling here
			panic(err)
		}
		cast_kph = util.MilesToKilometers(cast_mph)
	} else {
		cast_kph, err = strconv.ParseFloat(kph_string, 32)
		if err != nil {
			//TODO: better error handling here
			panic(err)
		}
		cast_mph = util.KilometersToMiles(cast_kph)
	}
	return &cast_mph, &cast_kph
}

func getDistancesFromRawSlice(rawSlice []string) (float64, float64) {
	miles_string := rawSlice[4]
	kilometers_string := rawSlice[5]

	if miles_string == "" && kilometers_string == "" {
		return 0, 0
	}

	var begin_distance_miles float64
	var begin_distance_kilometers float64
	var err error
	if miles_string == "" {
		begin_distance_kilometers, err = strconv.ParseFloat(kilometers_string, 64)
		if err != nil {
			panic(err)
		}
		begin_distance_miles = util.KilometersToMiles(begin_distance_kilometers)
	} else {
		begin_distance_miles, err = strconv.ParseFloat(miles_string, 64)
		if err != nil {
			panic(err)
		}
		begin_distance_kilometers = util.MilesToKilometers(begin_distance_miles)
	}
	return begin_distance_miles, begin_distance_kilometers
}
