package instructions

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestGetDistancesFromRawSlice(t *testing.T) {
	type testcase struct {
		rawSlice           []string
		expectedMiles      float64
		expectedKilometers float64
	}

	testcases := []testcase{
		//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
		{rawSlice: []string{"1", "", "", "", "1", "", "", "", ""}, expectedMiles: 1, expectedKilometers: 1.609344},
		{rawSlice: []string{"1", "", "", "", "2", "", "", "", ""}, expectedMiles: 2, expectedKilometers: 3.218688},
		{rawSlice: []string{"1", "", "", "", "", "1", "", "", ""}, expectedMiles: 0.6213712, expectedKilometers: 1},
	}

	for _, tc := range testcases {
		miles, kilometers := getDistancesFromRawSlice(tc.rawSlice)
		milesDifference := math.Abs(miles - tc.expectedMiles)
		if milesDifference > 0.0000001 {
			t.Errorf("getDistancesFromRawSlice(%v): miles = %f, want %f", tc.rawSlice, miles, tc.expectedMiles)
		}
		kilometersDifference := math.Abs(kilometers - tc.expectedKilometers)
		if kilometersDifference > 0.0000001 {
			t.Errorf("getDistancesFromRawSlice(%v): kilometers %f, want %f", tc.rawSlice, kilometers, tc.expectedKilometers)
		}
	}
}

func TestGetCastsFromRawSlice(t *testing.T) {
	type testcase struct {
		rawSlice    []string
		expectNil   bool
		expectedMph float64
		expectedKph float64
	}

	testcases := []testcase{
		//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
		{rawSlice: []string{"1", "", "1", "", "", "", "", "", ""}, expectNil: false, expectedMph: 1, expectedKph: 1.609344},
		{rawSlice: []string{"1", "", "", "1", "", "", "", "", ""}, expectNil: false, expectedMph: 0.6213712, expectedKph: 1},
		{rawSlice: []string{"1", "", "", "", "", "", "", "", ""}, expectNil: true, expectedMph: 0, expectedKph: 0},
	}

	for _, tc := range testcases {
		mph, kph := getCastsFromRawSlice(tc.rawSlice)
		if tc.expectNil {
			if mph != nil {
				t.Errorf("getCastsFromRawSlice(%v): mph = %f, want nil", tc.rawSlice, *mph)
			}
			if kph != nil {
				t.Errorf("getCastsFromRawSlice(%v): kph = %f, want nil", tc.rawSlice, *kph)
			}
			continue
		}

		mphDifference := math.Abs(*mph - tc.expectedMph)
		if mphDifference > 0.0000001 {
			t.Errorf("getCastsFromRawSlice(%v): mph = %f, want %f", tc.rawSlice, *mph, tc.expectedMph)
		}
		kphDifference := math.Abs(*kph - tc.expectedKph)
		if kphDifference > 0.0000001 {
			t.Errorf("getCastsFromRawSlice(%v): kph = %f, want %f", tc.rawSlice, *kph, tc.expectedKph)
		}
	}
}

