package sv4

func (s *SaveState) Cash() int {
	return int(s.readUint32(0x198834))
}

func (s *SaveState) Loan() int {
	return int(s.readUint32(0x198838))
}

func (s *SaveState) MaxLoan() int {
	return int(s.readUint32(0x199548))
}

func (s *SaveState) ParkEntryFee() int {
	return int(s.readUint16(0x198840))
}

type FinanceReport struct {
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

func (s *SaveState) FinanceReports() []FinanceReport {
	reports := make([]FinanceReport, 16)
	pos := uint32(0x198CA0)
	read := func() int32 {
		defer (func() { pos = pos + 4 })()

		return s.readInt32(pos)
	}

	for i := 0; i < len(reports); i = i + 1 {
		r := &reports[i]

		r.RideConstruction = read()
		r.RideOperation = read()
		r.LandPurchase = read()
		r.Landscaping = read()
		r.ParkTickets = read()
		r.RideTickets = read()
		r.ShopSales = read()
		r.ShopStock = read()
		r.FoodSales = read()
		r.FoodStock = read()
		r.StaffWages = read()
		r.Marketing = read()
		r.Research = read()
		r.LoanInterest = read()
	}

	return reports
}
