package gui

import (
	"github.com/zurek87/go-ddccontrol/goddcci"
	"github.com/zurek87/go-gtk3/glib"
	"github.com/zurek87/go-gtk3/gtk3"
	"github.com/zurek87/go-gtk3/gdk3"
	"github.com/zurek87/go-gtk3/gdkpixbuf"
	"unsafe"
	"fmt"
)

type GDDCci struct {
	ddcci *goddcci.DDCci
	visible bool
	icon *gtk3.StatusIcon
	menu *gtk3.Menu
	list []*GuiMonitorInfo
	selected int
}

func NewGui() GDDCci {

	ddcci, err := goddcci.InitDDCci()
	if err != nil {
		panic(err)
	}
	defaultMon, err := ddcci.DefaultMonitor()
	if err != nil {
		panic(err)
	}
	ddcciMonitorList := ddcci.MonitorList()

	list := guiMonitorList(ddcciMonitorList)
	return GDDCci{
		ddcci: &ddcci,
		selected: defaultMon.Id,
		list: list,
	}
}

func (gddcci *GDDCci) Main() {
	gtk3.Main()
}

func (gddcci *GDDCci) Show() {
	if gddcci.visible {
		return
	}
	gddcci.initGTK()
	gddcci.initIcon()
	gddcci.initMenu()
}

func (gddcci *GDDCci) initGTK() {
	gtk3.Init(nil)
}

func (gddcci *GDDCci) initIcon() {
	appName := fmt.Sprintf("Brightness on %v", gddcci.monitorName())
	glib.SetApplicationName(appName)

	pixBuff, _ := gdkpixbuf.NewPixbufFromBytes(GuiIcon)
	gddcci.icon = gtk3.NewStatusIconFromPixbuf(pixBuff)
	gddcci.icon.Connect("popup-menu", func(cbx *glib.CallbackContext) {
		gddcci.menu.Popup(nil, nil, gtk3.StatusIconPositionMenu, gddcci.icon, uint(cbx.Args(0)), uint32(cbx.Args(1)))
	})
	gddcci.icon.Connect("scroll-event", func(cbx *glib.CallbackContext) {
		gddcci.onScroll(cbx)
	})
}

func (gddcci *GDDCci) initMenu() {
	gddcci.menu = gtk3.NewMenu()
	gddcci.addMenuItems()
	gddcci.menu.Append(gtk3.NewSeparatorMenuItem())
	gddcci.menu.Append(closeMenuItem())

	gddcci.menu.ShowAll()
}

func (gddcci *GDDCci)onScroll(cbx *glib.CallbackContext) {
	arg := cbx.Args(0)
	event := *(**gdk3.EventScroll)(unsafe.Pointer(&arg))
	var stateUp uint   = 0x200000000
	var stateDown uint = 0x300000000
	var deltaVal int
	if event.State & stateUp == stateUp && event.State & stateDown != stateDown {
		deltaVal = 3
	}
	if event.State & stateDown == stateDown {
       deltaVal = -3
	}
	gddcci.ChangeBrightness(deltaVal)
}

func closeMenuItem() *gtk3.MenuItem {
	item := gtk3.NewMenuItemWithLabel("Exit!")
	item.Connect("activate", func() {
		gtk3.MainQuit()
	})
	return item
}

func (gddcci *GDDCci) addMenuItems() {
	// add select monitor
	group := glib.GSListAlloc()
	for _, info := range gddcci.list{
		item , gr := createChoseMonitorMenuItem(info, group, gddcci)
		gddcci.menu.Append(item)
		group = gr // go-gtk bug
	}

	// add config monitor
	for _, info := range gddcci.list{
		item := createConfigMonitorMenuItem(info)
		gddcci.menu.Append(item)
	}
}


func (gddcci *GDDCci)  monitorName() string {
	monInfo := gddcci.list[gddcci.selected]
	if monInfo.info == nil {
		return "No monitor found"
	}
	return monInfo.info.MonitorName()
}

func (gddcci *GDDCci)  ChangeBrightness(deltaVal int) {
	monInfo := gddcci.list[gddcci.selected]
	if monInfo.info == nil {
		panic(guiNoMonitorError{"monitor is <nil> can't change brightness"})
	}
	monInfo.LevelChan <- deltaVal
}

func (gddcci *GDDCci)  SetBrightness(value uint16) error {
	monInfo := gddcci.list[gddcci.selected]
	if monInfo.info == nil {
		return guiNoMonitorError{"monitor is <nil> can't set brightness"}
	}
	return monInfo.info.SetBrightness(value)
}