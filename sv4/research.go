// Copyright (c) 2015, xrstf | MIT licensed

//go:generate stringer -type=ResearchRate,ResearchFlag -output=research_strings.go

package sv4

type ResearchTask struct {
	Item     byte
	Ride     byte
	Category byte
	Flags    byte
}

type ResearchFlag byte

const (
	RollercoastersResearch   ResearchFlag = 0x01
	ThrillRidesResearch      ResearchFlag = 0x02
	GentleRidesResearch      ResearchFlag = 0x04
	ShopsResearch            ResearchFlag = 0x08
	ThemingResearch          ResearchFlag = 0x10
	RideImprovementsResearch ResearchFlag = 0x20
)

type ResearchRate uint8

const (
	NoResearch ResearchRate = iota
	MinimumResearch
	NormalResearch
	MaximumResearch
)

func (s *SaveState) ResearchRate() ResearchRate {
	return ResearchRate(s.readUint8(0x198857))
}

func (s *SaveState) ResearchFlag(flag ResearchFlag) bool {
	val := s.readByte(0x19914A)

	return val&byte(flag) > 0
}

/*
func (p *sv4parser) Parse() (*SaveState, error) {
	save.ResearchProgress = p.consumeByte()
	save.LastResearch = ResearchTask{
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
	}

	save.NextResearch = ResearchTask{
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
		p.consumeByte(),
	}
*/
