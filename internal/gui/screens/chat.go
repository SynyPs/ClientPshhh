package screens

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fhg/ClientPshhh/internal/gui/widgets"
	"github.com/fhg/ClientPshhh/internal/models"
	"github.com/fhg/ClientPshhh/internal/service"
)

func ChatScreen(svc *service.Orchestrator, data binding.ExternalUntypedList) fyne.CanvasObject {
	var currentChatKey string = ""

	contactsData := binding.BindStringList(&[]string{})
	contacts, _ := svc.GetContacts()
	for _, c := range contacts {
		contactsData.Append(c.Name)
	}

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return container.NewMax(widget.NewLabel("template"))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			val, _ := di.(binding.Untyped).Get()

			msg := val.(models.UIMessage)
			// binding.StringList хранит только текст.
			// В реальном проекте используйте binding.UntypedList с кастомными структурами.

			containerObj := co.(*fyne.Container)
			containerObj.Objects = nil

			bubble := widgets.NewMessageBubble(msg.Text, msg.IsMe)
			containerObj.Add(bubble)
			containerObj.Refresh()
		},
	)

	contactList := widget.NewListWithData(
		contactsData,
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("user name"))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			val, _ := di.(binding.String).Get()
			label := co.(*fyne.Container).Objects[1].(*widget.Label)
			label.SetText(val)
		})

	contactList.OnSelected = func(id widget.ListItemID) {
		log.Println("Выбран чат :", contacts[id].PublicKey)

		if id < 0 || id >= len(contacts) {
			return
		}
		selectedContact := contacts[id]
		currentChatKey = selectedContact.PublicKey

		data.Set([]interface{}{})

		history, err := svc.LoadChatHistory(currentChatKey)
		if err != nil {
			log.Println("Ошибка загрузки истории:", err)
			return
		}
		for _, msg := range history {
			data.Append(msg)
		}

		list.ScrollToBottom()
	}

	addContactBtn := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		// Создаем попап диалог
		keyEntry := widget.NewEntry()
		keyEntry.PlaceHolder = "Вставь Public Key друга..."

		nameEntry := widget.NewEntry()
		nameEntry.PlaceHolder = "Имя (например Bob)"

		dialog.ShowForm("Новый контакт", "Сохранить", "Отмена",
			[]*widget.FormItem{
				widget.NewFormItem("Key", keyEntry),
				widget.NewFormItem("Name", nameEntry),
			},
			func(confirm bool) {
				if confirm && keyEntry.Text != "" {
					// Вызываем сервис для сохранения в БД
					err := svc.AddContact(keyEntry.Text, nameEntry.Text)
					if err != nil {
						log.Println("Ошибка добавления:", err)
						return
					}
					contactsData.Append(nameEntry.Text)

					newContact := models.Contact{
						PublicKey: keyEntry.Text,
						Name:      nameEntry.Text,
					}
					contacts = append(contacts, newContact)
				}
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
	})

	input := widget.NewEntry()
	input.PlaceHolder = "Type send..."

	sendBtn := widget.NewButton("Send", func() {
		if input.Text == "" {
			return
		}

		text := input.Text
		go svc.SendMessage(input.Text, currentChatKey)

		data.Append(models.UIMessage{
			Text: text,
			IsMe: true,
		})
		input.SetText("")
	})

	inputArea := container.NewBorder(nil, nil, nil, sendBtn, input)

	leftPanel := container.NewBorder(addContactBtn, nil, nil, nil, contactList)
	rightPanel := container.NewBorder(nil, inputArea, nil, nil, list)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.3)

	return split
}
