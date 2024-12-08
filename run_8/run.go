package run_8

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
)

const (
	B  int = 1
	KB     = B << 10
	MB     = KB << 10
)

// maxLineLength does not need to be exact just > the longest possible line
const maxLineLength = 110
const maxStationCount = 10_000
const chunkSize = 16 * MB

// ==================================================================================== //
// Run
// ==================================================================================== //
type stationData struct {
	Min   int
	Max   int
	Sum   int
	Count uint
}

func Entrypoint(w io.Writer, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stations := make(map[string]*stationData, maxStationCount)

	scanner := newStationScanner(file)

	for scanner.Next() {
		name, temp := scanner.Line()

		if c, ok := stations[name]; ok { // update stationData
			c.Max = max(c.Max, temp)
			c.Min = min(c.Min, temp)
			c.Sum += temp
			c.Count++
		} else { // add stationData
			stations[name] = &stationData{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}
	}

	printCities(w, stations)
}

func ceilPrecision1(val float64) float64 {
	return math.Ceil(val*10) / 10
}

func printCities(w io.Writer, cities map[string]*stationData) {
	keys := make([]string, 0, len(cities))
	for k, _ := range cities {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	_, _ = fmt.Fprint(w, "{")
	for i, key := range keys {
		c := cities[key]
		_, _ = fmt.Fprintf(w, "%s=%.1f/%.1f/%.1f",
			key,
			ceilPrecision1(float64(c.Min)/10),
			ceilPrecision1(float64(c.Sum)/10/float64(c.Count)),
			ceilPrecision1(float64(c.Max)/10),
		)
		if i+1 < len(keys) {
			_, _ = fmt.Fprint(w, ", ")
		}
	}
	_, _ = fmt.Fprint(w, "}\n")
}
