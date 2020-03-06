package main

import (
	"xlang/testing"
)

func main() {
	testing.TestParser(`parser(x, y, z) * (6 + 3)`)
}
