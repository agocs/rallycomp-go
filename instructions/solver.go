package instructions

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type instruction struct {
	prevInstruction         *instruction
	nextInstruction         *instruction
	number                  int
	carZeroBeginTime        time.Time
	castKph                 float64
	beginDistanceKilometers float64
	instructionDuration     time.Duration
	instructionText         string
	pauseDuration           time.Duration
	cautions                int
}

func (i *instruction) setInstruction(raw rawInstruction, prev *instruction) *instruction {
	// do the linked list stuff
	i.prevInstruction = prev
	if prev != nil {
		prev.nextInstruction = i
	}

	i.number = raw.number

	//every instruction should have a distance
	i.beginDistanceKilometers = raw.begin_distance_kilometers

	// if the instruction doesn't have a CAST, use the previous instruction's CAST
	if raw.cast_kph == nil {
		i.castKph = prev.castKph
	} else {
		i.castKph = *raw.cast_kph
	}

	// we have a CAST and a distance, so we can calculate the duration of this instruction
	i.instructionDuration = time.Duration(i.beginDistanceKilometers/i.castKph*3600) * time.Second

	// if this instruction does not have a carZeroBeginTime, use the previous instruction's carZeroBeginTime plus the duration of the previous instruction
	if raw.carZeroBeginTime == nil {
		i.carZeroBeginTime = prev.carZeroBeginTime.Add(prev.instructionDuration).Add(prev.pauseDuration)
	} else {
		i.carZeroBeginTime = *raw.carZeroBeginTime
		// Validate the provided carZeroBeginTime
		theoreticalCarZeroBeginTime := prev.carZeroBeginTime.Add(prev.instructionDuration).Add(prev.pauseDuration)
		if i.carZeroBeginTime != theoreticalCarZeroBeginTime {
			log.Printf("Instruction %d has a carZeroBeginTime of %s, but the theoretical carZeroBeginTime is %s", i.number, i.carZeroBeginTime, theoreticalCarZeroBeginTime)
		}
	}

	i.instructionText = raw.instructionText
	i.pauseDuration = raw.pauseDuration
	i.cautions = raw.cautions
	return i

}

func (i instruction) String() string {
	outBuilder := strings.Builder{}
	for {
		outBuilder.WriteString(fmt.Sprintf("%d, %v, %v, %s", i.number, i.carZeroBeginTime, i.beginDistanceKilometers, i.instructionText))
		if i.nextInstruction == nil {
			break
		}
		i = *i.nextInstruction
	}
	return outBuilder.String()
}

func solveRawInstructions(rawInstructions []rawInstruction) instruction {
	headInstruction := &instruction{}
	headInstruction.setInstruction(rawInstructions[0], nil)
	tailInstruction := headInstruction
	for i := 1; i < len(rawInstructions); i++ {
		thisInstruction := &instruction{}
		thisInstruction.setInstruction(rawInstructions[i], tailInstruction)
		tailInstruction = thisInstruction
	}
	return *headInstruction
}
