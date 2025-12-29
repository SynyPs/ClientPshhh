package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoginScreen(onLogin func(server string)) fyne.CanvasObject {
	serverEntry := widget.NewEntry()
	serverEntry.PlaceHolder = "URL (ws://)"
	serverEntry.Text = "ws://localhost:8443/ws"

	loginBtn := widget.NewButton("Connect", func() {
		if serverEntry.Text != "" {
			onLogin(serverEntry.Text)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Pshhh Login"),
		serverEntry,
		loginBtn,
	)

	return container.NewCenter(form)
}
