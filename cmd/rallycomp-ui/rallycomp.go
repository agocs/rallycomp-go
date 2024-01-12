package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var DISPLAY_WIDTH float32 = 800
var DISPLAY_HEIGHT float32 = 480

func main() {
	a := app.New()
	w := a.NewWindow("RallyComp")
	w.Resize(fyne.NewSize(DISPLAY_WIDTH, DISPLAY_HEIGHT))
	w.SetFixedSize(true)
	w.SetContent(widget.NewLabel("Hello World"))
	w.ShowAndRun()
}
