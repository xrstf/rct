package rct

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const SaveStateSize = 2065676

var DaysInMonth = []int{31, 30, 31, 30, 31, 31, 30, 31}

type position2D struct {
	X uint16
	Y uint16
}

type position3D struct {
	position2D

	Z uint16
}

type report struct {
	RideConstruction int32
	RideOperation    int32
	LandPurchase     int32
	Landscaping      int32
	ParkTickets      int32
	RideTickets      int32
	ShopSales        int32
	ShopStock        int32
	FoodSales        int32
	FoodStock        int32
	StaffWages       int32
	Marketing        int32
	Research         int32
	LoanInterest     int32
}

type bannerMenu struct {
	Standard bool
	Jungle   bool
	Roman    bool
	Egyptian bool
	Mining   bool
	Jurassic bool
	Asian    bool
	Snow     bool
	Space    bool
}

type researchOptions struct {
	Rollercoasters   bool
	ThrillRides      bool
	GentleRides      bool
	Shops            bool
	Theming          bool
	RideImprovements bool
}

type researchTask struct {
	Item     byte
	Ride     byte
	Category byte
	Flags    byte
}

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

type SaveState struct {
	Month time.Month
	Year  int

	ticks int
	Day   int

	timeTicks uint32

	rngState []byte

	gameMap []byte

	counter uint32

	sprites                    []byte
	NextSpriteId               uint16
	FirstVehicleSpriteId       uint16
	FirstPeepSpriteId          uint16
	FirstDuckSpriteId          uint16
	FirstTrashSpriteId         uint16
	FirstOversizedRideSpriteId uint16
	SpriteSlotsAvailable       uint16
	VehicleCount               uint16
	PeepCount                  uint16
	DuckCount                  uint16
	TrashCount                 uint16
	OversizedRideCount         uint16

	Cash uint32
	Loan uint32

	ParkOpened                bool
	ProhibitLandModifitations bool
	ProhibitRemovingScenery   bool
	ShowRealNames             bool
	ProhibitAboveTreeLevel    bool
	LowIntensityPeeps         bool
	ProhibitAdvertising       bool
	CheatsDetected            bool
	HighIntensityPeeps        bool
	NoMoneyMode               bool
	GuestHighDifficulty       bool
	ForcedFreeEntry           bool
	RatingHighDifficulty      bool
	ProhibitToggleRealNames   bool

	ParkEntryFee       uint16
	ParkEntryDirection byte
	ParkEntryPos       position3D

	PeepEntry1Pos       position2D
	PeepEntry1Height    byte
	PeepEntry1Direction byte

	PeepEntry2Pos       position2D
	PeepEntry2Height    byte
	PeepEntry2Direction byte

	ResearchRate byte

	ridesMenu         []byte
	ridesMenuReserved []byte

	vehicleMenu         []byte
	vehicleMenuReserved []byte

	rideFeatures1         []byte
	rideFeatures1Reserved []byte

	rideFeatures2         []byte
	rideFeatures2Reserved []byte

	GuestCount uint16

	Reports []report

	ColorHandymen  byte
	ColorMechanics byte
	ColorSecurity  byte

	sceneryMenu []byte

	BannerMenu bannerMenu

	ParkRating        uint16
	ParkRatingHistory []uint8
	ParkGuestHistory  []uint8

	ResearchOptions  researchOptions
	ResearchProgress byte
	LastResearch     researchTask
	NextResearch     researchTask

	OwnedLand uint16

	MaxLoan uint32

	PocketCashBase int16
	HungerBase     uint8
	ThirstBase     uint8
	ScenarioGoal   scenarioGoal
}

type sv4parser struct {
	pos  int
	data []byte
}

// see http://tid.rctspace.com/Sv4/SV4.html
func ParseSaveState(state []byte) (*SaveState, error) {
	parser := sv4parser{0, state}

	return parser.Parse()
}

