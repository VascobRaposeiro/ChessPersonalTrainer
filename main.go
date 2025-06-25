package main

import (
	"ChessCoordinatesApp/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Chess Personal Trainer")

	window.SetContent(gui.MainMenu(app, window))

	window.Resize(fyne.NewSize(400, 300))

	window.ShowAndRun()
}
