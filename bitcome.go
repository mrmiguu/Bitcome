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
	Percent = 5.

	bitcome Bitcome
)

// Data maps points in time to dollar values.
type Data map[time.Time]float64

// Bitcome is the program's core structure.
type Bitcome struct {
	DataFile  string
	Pool      float64
	Threshold float64

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
}

// Run fires off the main program.
func Run() error {
	return bitcome.Run()
}

// Run fires off the main program.
func (b *Bitcome) Run() error {
	b.once.Do(b.init)

	buy, err := buy()
	if err != nil {
		return warn(err)
	}
	println("LTC:buy   $" + ftoa(buy))

	sell, err := sell()
	if err != nil {
		return warn(err)
	}
	println("LTC:sell  $" + ftoa(sell))

	net := sell - buy
	println("LTC:net   $" + ftoa(net))

	pct := percent(buy, sell)
	println("LTC:pct   " + ftoa(pct) + "%")

	// if pct > 10 {
	// 	b.dat[time.Now()] = sell
	// }

	return nil
}

// Save saves the program's collected data.
func Save() error {
	return bitcome.Save()
}

// Save saves the program's collected data.
func (b *Bitcome) Save() error {
	b.once.Do(b.init)
	return savedata(b.dat, b.DataFile)
}