func (p *sv4parser) Parse() (*SaveState, error) {
	if len(p.data) != SaveStateSize {
		return nil, errors.New("Save states must be exactly " + strconv.Itoa(SaveStateSize) + " bytes in size.")
	}

	save := &SaveState{}

	// first two bytes contain the number of months, counting up (Mar-Oct)
	m := int(p.consumeUint16())

	save.Month = time.Month((m % 8) + 3)
	save.Year = (m / 8) + 1

	// next two bytes are the tick count, representing the current day
	// see https://github.com/OpenRCT2/OpenRCT2/blob/90fcc6f18/src/windows/game_bottom_toolbar.c#L227
	ticks := int(p.consumeUint16())
	days := DaysInMonth[save.Month-3]

	save.ticks = ticks
	save.Day = (((ticks * days) >> 16) & 0xFF) + 1

	// game time counter, counts in steps of 1. does not begin at zero when scenario starts
	save.timeTicks = p.consumeUint32()

	// psuedo random numbers which are related to each other.
	save.rngState = p.consumeBytes(8)

	// the game map
	// save.gameMap = p.consumeBytes(393216)
	p.consumeBytes(393216)

	// an incrementing counter, two upper bits not used - not used as a trigger
	save.counter = p.consumeUint32()

	// sprite data structures
	p.consumeBytes(1280000)

	// sprite summary
	save.NextSpriteId = p.consumeUint16()
	save.FirstVehicleSpriteId = p.consumeUint16()
	save.FirstPeepSpriteId = p.consumeUint16()
	save.FirstDuckSpriteId = p.consumeUint16()
	save.FirstTrashSpriteId = p.consumeUint16()
	save.FirstOversizedRideSpriteId = p.consumeUint16()
	save.SpriteSlotsAvailable = p.consumeUint16()
	save.VehicleCount = p.consumeUint16()
	save.PeepCount = p.consumeUint16()
	save.DuckCount = p.consumeUint16()
	save.TrashCount = p.consumeUint16()
	save.OversizedRideCount = p.consumeUint16()

	// park name string offset
	p.consumeUint32()

	// ???
	p.consumeUint32()

	// cash & loan
	save.Cash = p.consumeUint32()
	save.Loan = p.consumeUint32()

	// park flags
	flags := p.consumeUint32()

	// yes, 0x2 and 0x400 are not used / defined here
	save.ParkOpened = flags&0x1 > 0
	save.ProhibitLandModifitations = flags&0x4 > 0
	save.ProhibitRemovingScenery = flags&0x8 > 0
	save.ShowRealNames = flags&0x10 > 0
	save.ProhibitAboveTreeLevel = flags&0x20 > 0
	save.LowIntensityPeeps = flags&0x40 > 0
	save.ProhibitAdvertising = flags&0x80 > 0
	save.CheatsDetected = flags&0x100 > 0
	save.HighIntensityPeeps = flags&0x200 > 0
	save.NoMoneyMode = flags&0x800 > 0
	save.GuestHighDifficulty = flags&0x1000 > 0
	save.ForcedFreeEntry = flags&0x2000 > 0
	save.RatingHighDifficulty = flags&0x4000 > 0
	save.ProhibitToggleRealNames = flags&0x8000 > 0

	// park info
	save.ParkEntryFee = p.consumeUint16()
	save.ParkEntryPos = position3D{
		position2D{
			p.consumeUint16(),
			p.consumeUint16(),
		},
		p.consumeUint16(),
	}
	save.ParkEntryDirection = p.consumeByte()

	// ???
	p.consumeByte()

	// peep entry #1
	save.PeepEntry1Pos = position2D{
		p.consumeUint16(),
		p.consumeUint16(),
	}
	save.PeepEntry1Height = p.consumeByte()
	save.PeepEntry1Direction = p.consumeByte()

	// peep entry #2
	save.PeepEntry2Pos = position2D{
		p.consumeUint16(),
		p.consumeUint16(),
	}
	save.PeepEntry2Height = p.consumeByte()
	save.PeepEntry2Direction = p.consumeByte()

	// ???
	p.consumeByte()

	// research rate
	save.ResearchRate = p.consumeByte()

	// ???
	p.consumeBytes(4)

	// menus
	save.ridesMenu = p.consumeBytes(11)
	save.ridesMenuReserved = p.consumeBytes(21)

	save.vehicleMenu = p.consumeBytes(12)
	save.vehicleMenuReserved = p.consumeBytes(20)

	// features
	save.rideFeatures1 = p.consumeBytes(356)
	save.rideFeatures1Reserved = p.consumeBytes(156)

	save.rideFeatures2 = p.consumeBytes(356)
	save.rideFeatures2Reserved = p.consumeBytes(156)

	save.GuestCount = p.consumeUint16()

	// ???
	p.consumeBytes(2)

	// part financial reports
	save.Reports = make([]report, 16)

	for i := 0; i < len(save.Reports); i = i + 1 {
		r := &save.Reports[i]

		r.RideConstruction = p.consumeInt32()
		r.RideOperation = p.consumeInt32()
		r.LandPurchase = p.consumeInt32()
		r.Landscaping = p.consumeInt32()
		r.ParkTickets = p.consumeInt32()
		r.RideTickets = p.consumeInt32()
		r.ShopSales = p.consumeInt32()
		r.ShopStock = p.consumeInt32()
		r.FoodSales = p.consumeInt32()
		r.FoodStock = p.consumeInt32()
		r.StaffWages = p.consumeInt32()
		r.Marketing = p.consumeInt32()
		r.Research = p.consumeInt32()
		r.LoanInterest = p.consumeInt32()
	}

	// ???
	p.consumeBytes(5)

	// color settings
	save.ColorHandymen = p.consumeByte()
	save.ColorMechanics = p.consumeByte()
	save.ColorSecurity = p.consumeByte()

	// scenery menu
	save.sceneryMenu = p.consumeBytes(128)

	// banner menu
	bannerFlags := p.consumeUint16()

	save.BannerMenu.Standard = bannerFlags&0x1 > 0
	save.BannerMenu.Jungle = bannerFlags&0x2 > 0
	save.BannerMenu.Roman = bannerFlags&0x4 > 0
	save.BannerMenu.Egyptian = bannerFlags&0x8 > 0
	save.BannerMenu.Mining = bannerFlags&0x10 > 0
	save.BannerMenu.Jurassic = bannerFlags&0x20 > 0
	save.BannerMenu.Asian = bannerFlags&0x40 > 0
	save.BannerMenu.Snow = bannerFlags&0x80 > 0
	save.BannerMenu.Space = bannerFlags&0x100 > 0

	// ???
	p.consumeBytes(94)

	// park rating
	save.ParkRating = p.consumeUint16()
	save.ParkRatingHistory = []uint8(p.consumeBytes(32))
	save.ParkGuestHistory = []uint8(p.consumeBytes(32))

	// research options
	rFlags := p.consumeByte()

	save.ResearchOptions.Rollercoasters = rFlags&0x1 > 0
	save.ResearchOptions.ThrillRides = rFlags&0x2 > 0
	save.ResearchOptions.GentleRides = rFlags&0x4 > 0
	save.ResearchOptions.Shops = rFlags&0x8 > 0
	save.ResearchOptions.Theming = rFlags&0x10 > 0
	save.ResearchOptions.RideImprovements = rFlags&0x20 > 0

	save.ResearchProgress = p.consumeByte()
	save.LastResearch = researchTask{
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
	}

	// skip research items for now
	p.consumeBytes(1000)

	save.NextResearch = researchTask{
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
	}

	// ???
	p.consumeBytes(6)

	// Cheat detection: count of owned land
	save.OwnedLand = p.consumeUint16()

	// ???
	p.consumeBytes(4)

	// Max loan amount
	save.MaxLoan = p.consumeUint32()

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

func (p *sv4parser) where() {
	fmt.Printf("current position = 0x%X\n", p.pos)
}

func (p *sv4parser) consumeByte() byte {
	return p.consumeBytes(1)[0]
}

func (p *sv4parser) consumeBytes(num int) []byte {
	defer (func() { p.pos = p.pos + num })()

	return p.data[p.pos:(p.pos + num)]
}

func (p *sv4parser) consumeUint8() uint8 {
	return uint8(p.consumeByte())
}

func (p *sv4parser) consumeUint16() uint16 {
	return binary.LittleEndian.Uint16(p.consumeBytes(2))
}

func (p *sv4parser) consumeInt16() int16 {
	b := p.consumeBytes(2)
	buf := bytes.NewReader(b)
	result := int16(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}

func (p *sv4parser) consumeUint32() uint32 {
	return binary.LittleEndian.Uint32(p.consumeBytes(4))
}

func (p *sv4parser) consumeInt32() int32 {
	b := p.consumeBytes(4)
	buf := bytes.NewReader(b)
	result := int32(0)

	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return result
}
