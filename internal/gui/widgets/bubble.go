package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewMessageBubble(text string, isMe bool) fyne.CanvasObject {
	label := widget.NewLabel(text)
	// label.Wrapping = fyne.TextWrapWord

	content := container.NewStack(container.NewPadded(label))
	if isMe {
		return container.NewHBox(layout.NewSpacer(), content)
	}
	return container.NewHBox(content, layout.NewSpacer())
}
