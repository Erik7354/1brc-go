package run_7

import (
	"bufio"
	"os"
	"testing"
)

// ==================================================================================== //
// File Reading & Processing Comparison
// ==================================================================================== //

func Test_StationScanner(t *testing.T) {
	f, err := os.Open("../samples/measurements-10.txt")
	if err != nil {
		panic(err)
	}

	sc := NewStationScanner(f)
	for sc.Next() {
		name, temp := sc.Line()
		t.Logf("%s -> %d", name, temp)
	}
}

// The comparison to just file.Read and bufio.Scanner is totally unfair since StationScanner already processes the values which takes most (or at least a lot of) the time
// StationScanner.Line() and StationScanner.intTemp() could be seperated and just wrapped by StationScanner so the other benchmarks could use them.
func Benchmark_StationScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("../measurements_10mio.txt")
		if err != nil {
			panic(err)
		}

		sc := NewStationScanner(f)
		for sc.Next() {
			_, _ = sc.Line()
		}
	}
}

// ==================================================================================== //
// Raw File Reading Comparison
// ==================================================================================== //

func Benchmark_Scanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("../measurements_10mio.txt")
		if err != nil {
			panic(err)
		}

		chunk := make([]byte, 16*MB)

		scanner := bufio.NewScanner(f)
		scanner.Buffer(chunk, 16*MB)

		for scanner.Scan() {
			_ = scanner.Bytes()
		}
	}
}

func Benchmark_FileRead_16MB(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, err := os.Open("../measurements_10mio.txt")
		if err != nil {
			panic(err)
		}

		chunk := make([]byte, 16*MB)

		for {
			_, err = f.Read(chunk)
			if err != nil {
				break
			}
		}
	}
}

func Benchmark_FileRead_1KB(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, err := os.Open("../measurements_10mio.txt")
		if err != nil {
			panic(err)
		}

		chunk := make([]byte, 1*KB)

		for {
			_, err = f.Read(chunk)
			if err != nil {
				break
			}
		}
	}
}
