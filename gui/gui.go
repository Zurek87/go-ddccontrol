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
}

func NewGui() GDDCci {

	ddcci, err := goddcci.InitDDCci()
	if err != nil {
		panic(err)
	}
	return GDDCci{ddcci: &ddcci, lastLevel: 101}
}

func (gddcci *GDDCci) Show() {
	if gddcci.visible {
		return
	}
	gddcci.initIcon()
}

func (gddcci *GDDCci) listMonitors() {
	list := gddcci.ddcci.MonitorList()
	for info := range list {
		fmt.Println(info)
	}
}

func (gddcci *GDDCci) initIcon() {
	gtk3.Init(nil)
	appName := fmt.Sprintf("Brightness on %v", gddcci.ddcci.MonitorName())
	glib.SetApplicationName(appName)

	menu := gtk3.NewMenu()
	menu.Append(closeMenuItem())

	menu.ShowAll()
	gddcci.icon = gtk3.NewStatusIconFromIconName("info")
	gddcci.icon.Connect("popup-menu", func(cbx *glib.CallbackContext) {
		menu.Popup(nil, nil, gtk3.StatusIconPositionMenu, gddcci.icon, uint(cbx.Args(0)), uint32(cbx.Args(1)))
	})
	gddcci.icon.Connect("scroll-event", func(cbx *glib.CallbackContext) {
		gddcci.onScroll(cbx)
	})

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
	gddcci.updateDefaultMonitor()
}

func (gddcci *GDDCci)updateDefaultMonitor() {
	if gddcci.level != gddcci.lastLevel {
		gddcci.ddcci.SetBrightness(gddcci.level)
		label := fmt.Sprintf("Brightness: %v", gddcci.level)

		gddcci.icon.SetTitle(label)
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



func (gddcci *GDDCci) Main() {
	gtk3.Main()
}