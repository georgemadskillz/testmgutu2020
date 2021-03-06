package ioctrl

import (
	"bufio"
	"flydb/datamdl"
	"fmt"
	"os"
	"strings"
)

const (
	// FdbFly is a special fdb key for Flight table
	FdbFly = iota
	// FdbAir is a special fdb key for Airport table
	FdbAir = iota
	// FdbPrc is a special fdb key for Price table
	FdbPrc = iota
	// FdbAmount is amount of fdb keys
	FdbAmount = iota
)

// FlyDbIO is a common type for I/O actions
type FlyDbIO struct {
	Db    datamdl.FlyDb
	files []string
}

// Init initializes io fr FlyDB
func (io *FlyDbIO) Init() {

	io.Db.Init(100)

	io.files = make([]string, 3, 3)
	io.files[FdbFly] = "database/flights.fdb"
	io.files[FdbAir] = "database/airports.fdb"
	io.files[FdbPrc] = "database/prices.fdb"

	for key := FdbFly; key < FdbAmount; key++ {
		io.LoadFile(key)
	}
}

// GetRange is
func (io *FlyDbIO) GetRange(key int, fromIndex, toIndex int) []interface{} {
	if key >= FdbAmount {
		return nil
	}

	switch key {
	case FdbFly:

		flights := make([]datamdl.Flight, 0)
		for i := fromIndex; i <= toIndex; i++ {
			flight, err := io.Db.GetRow(key, i)

			if err != nil {
				flights = append(flights, datamdl.Flight{})
			} else {
				f := flight.(datamdl.Flight)
				flights = append(flights, f)
			}
		}

		result := make([]interface{}, 0)
		for i := range flights {
			result = append(result, flights[i])
		}

		return result

	case FdbAir:

		airports := make([]datamdl.Airport, 0)
		for i := fromIndex; i <= toIndex; i++ {
			airport, err := io.Db.GetRow(key, i)

			if err != nil {
				airports = append(airports, datamdl.Airport{})
			} else {
				a := airport.(datamdl.Airport)
				airports = append(airports, a)
			}

		}

		result := make([]interface{}, 0)
		for i := range airports {
			result = append(result, airports[i])
		}

		return result

	case FdbPrc:
		prices := make([]datamdl.Price, 0)
		for i := fromIndex; i <= toIndex; i++ {

			price, err := io.Db.GetRow(key, i)
			if err != nil {
				prices = append(prices, datamdl.Price{})
			} else {
				p := price.(datamdl.Price)
				prices = append(prices, p)
			}

		}

		result := make([]interface{}, 0)
		for pr := range prices {
			result = append(result, prices[pr])
		}

		return result
	}

	return nil
}

// LoadFile is
func (io *FlyDbIO) LoadFile(key int) error {

	if key >= FdbAmount {
		return fmt.Errorf("unknown key=%v", key)
	}

	path := io.files[key]

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		switch key {
		case FdbFly:
			flight, _ := parseFdbRow(key, line)
			io.Db.AppendRow(flight.(datamdl.Flight))
		case FdbAir:
			air, _ := parseFdbRow(key, line)
			io.Db.AppendRow(air.(datamdl.Airport))
		case FdbPrc:
			price, _ := parseFdbRow(key, line)
			io.Db.AppendRow(price.(datamdl.Price))
		}
	}

	return nil
}

func parseFdbRow(key int, line string) (interface{}, error) {
	if key >= FdbAmount {
		return datamdl.Flight{}, fmt.Errorf("unknown key=%v", key)
	}

	tokens := strings.Split(line, ",")

	switch key {
	case FdbFly:
		flight := datamdl.Flight{}
		flight.TimeFrom = tokens[0]
		flight.FlightFrom = tokens[1]
		flight.FlightTo = tokens[2]
		flight.TimeTo = tokens[3]
		return flight, nil
	case FdbAir:
		airport := datamdl.Airport{}
		airport.AirID = tokens[0]
		airport.AirCity = tokens[1]
		airport.AirName = tokens[2]
		return airport, nil
	case FdbPrc:
		price := datamdl.Price{}
		price.FlightID = tokens[0]
		price.Seat = tokens[1]
		price.Price = tokens[2]
		return price, nil
	}

	return nil, fmt.Errorf("unknown key=%v", key)
}
