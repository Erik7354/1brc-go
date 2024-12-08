package main

import (
	"1brc/run_9"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()

	run_9.Entrypoint(os.Stdout, "measurements_1b.txt")

	elapsed := time.Since(start)
	fmt.Printf("took %s\n", elapsed)
}
