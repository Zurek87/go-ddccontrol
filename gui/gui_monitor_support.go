package gui

import (
	"github.com/zurek87/go-ddccontrol/goddcci"
	"fmt"
	"github.com/zurek87/go-gtk3/gtk3"
	"github.com/zurek87/go-gtk3/glib"
	"time"
)

type GuiMonitorInfo struct{
	info *goddcci.MonitorInfo
	run bool
	changeChan chan int
	LevelChan chan<- int
}


/*
Can run only in separate thread
 */
func (info *GuiMonitorInfo)listen() {
	for {
		level := info.getNewLevel()
		info.updateMonitor(level)
	}
}

func (info *GuiMonitorInfo)getNewLevel() uint16 {
	value := <- info.changeChan
	lenInChan := len(info.changeChan)
	if lenInChan > 0 {
		for i := 0; i < lenInChan; i++{
			val := <- info.changeChan
			value += val
		}
	}

	level := info.info.Level()

	if value + int(level) < 0{
		level = 0
	} else {
		level += uint16(value)
	}

	if level > info.info.MaxLevel() {
		level = info.info.MaxLevel()
	}
	return level
}

func (info *GuiMonitorInfo)updateMonitor(level uint16) {
	if info.info.Level() != level {
		info.info.SetBrightness(level)
	}
	time.Sleep(50*time.Millisecond)
}
func (info *GuiMonitorInfo)Run() {
	if info.run {
		return
	}
	go info.listen()
	info.run = true
}

func (gddcci *GDDCci) SelectMonitor(id int) error {
	if len(gddcci.list) <= id {
		return fmt.Errorf("monitor id out of range")
	}
	gddcci.selected = id
	return nil
}

func createChoseMonitorMenuItem(info *GuiMonitorInfo, group *glib.SList, gddcci *GDDCci) (*gtk3.RadioMenuItem, *glib.SList) {
	label := fmt.Sprintf("Default: %v", info.info.PnPid)
	item := gtk3.NewRadioMenuItemWithLabel(group, label)
	item.Connect("activate", func() {
		gddcci.SelectMonitor(info.info.Id)
	})
	gr := item.GetGroup()
	return item, gr
}

func guiMonitorList(ddcciMonitorList []*goddcci.MonitorInfo) []*GuiMonitorInfo {
	lenList := len(ddcciMonitorList)
	guiList := make([]*GuiMonitorInfo, lenList)
	for i, info := range ddcciMonitorList {
		ch := make(chan int, 32)
		guiInfo := GuiMonitorInfo{
			info: info,
			changeChan:ch,
			LevelChan:ch,
		}
		guiInfo.Run()
		guiList[i] = &guiInfo

	}
	return guiList
}