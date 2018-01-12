package goddcci

import (
	"testing"
	"time"
)

func TestInitDDCci2(t *testing.T) {
	err := initDDCci()
	if err != nil{
		t.Error("Expected Success","Got error:", err)
	}
}


func TestDDCci(t *testing.T) {
	err := initDDCci()
	if err != nil {
		t.Error("Expected initialized success","Got error:", err)
	}

	monitorList, err := probeDDCci()
	if err != nil {
		t.Error("Expected monitorList, got :", monitorList, "with error:", err)
	}
	ddcci := makeDDCci(monitorList)
	ddcci.detectSupportedMonitors()
	ddcci.openSupportedMonitors()
	if ddcci.count == 0 {
		t.Error("Expected to find supported monitor")
	}

	defaultMonitor, err := ddcci.DefaultMonitor()
	if err != nil {
		t.Error("Expected get default monitor, got error:", err)
	}
	if defaultMonitor.Name == ""{
		t.Error("Expected to first monitor have name")
	}

	list := ddcci.MonitorList()
	if len(list) == 0 {
		t.Error("Empty monitor list")
	}
	// reading of current monitor value
	defaultVal, _, err := defaultMonitor.ReadCtrl(0x10)
	if err != nil {
		t.Error("Expected to get brightness value, got error:", err)
	}

	var testValue uint16 = 14
	if defaultVal == testValue {
		testValue++
	}
	err = defaultMonitor.SetBrightness(testValue)
	if err != nil {
		t.Error("Expected to set brightness, got error:", err)
	} else {
		time.Sleep(1 * time.Second)
		tmp, _, _ := defaultMonitor.ReadCtrl(0x10)
		if tmp != uint16(testValue){
			t.Error("Expected to confirm brightness, got", tmp)
		}
		defaultMonitor.SetBrightness(defaultVal)
	}


}