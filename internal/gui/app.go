package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/fhg/ClientPshhh/internal/gui/screens"
	"github.com/fhg/ClientPshhh/internal/models"
	"github.com/fhg/ClientPshhh/internal/service"
)

type MainWindow struct {
	app     fyne.App
	window  fyne.Window
	service *service.Orchestrator

	// message list
	messages binding.ExternalUntypedList
}

func NewMainWindow(app fyne.App, svc *service.Orchestrator) *MainWindow {
	w := app.NewWindow("Pshhh")
	w.Resize(fyne.NewSize(400, 600))

	var datalist []interface{}
	data := binding.BindUntypedList(&datalist)

	appUI := &MainWindow{
		app:      app,
		window:   w,
		service:  svc,
		messages: data,
	}

	appUI.showLogin()

	return appUI
}

func (mv *MainWindow) ShowAndRun() {
	mv.window.ShowAndRun()
}

func (mv *MainWindow) showLogin() {
	content := screens.LoginScreen(func(server string) {
		err := mv.service.Connect(server)
		if err != nil {
			return
		}
		mv.showChat()
	})
	mv.window.SetContent(content)
}

func (mv *MainWindow) showChat() {
	mv.service.SetOnMessageReceived(func(msg string) {
		log.Printf("GUI: Получено сообщение для отображения: %s", msg)
		mv.messages.Append(models.UIMessage{
			Text: msg,
			IsMe: false,
		})

		mv.messages.Reload()
	})
	content := screens.ChatScreen(mv.service, mv.messages)
	mv.window.SetContent(content)
}
