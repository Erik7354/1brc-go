package run_4

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
)

const maxCityCount = 10_000

type city struct {
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

	cities := make(map[string]*city, maxCityCount)

	for scanner.Scan() {
		line := scanner.Bytes()

		nameb, tempb, _ := bytes.Cut(line, []byte{';'})
		name := string(nameb)
		temp := lineToInt(tempb)

		if c, ok := cities[name]; ok { // update city
			c.Max = max(c.Max, temp)
			c.Min = min(c.Min, temp)
			c.Sum += temp
			c.Count++
		} else { // add city
			cities[name] = &city{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}
	}

	printCities(w, cities)
}

// lineToInt converts the second part of a line to int.
// "-77.7" => -777
// "77.7" => 777
func lineToInt(bs []byte) int {
	neg := false
	if bs[0] == '-' {
		neg = true
		bs = bs[1:]
	}

	var res int
	for _, b := range bs {
		if b == '.' {
			continue
		}

		res = res*10 + int(b-48) // numbers start at 48 in ascii
	}

	if neg {
		return -res
	}
	return res
}

func ceilPrecision1(val float64) float64 {
	return math.Ceil(val*10) / 10
}

func printCities(w io.Writer, cities map[string]*city) {
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
