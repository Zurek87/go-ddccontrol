package main

import (
	"github.com/zurek87/go-ddccontrol/gui"
)

func main() {
	gddcci := gui.NewGui()
	gddcci.Show()

	gddcci.Main()

}