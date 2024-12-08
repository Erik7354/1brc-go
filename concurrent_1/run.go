package concurrent_1

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"sync"
)

const nMaxCities = 10_000

const (
	B  int = 1
	KB     = B << 10
	MB     = KB << 10
	GB     = MB << 10
)

// tuning parameters
const chunkSize = 16 * MB
const nConsumer = 25
const nInBuffer = 25

type city struct {
	Min   int
	Max   int
	Sum   int
	Count uint
}

func Entrypoint(w io.Writer, filepath string) {
	inChans := make([]chan []byte, nConsumer)
	outChans := make([]chan map[string]*city, nConsumer)

	var wg sync.WaitGroup
	wg.Add(nConsumer)

	// Create workers
	for i := range nConsumer {
		input := make(chan []byte, nInBuffer)
		output := make(chan map[string]*city, 1)

		go consumer(input, output, &wg)

		inChans[i] = input
		outChans[i] = output
	}

	// read file
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	chunk := make([]byte, chunkSize)
	leftover := make([]byte, 110) // should be able to hold the longest possible line
	leftoverSize := 0
	iConsumer := 0
	var n int
	for {
		n, err = file.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		i := bytes.LastIndexByte(chunk[:n], '\n')
		if i == -1 {
			i = 0
		}

		lines := make([]byte, i+leftoverSize)
		copy(lines, leftover[:leftoverSize])
		copy(lines[leftoverSize:], chunk[:i])

		copy(leftover, chunk[i+1:n])
		leftoverSize = n - i - 1

		inChans[iConsumer%nConsumer] <- lines
		iConsumer++
	}

	// stop consumers
	for i := range nConsumer {
		close(inChans[i])
	}

	// wait for pending consumers
	wg.Wait()
	for i := range nConsumer {
		close(outChans[i])
	}

	// collect results
	cities := make(map[string]*city, nMaxCities)
	for _, outChan := range outChans {
		for name, v1 := range <-outChan {
			if v2, ok := cities[name]; ok {
				v2.Min = min(v1.Min, v2.Min)
				v2.Max = max(v1.Max, v2.Max)
				v2.Sum = v1.Sum + v2.Sum
				v2.Count = v1.Count + v2.Count
			} else {
				cities[name] = &city{
					Min:   v1.Min,
					Max:   v1.Max,
					Sum:   v1.Sum,
					Count: v1.Count,
				}
			}
		}
	}

	// print
	printCities(w, cities)
}

func consumer(in chan []byte, out chan map[string]*city, wg *sync.WaitGroup) {
	defer wg.Done()
	cities := make(map[string]*city, 100)

	for lines := range in {
		var offset int
		for offset < len(lines) {
			used, name, temp := processLine(lines[offset:])
			if used == 0 {
				break
			}
			offset += used

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
	}

	out <- cities
}

// processLines takes whatever amount of lines and processes the first one.
// [used] gives the byte length of the first row.
func processLine(lines []byte) (used int, city string, temp int) {
	l := bytes.IndexByte(lines, '\n')
	if l == -1 {
		l = len(lines)
	}

	switch {
	case lines[l-4] == ';': // 1.2
		return l + 1, string(lines[:l-4]), intTemp(lines[l-3 : l])
	case lines[l-5] == ';': // 12.3 or -1.2
		return l + 1, string(lines[:l-5]), intTemp(lines[l-4 : l])
	case lines[l-6] == ';': // -12.3
		return l + 1, string(lines[:l-6]), intTemp(lines[l-5 : l])
	default:
		panic("unknown format")
	}
}

// intTemp converts second part of a line to int.
// "-77.7" => -777
// "77.7" => 777
func intTemp(bs []byte) int {
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
