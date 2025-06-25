package gui

import (
	"ChessCoordinatesApp/game"

	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PracticeGame struct {
	// UI Widgets
	CoordinateLabel *widget.Label
	ScoreLabel      *widget.Label
	FeedbackRect    *canvas.Rectangle

	// Game State
	CurrentSquare game.Square
	RightScore    int
	WrongScore    int

	// Communication & Control
	AnswerChannel chan int

	CancelFunc context.CancelFunc
	App        fyne.App
	Window     fyne.Window
}

func MainMenu(app fyne.App, window fyne.Window) fyne.CanvasObject {

	practiceModeButton := widget.NewButton("Practice Mode", func() {

		window.SetContent(freePracticeModeStartScreen(app, window))

	})

	blitzModeButton := widget.NewButton("Blitz Mode", func() {
		dialog.ShowInformation("Not Possible", "Still in development", window)
	})
	quitButton := widget.NewButton("Quit", func() {
		app.Quit()
	})

	menuContent := container.NewCenter(
		container.NewVBox(
			widget.NewLabel("Chess Personal Trainer"),
			practiceModeButton,
			blitzModeButton,
			quitButton,
		))

	return menuContent
}

func freePracticeModeStartScreen(app fyne.App, window fyne.Window) fyne.CanvasObject {

	startGameButton := widget.NewButton("Start Game", func() {

		window.SetContent(countdownMenu(app, window))
	})

	startGameMenu := container.NewCenter(
		container.NewVBox(
			startGameButton,
		))

	return startGameMenu
}

func countdownMenu(app fyne.App, window fyne.Window) fyne.CanvasObject {

	countdownText := canvas.NewText("", color.White)
	countdownText.Alignment = fyne.TextAlignCenter
	countdownText.TextStyle = fyne.TextStyle{Bold: true}
	countdownText.TextSize = 72

	content := container.NewCenter(countdownText)

	go func() {
		for i := 3; i > 0; i-- {

			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Content: fmt.Sprintf("%d", i),
			})
			fyne.DoAndWait(func() {
				fyne.CurrentApp().Driver().CanvasForObject(countdownText).Refresh(countdownText)
				countdownText.Text = fmt.Sprintf("%d", i)
				canvas.Refresh(countdownText)
				time.Sleep(1 * time.Second)
			})
		}

		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Content: "GO!",
		})

		fyne.DoAndWait(func() {
			countdownText.Text = "GO!"
			canvas.Refresh(countdownText)
			time.Sleep(500 * time.Millisecond)
		})

		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Content: "Starting game...",
		})
		fyne.DoAndWait(func() {
			window.SetContent(FreePracticeMode(app, window))
		})
	}()

	return content
}

func FreePracticeMode(app fyne.App, window fyne.Window) fyne.CanvasObject {

	gameSession := &PracticeGame{
		RightScore: 0,
		WrongScore: 0,
		App:        app,
		Window:     window,
	}

	gameSession.ScoreLabel = widget.NewLabel(fmt.Sprintf("Right: %d | Wrong: %d", gameSession.RightScore, gameSession.WrongScore))
	gameSession.CoordinateLabel = widget.NewLabelWithStyle("Loading...", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	//gameSession.CoordinateLabel.TextSize = 48

	gameSession.FeedbackRect = canvas.NewRectangle(color.Transparent)
	gameSession.FeedbackRect.SetMinSize(fyne.NewSize(200, 50))
	gameSession.AnswerChannel = make(chan int)

	whiteButton := widget.NewButton("White", func() {

		if !gameSession.CurrentSquare.IsBlack() {

			gameSession.AnswerChannel <- 1
		} else {

			gameSession.AnswerChannel <- 2
		}

	})

	blackButton := widget.NewButton("Black", func() {

		if gameSession.CurrentSquare.IsBlack() {

			gameSession.AnswerChannel <- 3
		} else {
			gameSession.AnswerChannel <- 4
		}

	})

	quitButton := widget.NewButton("Quit", func() {

		if gameSession.CancelFunc != nil {
			gameSession.CancelFunc()
		}
		gameSession.Window.SetContent(MainMenu(app, window))
	})

	gameMenu := container.NewVBox(

		gameSession.ScoreLabel,
		layout.NewSpacer(),
		gameSession.FeedbackRect,
		gameSession.CoordinateLabel,
		layout.NewSpacer(),
		container.NewGridWithColumns(2, whiteButton, blackButton),
		quitButton,
	)

	var ctx context.Context
	ctx, gameSession.CancelFunc = context.WithCancel(context.Background())
	go gameSession.runGameLoop(ctx)
	return gameMenu

}

func (g *PracticeGame) runGameLoop(ctx context.Context) {

	for {

		select {

		case <-ctx.Done():
			close(g.AnswerChannel)

		default:
			fyne.DoAndWait(func() {

				g.CurrentSquare = game.GenerateRandomSquare()
				g.CoordinateLabel.SetText(g.CurrentSquare.String())
				g.FeedbackRect.FillColor = color.Transparent
				g.CoordinateLabel.Refresh()
				canvas.Refresh(g.FeedbackRect)
			})
			select {

			case guess := <-g.AnswerChannel:
				fyne.DoAndWait(func() {
					if guess == 1 || guess == 3 {

						g.RightScore++
						g.FeedbackRect.FillColor = color.RGBA{0, 255, 0, 128}
					} else {

						g.WrongScore++
						g.FeedbackRect.FillColor = color.RGBA{255, 0, 0, 128}
					}

					g.ScoreLabel.SetText(fmt.Sprintf("Right: %d | Wrong: %d", g.RightScore, g.WrongScore))

					g.ScoreLabel.Refresh()
					canvas.Refresh(g.FeedbackRect)
				})

			case <-ctx.Done():
				return

			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
