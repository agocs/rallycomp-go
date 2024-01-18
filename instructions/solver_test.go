package instructions

import (
	"testing"
	"time"
)

// type instruction struct {
// 	prevInstruction         *instruction
// 	nextInstruction         *instruction
// 	number                  int
// 	carZeroBeginTime        time.Time
// 	castKph                 float64
// 	beginDistanceKilometers float64
// 	instructionDuration     time.Duration
// 	instructionText         string
// 	pauseDuration           time.Duration
// 	cautions                int
// }

// type rawInstruction struct {
// 	number                    int
// 	carZeroBeginTime          *time.Time
// 	cast_mph                  *float64
// 	cast_kph                  *float64
// 	begin_distance_miles      float64
// 	begin_distance_kilometers float64
// 	instructionText           string
// 	pauseDuration             time.Duration
// 	cautions                  int
// }

func TestInstruction_SetInstruction(t *testing.T) {
	instructonZero := &instruction{
		prevInstruction:         nil,
		nextInstruction:         nil,
		number:                  1,
		carZeroBeginTime:        time.Date(2023, 1, 17, 13, 00, 0, 0, time.UTC),
		castKph:                 10,
		beginDistanceKilometers: 0,
		instructionDuration:     0,
		instructionText:         "Start",
		pauseDuration:           0,
		cautions:                0,
	}

	rawInstructionOne := rawInstruction{
		number:                    2,
		carZeroBeginTime:          nil,
		begin_distance_kilometers: 10,
		instructionText:           "Go 10km",
		pauseDuration:             0,
		cautions:                  0,
	}

	instructionOne := &instruction{}
	instructionOne.setInstruction(rawInstructionOne, instructonZero)

	if instructionOne.prevInstruction != instructonZero {
		t.Errorf("instructionOne.prevInstruction = %v; want %v", instructionOne.prevInstruction, instructonZero)
	}

	if instructonZero.nextInstruction != instructionOne {
		t.Errorf("instructionZero.nextInstruction = %v; want %v", instructonZero.nextInstruction, instructionOne)
	}

	//TODO: confirm the instruction durations are set correctly
	if instructonZero.instructionDuration != 0 {
		t.Errorf("instructionZero.instructionDuration = %v; want %v", instructionZero.instructionDuration, 0)
	}

}
