package run_9

import (
	"bytes"
	"testing"
)

var longSlice = []byte("igButeboJuršinciKoaniImdinaNova VasDestrnikVarvarinSkopunGornji PetrovciRibnicaKon TumŠavnikPodl;11.5")
var shortSlice = []byte("Cabo San Lucas;14.9")
var testSlice = longSlice

func BenchmarkBytesIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = bytes.IndexByte(testSlice, ';')
	}
}

func BenchmarkCustomIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = indexByte(testSlice, ';')
	}
}
