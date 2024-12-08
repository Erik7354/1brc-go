package run_6

import "testing"

// ==================================================================================== //
// Loop vs Switch vs If
// ==================================================================================== //

// 2.5 ns/op
func intTempLoop(bs []byte) int {
	neg := bs[0] == '-'
	if neg {
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

func BenchmarkIntTempLoop(b *testing.B) {
	num := []byte("-77.7")

	for i := 0; i < b.N; i++ {
		_ = intTempLoop(num)
	}
}

// 2.1 ns/op
func intTempSwitch(bs []byte) int {
	neg := bs[0] == '-'
	if neg {
		bs = bs[1:]
	}

	var res int

	switch {
	case bs[1] == '.': // X.X
		res = int(bs[0]-48)*10 + int(bs[2]-48)
	case bs[2] == '.': // XX.X
		res = int(bs[0]-48)*100 + int(bs[1]-48)*10 + int(bs[3]-48)
	}

	if neg {
		return -res
	}
	return res
}

func BenchmarkIntTempSwitch(b *testing.B) {
	num := []byte("-77.7")

	for i := 0; i < b.N; i++ {
		_ = intTempSwitch(num)
	}
}

// 2.1 ns/op
func intTempIfElseIf(bs []byte) int {
	neg := bs[0] == '-'
	if neg {
		bs = bs[1:]
	}

	var res int
	if bs[1] == '.' {
		res = int(bs[0]-48)*10 + int(bs[2]-48)
	} else if bs[2] == '.' {
		res = int(bs[0]-48)*100 + int(bs[1]-48)*10 + int(bs[3]-48)
	}

	if neg {
		return -res
	}
	return res
}

func BenchmarkIntTempIfElseIf(b *testing.B) {
	num := []byte("-77.7")

	for i := 0; i < b.N; i++ {
		_ = intTempIfElseIf(num)
	}
}

// 0.7 ns/op
func intTempIfElse(bs []byte) int {
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

func BenchmarkIntTempIfElse(b *testing.B) {
	num := []byte("-77.7")

	for i := 0; i < b.N; i++ {
		_ = intTempIfElse(num)
	}
}
