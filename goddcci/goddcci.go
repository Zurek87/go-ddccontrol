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
	// structures from ddcci
	monitorList *C.struct_monitorlist
	supported []*C.struct_monitorlist
	selected *C.struct_monitorlist
	monitor *C.struct_monitor
	// go friendly ;)
	list []MonitorInfo
	// selected monitor info:
	monitorName string
	pnpid string
}

type MonitorInfo struct {
	Id int
	Name string
	PnPid string
	monitor *C.struct_monitor
}

func InitDDCci() (DDCci, error) {
	err := initDDCci()
	if err != nil {
		return DDCci{}, err
	}
	fmt.Println("Initialized DDCci")

	monitorList, err := probeDDCci()
	if err != nil {
		return DDCci{}, err
	}
	ddcci := makeDDCci(monitorList)
	ddcci.detectSupportedMonitors()
	err = ddcci.openMonitor()

	return ddcci, err
}

func probeDDCci() (*C.struct_monitorlist, error) {
	monitorList := C.ddcci_probe()
	if monitorList == nil {
		err := fmt.Errorf("monitor list is empty, is mod 'i2c-dev' loaded? ")
		return monitorList, err
	}
	fmt.Println("monitor list: ", monitorList)
	return monitorList, nil
}

func makeDDCci(monitorList *C.struct_monitorlist) DDCci {
	return DDCci{
		monitorList:monitorList,
		supported:make([]*C.struct_monitorlist, 0),
	}
}

func initDDCci() error {
	res := int(C.ddcci_init(nil))
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
		C.ddcci_open(&mon, fileName, 0)
		monitorName := ""
		pnpid := "UnKnow"
		if mon.db != nil {
			name := C.xmlCharToChar(mon.db.name)
			monitorName = C.GoString(name)
			pnpid = C.GoString(&mon.pnpid[0])
		}
		fmt.Printf("Opened monitor: %v [%v]\n", pnpid, monitorName)
		ddcci.monitor = &mon
		ddcci.monitorName = monitorName
		ddcci.pnpid = pnpid
		return nil
	}
	return fmt.Errorf("DDCCi no supported monitor found")
}



func (ddcci *DDCci) MonitorList() []MonitorInfo{
	return ddcci.list
}

func (ddcci *DDCci) MonitorName() string{
	return ddcci.pnpid
}
func (ddcci *DDCci) MonitorFullName() string{
	return fmt.Sprintf("%v [%v]\n", ddcci.pnpid, ddcci.monitorName)
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
