package main

import (
	"github.com/zurek87/go-ddccontrol/gui"
	"fmt"
)

func main() {
	defer onClose()

	gddcci := gui.NewGui()
	gddcci.Show()
	gddcci.Main()
}

func onClose(){
	// check is no panic message
	if r := recover(); r != nil {
		msg := fmt.Sprintf("Error in Applicaton:\n %v", r)
		gui.DialogError(msg)
	}
}