package gui

import (
	"github.com/zurek87/go-ddccontrol/goddcci"
	"github.com/zurek87/go-gtk3/glib"
	"github.com/zurek87/go-gtk3/gtk3"
	"github.com/zurek87/go-gtk3/gdk3"
	"unsafe"
	"fmt"
)

func DialogError(errorMessage string) {
	gtk3.Init(nil)
	fmt.Println("------------------------------------------")
	fmt.Println(errorMessage)
	fmt.Println("------------------------------------------")
	dialog := gtk3.NewMessageDialog(
		nil,
		gtk3.DIALOG_MODAL,
		gtk3.MESSAGE_ERROR,
		gtk3.BUTTONS_OK,
		errorMessage,
	)
	dialog.Response(func() {
		dialog.Destroy()
	})
	dialog.Run()
}

type GDDCci struct {
	ddcci *goddcci.DDCci
	visible bool
	level int8
	lastLevel int8
	icon *gtk3.StatusIcon
	menu *gtk3.Menu
	list []*goddcci.MonitorInfo
	selected *goddcci.MonitorInfo
}

func NewGui() GDDCci {

	ddcci, err := goddcci.InitDDCci()
	if err != nil {
		panic(err)
	}
	info, err := ddcci.DefaultMonitor()
	if err != nil {
		panic(err)
	}
	list := ddcci.MonitorList()

	return GDDCci{
		ddcci: &ddcci,
		lastLevel: 101,
		selected: info,
		list: list,
	}
}

func (gddcci *GDDCci) Show() {
	if gddcci.visible {
		return
	}
	gddcci.initGTK()
	gddcci.initIcon()
	gddcci.initMenu()

}


func (gddcci *GDDCci) SelectMonitor(id int) error {
	if len(gddcci.list) <= id {
		return fmt.Errorf("monitor id out of range")
	}
	gddcci.selected = gddcci.list[id]
	return nil
}

func (gddcci *GDDCci) initGTK() {
	gtk3.Init(nil)
}

func (gddcci *GDDCci) initIcon() {

	appName := fmt.Sprintf("Brightness on %v", gddcci.selected.MonitorName())
	glib.SetApplicationName(appName)


	gddcci.icon = gtk3.NewStatusIconFromFile("./gui/icons/brightness.png")
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

func (gddcci *GDDCci) addMenuItems() {
	// add select monitor
	group := glib.GSListAlloc()
	for _, info := range gddcci.list{
		item , gr := createChoseMonitorMenuItem(info, group, gddcci)
		gddcci.menu.Append(item)
		group = gr // go-gtk bug
	}
}



func (gddcci *GDDCci)onScroll(cbx *glib.CallbackContext) {
	arg := cbx.Args(0)
	event := *(**gdk3.EventScroll)(unsafe.Pointer(&arg))
	var stateUp uint   = 0x200000000
	var stateDown uint = 0x300000000

	if event.State & stateUp == stateUp && event.State & stateDown != stateDown {
		gddcci.level += 3
		if gddcci.level > 100 {
			gddcci.level = 100
		}
	}
	if event.State & stateDown == stateDown {
		gddcci.level -= 3
		if gddcci.level < 0 {
			gddcci.level = 0
		}
	}
	gddcci.updateSelectedMonitor()
}

func (gddcci *GDDCci)updateSelectedMonitor() {
	if gddcci.level != gddcci.lastLevel {
		if gddcci.selected == nil {
			panic(fmt.Errorf("monitor not selected"))
		}
		gddcci.selected.SetBrightness(gddcci.level)
		//label := fmt.Sprintf("Brightness: %v", gddcci.level)
		//
		//gddcci.icon.SetTitle(label)
		gddcci.lastLevel = gddcci.level
	}
}

func closeMenuItem() *gtk3.MenuItem {
	item := gtk3.NewMenuItemWithLabel("Exit!")
	item.Connect("activate", func() {
		gtk3.MainQuit()
	})
	return item
}

func createChoseMonitorMenuItem(info *goddcci.MonitorInfo, group *glib.SList, gddcci *GDDCci) (*gtk3.RadioMenuItem, *glib.SList) {
	label := fmt.Sprintf("Default: %v", info.PnPid)
	item := gtk3.NewRadioMenuItemWithLabel(group, label)
	item.Connect("activate", func() {
		gddcci.SelectMonitor(info.Id)
	})
	gr := item.GetGroup()
	return item, gr
}


func (gddcci *GDDCci) Main() {
	gtk3.Main()
}