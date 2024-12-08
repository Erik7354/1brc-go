package run_6

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
)

const maxCityCount = 10_000

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

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	stations := make(map[string]*stationData, maxCityCount)

	for scanner.Scan() {
		line := scanner.Bytes()

		name, tempb := splitLine(line)
		temp := lineToInt(tempb)

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

// splitLine splits a line into stationData name and temperature.
// "München;23.5" => "München", []byte("23.5")
// "München;-10.5" => "München", []byte("-10.5")
func splitLine(line []byte) (string, []byte) {
	l := len(line)

	switch {
	case line[l-4] == ';': // 1.2
		return string(line[:l-4]), line[l-3:]
	case line[l-5] == ';': // 12.3 or -1.2
		return string(line[:l-5]), line[l-4:]
	case line[l-6] == ';': // -12.3
		return string(line[:l-6]), line[l-5:]
	default:
		panic("unknown format")
	}
}

// lineToInt converts the second part of a line to int.
// "-77.7" => -777
// "77.7" => 777
func lineToInt(bs []byte) int {
	neg := bs[0] == '-'
	if neg {
		bs = bs[1:]
	}

	var res int
	if bs[1] == '.' {
		res = int(bs[0]-48)*10 + int(bs[2]-48)
	} else { // bs[2] == '.'
		res = int(bs[0]-48)*100 + int(bs[1]-48)*10 + int(bs[3]-48)
	}

	if neg {
		return -res
	}
	return res
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
