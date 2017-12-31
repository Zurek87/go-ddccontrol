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
	monitorlist *C.struct_monitorlist
}

func InitDDCci() (DDCci, error) {
	err := initDDCci()
	if err != nil {
		return DDCci{}, err
	}
	fmt.Println("Initialized DDCci")
	monitorList := C.ddcci_probe()
	return DDCci{monitorlist:monitorList}, nil
}

func initDDCci() error {
	res := int(C.ddcci_init(nil))
	if res == 0 {
		return fmt.Errorf("ddcci initialize error")
	}
	return nil
}

func Probe() {
	k := C.int(4)

	C.ddcci_verbosity(k)
}