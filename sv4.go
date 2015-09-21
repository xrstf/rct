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

	ColorHandymen byte
	ColorMechanics byte
	ColorSecurity byte

	sceneryMenu []byte
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

	p.where()

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

	return save, nil
}

func (p *sv4parser) where() {
	fmt.Printf("position = %X\n", p.pos)
}

func (p *sv4parser) consumeByte() byte {
	return p.consumeBytes(1)[0]
}

func (p *sv4parser) consumeBytes(num int) []byte {
	defer (func() { p.pos = p.pos + num })()

	return p.data[p.pos:(p.pos + num)]
}

func (p *sv4parser) consumeUint16() uint16 {
	return binary.LittleEndian.Uint16(p.consumeBytes(2))
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
