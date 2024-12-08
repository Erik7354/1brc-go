package run_1

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const maxCityCount = 10_000

type city struct {
	Min   float64
	Max   float64
	Sum   float64
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
		line := scanner.Text()

		split := strings.Split(line, ";")
		name := split[0]
		temp, _ := strconv.ParseFloat(split[1], 64)

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
			ceilPrecision1(c.Min),
			ceilPrecision1(c.Sum/float64(c.Count)),
			ceilPrecision1(c.Max),
		)
		if i+1 < len(keys) {
			_, _ = fmt.Fprint(w, ", ")
		}
	}
	_, _ = fmt.Fprint(w, "}\n")
}
