/*
	Implementation ddccontrol for GO

	Some hint:
	- C.struct_monitorlist
		- contains information about detected monitor:
			- filename: path to i2c dev

		- is first monitor info returned from C.ddcci_probe()
		- to get(need assign) next use .next
	- C.struct_monitor
		- opened monitor (C.ddcci_open)
		- have some information about monitor:
			- pnpid: eg: "DELD072"
			- db (more specific information):
				-name: eg: "VESA standard monitor"

 */
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

	list []*MonitorInfo
	count int
}

type MonitorInfo struct {
	Id int
	Name string
	PnPid string
	monitor *C.struct_monitor
	monitorList *C.struct_monitorlist
	BrightnessLevel int
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
	ddcci.openSupportedMonitors()


	return ddcci, err
}

func probeDDCci() (*C.struct_monitorlist, error) {
	monitorList := C.ddcci_probe()
	if monitorList == nil {
		err := fmt.Errorf("monitor list is empty, is mod 'i2c-dev' loaded? ")
		return monitorList, err
	}
	return monitorList, nil
}

func makeDDCci(monitorList *C.struct_monitorlist) DDCci {
	return DDCci{
		monitorList:monitorList,
		list:make([]*MonitorInfo, 0),
		count: 0,
	}
}

func initDDCci() error {
	res := int(C.ddcci_init(nil))
	if res == 0 {
		return fmt.Errorf("ddcci initialize error")
	}
	return nil
}

func newMonitorInfo(mon *C.struct_monitorlist) MonitorInfo {
	return MonitorInfo{
		monitorList:mon,
		Name: C.GoString(mon.name),
	}
}

func (ddcci *DDCci)detectSupportedMonitors() {
	current := ddcci.monitorList
	for {
		if current == nil {
			break
		}
		printInfo(current)
		if current.supported == 1 {
			info := newMonitorInfo(current)
			ddcci.list = append(ddcci.list, &info)
			ddcci.count++
		}
		current = current.next
	}
}

func (ddcci *DDCci)openSupportedMonitors() {
	for _, mon := range ddcci.list {
		mon.openMonitor()
	}
}

func (info *MonitorInfo)openMonitor() error {

	if info.monitor != nil {
		return nil
	}

	fileName := info.monitorList.filename
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
	info.monitor = &mon
	info.Name = monitorName
	info.PnPid = pnpid

	return nil
}

func (ddcci *DDCci) MonitorList() []*MonitorInfo{
	return ddcci.list
}

func (ddcci *DDCci) DefaultMonitor() (*MonitorInfo, error){
	info := ddcci.list[0]
	return info, nil
}


func (info *MonitorInfo) MonitorName() string{
	return info.PnPid
}
func (info *MonitorInfo) MonitorFullName() string{
	return fmt.Sprintf("%v [%v]\n", info.PnPid, info.Name)
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

func (info *MonitorInfo) SetBrightness(value int8) error {
	if info == nil || info.monitor == nil {
		return fmt.Errorf("monitor closed. Please open first")
	}
	var cc C.char = 0x10
	delay := C.find_write_delay(info.monitor, cc)
	cval := C.ushort(value)
	C.ddcci_writectrl(info.monitor, 0x10, cval, delay)
	return nil
}
