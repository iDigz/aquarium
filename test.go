package main

import (
	"fmt"
)

type pins []struct {
	start struct{
		hours int
		minutes int
	}
	stop struct{
		hours int
		minutes int
	}
}

var pips struct{
	pip [0]struct{
		start, 10
	}
}

func main()  {
	fmt.Println(pins)
}