//go:generate stringer -type=SaveStateType,scenarioGoalType,ParkFlag -output=savestate_strings.go

// Package sv4 implements an API to read/write savestates.
//
// Savestates are manipulated by having an in-memory byte slice that can be dumped to
// a file when needed. Not all bytes are known or implemented (yet). Checksums embedded
// in the savestate will be automatically re-calculated, everything else is left
// untouched.
//
// All byte positions and their meaning are derived from the excellent page at
// http://tid.rctspace.com/Sv4/SV4.html by James Hughes. A few bugs were found by other
// people and have been integrated in https://github.com/UnknownShadow200/RCTTechDepot-Archive.
package sv4

import (
	"errors"
	"strconv"
	"time"
)

// All savestates are fixed-sized bytes slices, only their RLE encoded files vary in size.
// This constant represents the size of an uncompressed savestate.
const SaveStateSize = 2065676

// Number of days in a month, ranging from March to October.
var DaysInMonth = []int{31, 30, 31, 30, 31, 31, 30, 31}

// A SaveState contains all information the game needs to restore a session.
type SaveState struct {
	data []byte
}

type SaveStateType int32

const (
	TypeRCT  SaveStateType = 108156
	TypeAACF SaveStateType = 110001
	TypeLL   SaveStateType = 120001
)

type scenarioGoalType uint8

const (
	GoalGuestsAndRating scenarioGoalType = iota
	ParkValue
	HaveFun
	Competition
	TenCoasters6Excitement
	Maintain
	MonthlyRideIncome
	TenCoasters7Excitement
	FiveCoasters
)

type ParkFlag uint32

const (
	ParkOpened                ParkFlag = 0x0001
	ProhibitLandModifitations ParkFlag = 0x0004
	ProhibitRemovingScenery   ParkFlag = 0x0008
	ShowRealNames             ParkFlag = 0x0010
	ProhibitAboveTreeLevel    ParkFlag = 0x0020
	LowIntensityPeeps         ParkFlag = 0x0040
	ProhibitAdvertising       ParkFlag = 0x0080
	CheatsDetected            ParkFlag = 0x0100
	HighIntensityPeeps        ParkFlag = 0x0200
	NoMoneyMode               ParkFlag = 0x0800
	GuestHighDifficulty       ParkFlag = 0x1000
	ForcedFreeEntry           ParkFlag = 0x2000
	RatingHighDifficulty      ParkFlag = 0x4000
	ProhibitToggleRealNames   ParkFlag = 0x8000
)

type scenarioGoal struct {
	Type  scenarioGoalType
	Years uint8

	// set one of them, never both
	MoneyGoal         uint32
	CoasterExcitement uint32

	// set one of them, never both
	GuestGoal        uint16
	MinCoasterLength uint16
}

// Creates a new SaveState from an uncompressed state file.
//
// Use the RLE decoder to uncompress before handing the raw data to this function.
// The function errors out if state is of the wrong length.
func NewSaveState(state []byte) (*SaveState, error) {
	s := &SaveState{state}

	if len(state) != SaveStateSize {
		return s, errors.New("Save states must be exactly " + strconv.Itoa(SaveStateSize) + " bytes in size.")
	}

	return s, nil
}

// Returns the current year.
func (s *SaveState) Year() int {
	val := int(s.readUint16(0x000000))

	return (val / 8) + 1
}

func (s *SaveState) Month() time.Month {
	val := int(s.readUint16(0x000000))

	return time.Month((val % 8) + 3)
}

// see https://github.com/OpenRCT2/OpenRCT2/blob/90fcc6f18/src/windows/game_bottom_toolbar.c#L227
func (s *SaveState) Days() int {
	ticks := int(s.readUint16(0x000002))
	days := DaysInMonth[s.Month()-3]

	return (((ticks * days) >> 16) & 0xFF) + 1
}

func (s *SaveState) ParkFlag(flag ParkFlag) bool {
	val := s.readUint32(0x19883C)

	return val&uint32(flag) > 0
}

func (s *SaveState) ParkRating() int {
	return int(s.readUint16(0x199108))
}

// func (s *SaveState) ParkRatingHistory() []int {
// 	return []int([]uint8(s.readBytes(0x19910A, 32)))
// }

func (s *SaveState) GuestCount() int {
	return int(s.readUint16(0x198C9C))
}

// func (s *SaveState) GuestCountHistory() []int {
// 	return []int([]uint8(s.readBytes(0x19912A, 32)))
// }

func (s *SaveState) HandymenColor() byte {
	return s.readByte(0x199025)
}

func (s *SaveState) MechanicsColor() byte {
	return s.readByte(0x199026)
}

func (s *SaveState) SecurityGuardsColor() byte {
	return s.readByte(0x199027)
}

/*
func (p *sv4parser) Parse() (*SaveState, error) {
	save.PocketCashBase = p.consumeInt16()
	save.HungerBase = p.consumeUint8()
	save.ThirstBase = p.consumeUint8()

	// scenario goal
	scGoal := p.consumeByte()

	save.ScenarioGoal.Type = scenarioGoalType(scGoal)
	save.ScenarioGoal.Years = p.consumeUint8()

	p.consumeBytes(2)

	if save.ScenarioGoal.Type == FiveCoasters {
		save.ScenarioGoal.CoasterExcitement = p.consumeUint32()
	} else {
		save.ScenarioGoal.MoneyGoal = p.consumeUint32()
	}

	if save.ScenarioGoal.Type == TenCoasters7Excitement {
		save.ScenarioGoal.MinCoasterLength = p.consumeUint16()
	} else {
		save.ScenarioGoal.GuestGoal = p.consumeUint16()
	}

	// 00 00 00 00 06 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 | 00000000ff00 | dynamite dunes, noch 6 wochen
	// 82 00 00 00 06 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 | ff000000ff00 | + freier eintritt für 2 wochen
	// 82 00 00 00 06 84 00 00 00 00 00 00 00 00 00 00 00 00 00 00 | ff000000ff0e | + minieisenbahn für 4 wochen
	// 82 00 00 82 06 84 00 00 00 00 00 00 00 00 00 00 00 00 00 00 | ff000006ff0e | + free hamburger für 2 wochen
	// 82 83 00 82 06 84 00 00 00 00 00 00 00 00 00 00 00 00 00 00 | ff180006ff0e | + free gokarts für 3 wochen

	fmt.Printf("%x\n", p.consumeBytes(20))
	fmt.Printf("%x\n", p.consumeBytes(6))

	// ???
	p.consumeBytes(16)

	p.where()

	return save, nil
}
*/

// see http://tid.rctspace.com/Checksum.html
func Checksum(encodedSavestate []byte, gameType SaveStateType) []byte {
	checksum := uint32(0)

	for i := 0; i < len(encodedSavestate); i++ {
		b := encodedSavestate[i]
		temp := checksum + uint32(b)

		checksum = (checksum & 0xFFFFFF00) | (temp & 0x000000FF)
		checksum = rol32(checksum, 3)
	}

	result := int32(checksum) + int32(gameType)

	return uint32ToBytes(result)
}
