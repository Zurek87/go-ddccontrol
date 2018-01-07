package goddcci

import (
	"fmt"
)

/*
#cgo pkg-config: --static libxml-2.0
#cgo pkg-config: ddccontrol
#include "goddcci.go.h"
#include <ddccontrol/ddcci.h>

*/
import "C"


type DDCci struct {
	monitorList *C.struct_monitorlist
	selected *C.struct_monitorlist
	monitor *C.struct_monitor
	supported []*C.struct_monitorlist
	monitorName string
}

func InitDDCci() (DDCci, error) {
	err := initDDCci()
	if err != nil {
		return DDCci{}, err
	}
	fmt.Println("Initialized DDCci")
	monitorList := C.ddcci_probe()
	ddcci := DDCci{
		monitorList:monitorList,
		supported:make([]*C.struct_monitorlist, 0),
		}
	ddcci.detectSupportedMonitors()
	err = ddcci.openMonitor()

	return ddcci, err
}

func initDDCci() error {
	res := int(C.ddcci_init(nil))
	fmt.Println("Ten teges odpowiedz z init: ", res)
	if res == 0 {
		return fmt.Errorf("ddcci initialize error")
	}
	return nil
}


func (ddcci *DDCci)detectSupportedMonitors() {
	current := ddcci.monitorList
	for {
		if current == nil {
			break
		}
		printInfo(current)
		if current.supported == 1 {
			ddcci.supported = append(ddcci.supported, current)
			if ddcci.selected == nil {
				ddcci.selected = current
			}
		}
		current = current.next
	}

}

func (ddcci *DDCci)openMonitor() error {
	if ddcci.selected != nil {
		fileName := ddcci.selected.filename
		var mon C.struct_monitor
		i := C.ddcci_open(&mon, fileName, 0)
		fmt.Println("Ten teges odpowiedz z open: ", i)
		monName := "UnKnow"
		pnpid := "UnKnow"
		if mon.db != nil {
			name := C.xmlCharToChar(mon.db.name)
			monName = C.GoString(name)
			pnpid = C.GoString(&mon.pnpid[0])
		}
		fmt.Printf("Opened monitor: %v [%v]\n", pnpid, monName)
		ddcci.monitor = &mon
		ddcci.monitorName = monName
		return nil
	}
	return fmt.Errorf("DDCCi no supported monitor found, is mod 'i2c-dev' loaded? ")
}

func printInfo(monList *C.struct_monitorlist) {
	fmt.Println("\nMonitor:")
	name := C.GoString(monList.name)
	fmt.Printf("Name: %v\n", name)
	fileName := C.GoString(monList.filename)
	fmt.Printf("I2C File name: %v\n", fileName)
	supported := "Yes"
	if monList.supported == 0 {
		supported = "No"
	}
	fmt.Printf("DDC/CI supported: %v\n", supported)
	input := "Analog"
	if monList.digital > 0 {
		input = "Digital"
	}
	fmt.Printf("Input type:: %v\n", input)
}

func (ddcci *DDCci)SetBrightness(value int8) {
	var cc C.char = 0x10
	delay := C.find_write_delay(ddcci.monitor, cc)
	cval := C.ushort(value)
	C.ddcci_writectrl(ddcci.monitor, 0x10, cval, delay)
}
