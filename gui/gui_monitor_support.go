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
	labelItem *gtk3.MenuItem
}

// for fast creation of menu items
type menuValueName struct {
	Name string
	Value int
	IsDelta bool
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
		if info.labelItem != nil {
			label := fmt.Sprintf("%v [ %v ]", info.info.PnPid, level)
			info.labelItem.SetLabel(label)
		}
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

func createConfigMonitorMenuItem(info *GuiMonitorInfo) (*gtk3.MenuItem) {
	label := fmt.Sprintf("%v [ %v ]", info.info.PnPid, info.info.Level())

	item := gtk3.NewMenuItemWithLabel(label)
	submenu := gtk3.NewMenu()

	itemLabel := gtk3.NewMenuItemWithLabel(info.info.PnPid)
	itemLabel.SetSensitive(false)

	submenu.Append(itemLabel)
	submenu.Append(gtk3.NewSeparatorMenuItem())
	info.labelItem = item

	constValues := listOfMenuValueName(info)
	for _, valInfo := range constValues{
		if valInfo.Name == "---" {
			submenu.Append(gtk3.NewSeparatorMenuItem())
		} else {
			subItem := createSetBrightnessMenuItem(valInfo)
			info.bindOnClick(subItem, valInfo)
			submenu.Append(subItem)
		}

	}
	item.SetSubmenu(submenu)
	return item
}

func listOfMenuValueName(info *GuiMonitorInfo) []menuValueName{
	list := make([]menuValueName, 9)
	val := int(info.info.MaxLevel())
	list[0] = menuValueName{"Max", val, false}
	list[1] = menuValueName{"Middle", val/2, false}
	list[2] = menuValueName{"Min", 0, false}
	list[3] = menuValueName{"---", 0, false} // separator
	list[4] = menuValueName{"+", 5, true}
	list[5] = menuValueName{"-", -5, true}
	list[6] = menuValueName{"---", 0, false} // separator
	list[7] = menuValueName{"Video day", 14, false}
	list[8] = menuValueName{"Video night", 8, false}
	return list
}

func (info *GuiMonitorInfo)bindOnClick(item *gtk3.MenuItem, menuVal menuValueName) {

	item.Connect("activate", func() {
		if menuVal.IsDelta {
			info.LevelChan <- menuVal.Value
		} else {
			info.updateMonitor(uint16(menuVal.Value))
		}
	})
}

func createSetBrightnessMenuItem(valInfo menuValueName) (*gtk3.MenuItem){

	label := fmt.Sprintf("%v (%v)", valInfo.Name, valInfo.Value)
	if valInfo.Name == "-" || valInfo.Name == "+" {
		val := valInfo.Value
		if val < 0 {
			val = -val
		}
		label = fmt.Sprintf("%v %v", valInfo.Name, val)
	}
	item := gtk3.NewMenuItemWithLabel(label)

	return item
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