package instructions

import (
	"os"
	"rallycomp-go/util"
	"strconv"
	"time"

	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type rawInstruction struct {
	number int
	carZeroBeginTime *time.Time
	cast_mph *float32
	cast_kph *float32
	begin_distance_miles float32
	begin_distance_kilometers float32
	instructionText string
	pauseDuration time.Duration
	cautions int
}

func parseCSVStringSlices(csvRaw [][]string) []rawInstruction {
	var instructions []rawInstruction
	for i := 1; i < len(csvRaw); i++ {
		// csv header is 
		//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
		instructions = append(instructions, rawInstruction{
			number: ,
			carZeroBeginTime: time.Now(),
			cast_mph: 0,
			cast_kph: 0,
			begin_distance_miles: 0,
			begin_distance_kilometers: 0,
			instructionText: csvRaw[i][0],
			pauseDuration: time.Duration(0),
			cautions: 0,
		})
	}
	return instructions
}

func parseSingleStringSlice(rawSlice []string) rawInstruction{
	// csv header is 
	//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
	instructionNumber, _ := strconv.Atoi(rawSlice[0])

	carZeroBeginTime = nil
	if rawSlice[1] != "" {
		carZeroBeginTime, _ := time.Parse(time.RFC3339, rawSlice[1])
	}

	cast_mph, cast_kph := getCastsFromRawSlice(rawSlice)
	begin_distance_miles, begin_distance_kilometers := getDistancesFromRawSlice(rawSlice)

	pauseDuration time.Duration = nil
	if rawSlice[7]	!= "" {
		pauseDuration, _ := time.ParseDuration(rawSlice[7])
	}
	cautions = 0
	if rawSlice[8] != "" {
		cautions, _ := strconv.Atoi(rawSlice[8])
	}
	return rawInstruction{
		number: instructionNumber,
		carZeroBeginTime: carZeroBeginTime,
		cast_mph: float32(cast_mph),
		cast_kph: float32(cast_kph),
		begin_distance_miles: float32(begin_distance_miles),
		begin_distance_kilometers: float32(begin_distance_kilometers),
		instructionText: rawSlice[6],
		pauseDuration: pauseDuration,
		cautions: cautions,
	}
}

func getCastsFromRawSlice(rawSlice []string) (*float32, *float32) {
	mph_string := rawSlice[2]
	kph_string := rawSlice[3]

	if mph_string == "" && kph_string == "" {
		return nil, nil
	}

	cast_mph *float32 = 0
	cast_kph *float32 = 0
	if mph_string == "" {
		cast_mph, err := strconv.ParseFloat(mph_string, 32)
		if err != nil {	
			//TODO: better error handling here
			panic(err)
		}
		cast_kph = util.MilesToKilometers(float32(cast_mph))
	} else {
		cast_kph, err := strconv.ParseFloat(kph_string, 32)
		if err != nil {
			//TODO: better error handling here
			panic(err)
		}
		cast_mph = util.KilometersToMiles(float32(cast_kph))
	}
	return cast_mph, cast_kph
}

func getDistancesFromRawSlice(rawSlice []string) (float32, float32) {
	miles_string := rawSlice[4]
	kilometers_string := rawSlice[5]

	if miles_string == "" && kilometers_string == "" {
		return 0, 0
	}

	begin_distance_miles float32 = 0
	begin_distance_kilometers float32 = 0
	if miles_string == "" {
		begin_distance_kilometers, err := strconv.ParseFloat(kilometers_string, 32)
		if err != nil {
			panic(err)
		}
		begin_distance_miles = util.KilometersToMiles(float32(begin_distance_kilometers))
	} else {
		begin_distance_miles, err := strconv.ParseFloat(miles_string, 32)
		if err != nil {
			panic(err)
		}
		begin_distance_kilometers = util.MilesToKilometers(float32(begin_distance_miles))
	}
	return begin_distance_miles, begin_distance_kilometers
}