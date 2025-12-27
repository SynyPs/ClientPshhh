package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoginScreen(onLogin func(userName, server string)) fyne.CanvasObject {
	userNameEntry := widget.NewEntry()
	userNameEntry.PlaceHolder = "UserName"

	serverEntry := widget.NewEntry()
	serverEntry.PlaceHolder = "URL (ws://)"
	serverEntry.Text = "ws://localhost:8443/ws"

	loginBtn := widget.NewButton("Connect", func() {
		if userNameEntry.Text != "" {
			onLogin(userNameEntry.Text, serverEntry.Text)
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Pshhh Login"),
		userNameEntry,
		serverEntry,
		loginBtn,
	)

	return container.NewCenter(form)
}
