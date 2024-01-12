package instructions

import (
	"log"
	"time"
)

// type rawInstruction struct {
// 	number int
// 	carZeroBeginTime *time.Time
// 	cast_mph *float32
// 	cast_kph *float32
// 	begin_distance_miles float32
// 	begin_distance_kilometers float32
// 	instructionText string
// 	pauseDuration time.Duration
// 	cautions int
// }

type instruction struct {
	prevInstruction           *instruction
	nextInstruction           *instruction
	number                    int
	carZeroBeginTime          time.Time
	cast_kph                  float32
	begin_distance_kilometers float32
	instructionDuration       time.Duration
	instructionText           string
	pauseDuration             time.Duration
	cautions                  int
}

func (i *instruction) setInstruction(raw rawInstruction, prev *instruction) *instruction {
	// do the linked list stuff
	i.prevInstruction = prev
	if prev != nil {
		prev.nextInstruction = i
	}

	i.number = raw.number

	//every instruction should have a distance
	i.begin_distance_kilometers = raw.begin_distance_kilometers

	// if the instruction doesn't have a CAST, use the previous instruction's CAST
	if raw.cast_kph == nil {
		i.cast_kph = prev.cast_kph
	} else {
		i.cast_kph = *raw.cast_kph
	}

	// we have a CAST and a distance, so we can calculate the duration of this instruction
	i.instructionDuration = time.Duration(i.begin_distance_kilometers/i.cast_kph*3600) * time.Second

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
