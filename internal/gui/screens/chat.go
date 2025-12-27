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
	"github.com/fhg/ClientPshhh/internal/service"
)

func ChatScreen(svc *service.Orchestrator, data binding.ExternalStringList) fyne.CanvasObject {
	contactsData := binding.BindStringList(&[]string{})
	contacts, _ := svc.GetContacts()
	for _, c := range contacts {
		contactsData.Append(c.Name)
	}

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
		val, _ := contactsData.GetValue(id)
		log.Println("Выбран чат :", val)
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
					} else {
						// Обновляем список (в идеале через binding, но пока можно так)
						log.Println("Контакт добавлен!")
						// contactsData.Append(nameEntry.Text)
					}
				}
			},
			// Тут нужно передать родительское окно.
			// Либо прокинь window в ChatScreen, либо используй глобальную переменную (плохо, но быстро)
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
	})

	input := widget.NewEntry()
	input.PlaceHolder = "Type send..."

	sendBtn := widget.NewButton("Send", func() {
		if input.Text == "" {
			return
		}

		go svc.SendMessage(input.Text)
		input.SetText("")
	})

	inputArea := container.NewBorder(nil, nil, nil, sendBtn, input)

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return container.NewMax(widget.NewLabel("template"))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			val, _ := di.(binding.String).Get()

			// binding.StringList хранит только текст.
			// В реальном проекте используйте binding.UntypedList с кастомными структурами.

			containerObj := co.(*fyne.Container)
			containerObj.Objects = nil

			isMe := false

			bubble := widgets.NewMessageBubble(val, isMe)
			containerObj.Add(bubble)
			containerObj.Refresh()
		},
	)

	leftPanel := container.NewBorder(addContactBtn, nil, nil, nil, contactList)
	rightPanel := container.NewBorder(nil, inputArea, nil, nil, list)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.3)

	return split
}
