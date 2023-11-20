package main

import (
	"fmt"
	"strconv"
	"time"

	juicer "github.com/dankobg/juicer/engine"
)

func main() {
	now := time.Now()

	for i := 0; i < 64; i++ {
		m := juicer.FindMagicNum(i, true)
		fmt.Printf("bishop: %v 0x%v\n", i, strconv.FormatUint(uint64(m), 16))
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()

	for i := 0; i < 64; i++ {
		m := juicer.FindMagicNum(i, false)
		fmt.Printf("bishop: %v 0x%v\n", i, strconv.FormatUint(uint64(m), 16))
	}

	fmt.Println(time.Since(now))
}
