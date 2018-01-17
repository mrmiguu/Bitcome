package bitcome

import (
	"sync"
	"time"
)

var (
	// DataFile is the name of the file used to
	// store price points.
	DataFile = "data.json"

	// Pool is the fixed amount of money used
	// to calculate potential net profit.
	//
	// This acts as a cap.
	Pool = 2.00

	// Percent is the percentage difference from
	// bought to sold needed for an alert.
	//
	// This takes fees into account.
	Percent = 8.

	bitcome Bitcome
)

// Data is a list of prices at specific times.
type Data []Snapshot

// Snapshot is a price at specific time.
type Snapshot struct {
	USD  float64
	Time time.Time
}

// Bitcome is the program's core structure.
type Bitcome struct {
	DataFile string
	Pool     float64
	Percent  float64

	once sync.Once
	dat  Data
}

func (b *Bitcome) init() {
	var err error

	if len(b.DataFile) == 0 {
		b.DataFile = DataFile
	}
	if b.Pool == 0 {
		b.Pool = Pool
	}

	b.dat, err = appendfile(b.dat, b.DataFile)
	warn(err)
	println(b.dat)
}

// Passed checks current cost, returning whether
// the set percentage has been reached.
func Passed() (bool, error) {
	return bitcome.Passed()
}

// Passed checks current cost, returning whether
// the set percentage has been reached.
func (b *Bitcome) Passed() (bool, error) {
	b.once.Do(b.init)

	// safe copy
	dat := make(Data, len(b.dat))
	for i, snap := range b.dat {
		dat[i] = snap
	}

	buynow, err := buy()
	if err != nil {
		return false, warn(err)
	}
	now := time.Now()
	bought := buynow
	if len(dat) > 0 {
		bought = dat[len(dat)-1].USD
	}
	dat = append(dat, Snapshot{buynow, now})
	println("LTC:buy   " + ftousd(bought))

	sell, err := sell()
	if err != nil {
		return false, warn(err)
	}
	println("LTC:sell  " + ftousd(sell))

	fees := sell - buynow
	println("LTC:fees  " + ftousd(fees))

	pct := percent(bought, sell)
	println("LTC:pct   " + ftoa(pct) + "%")

	// safe mutation
	b.dat = dat

	return pct >= b.Percent, nil
}

// Save saves the program's collected data.
func Save() error {
	return bitcome.Save()
}

// Save saves the program's collected data.
func (b *Bitcome) Save() error {
	b.once.Do(b.init)
	println(b.dat)
	return savedata(b.dat, b.DataFile)
}
