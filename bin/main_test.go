package main

import (
	"os"
	"testing"
)

func TestMainBin(t *testing.T) {
	os.Args = append(os.Args, "../tests/car.ldr")
	main()
}