func TestParseSingleStringSlice(t *testing.T) {
	type testcase struct {
		rawSlice            []string
		expectedInstruction rawInstruction
	}
	//number,car_zero_begin_time,cast_mph,cast_kph,begin_distance_miles,begin_distance_km,instruction,pause_time,caution
	// 52,17:18:00,34,,0,,R Eastside Road,,
	exampleToday := time.Date(2023, 1, 17, 0, 0, 0, 0, time.UTC)
	aBeginTime := time.Date(2023, 1, 17, 17, 18, 0, 0, time.UTC)
	aCast := 34.0
	aCastKph := 54.717696
	cast22 := 22.0
	cast22Kph := 35.405568
	testcases := []testcase{
		{rawSlice: []string{"52", "17:18:00", "34", "", "0", "", "R Eastside Road", "", ""}, expectedInstruction: rawInstruction{
			number:                    52,
			carZeroBeginTime:          &aBeginTime,
			cast_mph:                  &aCast,
			cast_kph:                  &aCastKph,
			begin_distance_miles:      0,
			begin_distance_kilometers: 0,
			instructionText:           "R Eastside Road",
			pauseDuration:             0,
			cautions:                  0,
		}},

		{rawSlice: []string{"58", "", "22", "", "6", "", "Change speed at mileage", "", ""}, expectedInstruction: rawInstruction{
			number:                    58,
			carZeroBeginTime:          nil,
			cast_mph:                  &cast22,
			cast_kph:                  &cast22Kph,
			begin_distance_miles:      6,
			begin_distance_kilometers: 9.656064,
			instructionText:           "Change speed at mileage",
			pauseDuration:             0,
			cautions:                  0,
		}},
		{rawSlice: []string{"69", "", "", "", "12.436", "", "L at Yield Caution!!", "", "2"}, expectedInstruction: rawInstruction{
			number:                    69,
			carZeroBeginTime:          nil,
			cast_mph:                  nil,
			cast_kph:                  nil,
			begin_distance_miles:      12.436,
			begin_distance_kilometers: 20.013802,
			instructionText:           "L at Yield Caution!!",
			pauseDuration:             0,
			cautions:                  2,
		}},
		{rawSlice: []string{"69", "", "", "", "12.436", "", "L at stop. pause 12", "12", "2"}, expectedInstruction: rawInstruction{
			number:                    69,
			carZeroBeginTime:          nil,
			cast_mph:                  nil,
			cast_kph:                  nil,
			begin_distance_miles:      12.436,
			begin_distance_kilometers: 20.013802,
			instructionText:           "L at stop. pause 12",
			pauseDuration:             12 * time.Second,
			cautions:                  2,
		}},
	}

	for _, tc := range testcases {
		actualRaw := parseSingleStringSlice(tc.rawSlice, exampleToday)
		isEqual, msg := actualRaw.equals(tc.expectedInstruction)
		if !isEqual {
			t.Errorf("actualRaw not equal to expectedInstruction: %s", msg)
		}
	}
}

func (r rawInstruction) equals(other rawInstruction) (bool, string) {
	if r.number != other.number {
		return false, fmt.Sprintf("number: %d != %d", r.number, other.number)
	}
	if r.carZeroBeginTime != nil && r.carZeroBeginTime.Sub(*other.carZeroBeginTime) != 0 {
		return false, fmt.Sprintf("carZeroBeginTime: %v != %v", r.carZeroBeginTime, other.carZeroBeginTime)
	}
	if r.cast_mph != nil && math.Abs(*r.cast_mph-*other.cast_mph) > 0.0000001 {
		return false, fmt.Sprintf("cast_mph: %f != %f", *r.cast_mph, *other.cast_mph)
	}
	if r.cast_kph != nil && math.Abs(*r.cast_kph-*other.cast_kph) > 0.0000001 {
		return false, fmt.Sprintf("cast_kph: %f != %f", *r.cast_kph, *other.cast_kph)
	}
	if math.Abs(r.begin_distance_miles-other.begin_distance_miles) > 0.0000001 {
		return false, fmt.Sprintf("begin_distance_miles: %f != %f", r.begin_distance_miles, other.begin_distance_miles)
	}
	if math.Abs(r.begin_distance_kilometers-other.begin_distance_kilometers) > 0.0000001 {
		return false, fmt.Sprintf("begin_distance_kilometers: %f != %f", r.begin_distance_kilometers, other.begin_distance_kilometers)
	}
	if r.instructionText != other.instructionText {
		return false, fmt.Sprintf("instructionText: %s != %s", r.instructionText, other.instructionText)
	}
	if r.pauseDuration != other.pauseDuration {
		return false, fmt.Sprintf("pauseDuration: %v != %v", r.pauseDuration, other.pauseDuration)
	}
	if r.cautions != other.cautions {
		return false, fmt.Sprintf("cautions: %d != %d", r.cautions, other.cautions)
	}
	return true, ""
}
