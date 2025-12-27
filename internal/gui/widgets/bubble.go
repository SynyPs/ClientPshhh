package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewMessageBubble(text string, isMe bool) fyne.CanvasObject {
	label := widget.NewLabel(text)
	// label.Wrapping = fyne.TextWrapWord

	var bgCol color.Color
	if isMe {
		bgCol = theme.PrimaryColor()
	} else {
		bgCol = theme.ButtonColor()
	}

	bg := canvas.NewRectangle(bgCol)
	bg.CornerRadius = 10

	content := container.NewStack(bg, container.NewPadded(label))
	if isMe {
		return container.NewHBox(layout.NewSpacer(), content)
	}
	return container.NewHBox(content, layout.NewSpacer())
}
