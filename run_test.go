package main

import (
	"1brc/concurrent_1"
	"1brc/run_1"
	"1brc/run_2"
	"1brc/run_3"
	"1brc/run_4"
	"1brc/run_5"
	"1brc/run_6"
	"1brc/run_7"
	"1brc/run_8"
	"1brc/run_9"
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	B  int = 1
	KB     = B << 10
	MB     = KB << 10
	GB     = MB << 10
)

// ==================================================================================== //
// Benchmark
// ==================================================================================== //
// go test -run=XXX -benchmem -v -bench=BenchmarkRun9

func BenchmarkRun9(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_9.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_8.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_7.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_6.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_5.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_4.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_3.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_2.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkRun1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run_1.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

func BenchmarkConcurrent1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concurrent_1.Entrypoint(os.Stdout, "measurements_1b.txt")
	}
}

// ==================================================================================== //
// Test
// ==================================================================================== //

func TestAll(t *testing.T) {
	funcsToTest := []func(io.Writer, string){
		concurrent_1.Entrypoint,
		run_1.Entrypoint,
		run_2.Entrypoint,
		run_3.Entrypoint,
		run_4.Entrypoint,
		run_5.Entrypoint,
		run_6.Entrypoint,
		run_7.Entrypoint,
		run_8.Entrypoint,
		run_9.Entrypoint,
	}

	matches, _ := filepath.Glob("samples/*.txt")
	t.Logf("testing with files: %v \n", matches)

	for _, match := range matches {
		t.Logf("testing file: %s", match)

		expectedPath := strings.TrimSuffix(match, ".txt") + ".out"
		expectedb, _ := os.ReadFile(expectedPath)
		expected := md5.Sum(expectedb)

		for i, fun := range funcsToTest {
			t.Logf("\t testing run: %d", i+1)

			var buf bytes.Buffer
			fun(&buf, match)

			res := md5.Sum(buf.Bytes())

			if expected != res {
				t.Logf("\t\t run_%d failed", i+1)
				t.Logf("\t\t produced hash %x expected %x \n", res, expected)
				t.Logf("\t\t produced %s \n", buf.String())
				t.Fail()
			}
		}

	}
}
