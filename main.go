package main

import (
	"homevision/cmd/tui"

	"github.com/common-nighthawk/go-figure"
)

func main() {

	myFigure := figure.NewFigure("Homevision", "", true)
	myFigure.Print()

	tui.StartTUI()
}
