package gui

import (
	"github.com/zurek87/go-gtk3/gtk3"
	"fmt"
)

type guiNoMonitorError struct {
	msg string
}

func (err guiNoMonitorError) Error() string {
	return fmt.Sprintf("No monitor: %v", err.msg)
}

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

