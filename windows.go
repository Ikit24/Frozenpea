package main

import (
	"fmt"
	"errors"
	"strconv"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type SoundConfig struct {
	NotificationSound bool
	BreakReminder     bool
	BreakStartSound   bool
	BreakEndSound     bool
}

var soundConfig SoundConfig

func showBreakWindow(a fyne.App) fyne.Window {
	var w fyne.Window
	fyne.DoAndWait(func() {
		w = a.NewWindow("Break time!")

		img := canvas.NewImageFromFile("./assets/peakpx.jpg")
		img.FillMode = canvas.ImageFillStretch
		w.SetContent(img)
		w.SetFullScreen(true)
		w.SetCloseIntercept(func() {
		})
		w.Show()

		if nativeWin, ok := w.(driver.NativeWindow); ok {
			nativeWin.RunNative(func(ctx any) {
				if x11ctx, ok := ctx.(driver.X11WindowContext); ok {
					fmt.Println("Got X11 context, window:", x11ctx.WindowHandle)
					conn, err := xgb.NewConn()
					if err != nil {
						fmt.Println("Connection failed:", err)
						return
					}

					cookie := xproto.GrabKeyboard(conn, false, xproto.Window(x11ctx.WindowHandle), 0, xproto.GrabModeAsync, xproto.GrabModeAsync)
					reply, err := cookie.Reply()
					if err != nil {
						fmt.Println("Cannot grab keyboard:", err)
						return
					}
					if reply.Status != xproto.GrabStatusSuccess {
						fmt.Println("Grab failed, status:", reply.Status)
					}
				}
			})
		} else {
			fmt.Println("NOT a native window!")
		}
	})
	return w
}

func showNotification(a fyne.App) fyne.Window {
	var n fyne.Window

	fyne.DoAndWait(func() {
		n = a.NewWindow("Break in 1 minute!")
		img := canvas.NewImageFromFile("./assets/notify.jpg")
		img.FillMode = canvas.ImageFillStretch

		topText := canvas.NewText("Break starts in 1 minute! Save your work.", color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		topText.Alignment = fyne.TextAlignCenter
		topText.TextSize = 24
		topText.TextStyle = fyne.TextStyle{Bold: true}

		textContainer := container.NewVBox(
			topText,
			widget.NewLabel(""),
		)
		
		content := container.NewStack(img, container.NewCenter(textContainer))

		n.SetContent(content)
		n.Resize(fyne.NewSize(1265, 650))
		n.SetFixedSize(true)
		n.SetCloseIntercept(func() {
		})
		n.Show()
	})

	return n
}

func createSoundCheckBox(label string, configField *bool) *fyne.Container {
	check := widget.NewCheck("", func(checked bool) {
	*configField = checked
	})
	text := canvas.NewText(label, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	text.TextStyle = fyne.TextStyle{Bold: true}
	textBox := container.NewHBox(check, text)
	return textBox
}

func startupWindow(a fyne.App, setupDone chan bool) {
	var start fyne.Window
	
	fyne.DoAndWait(func () {
		start = a.NewWindow("Welcome to FrozenPea")
		workEntry := widget.NewEntry()
		workEntry.Resize(fyne.NewSize(150, 60))
		
		breakEntry := widget.NewEntry()
		breakEntry.Resize(fyne.NewSize(150, 60))
		workDur := widget.NewFormItem("Session duration (minutes):", workEntry)
		breakDur := widget.NewFormItem("Break duration (minutes):   ", breakEntry)

		form := widget.NewForm(workDur, breakDur)

		notifyBreakReminder := createSoundCheckBox("Break reminder  ", &soundConfig.BreakReminder)
		playBreakReminderSample := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {playSound("./assets/before_break.mp3")})
		beforeBreakReminder := container.NewHBox(notifyBreakReminder, playBreakReminderSample)

		notifyStartBreak := createSoundCheckBox("Break starting    ", &soundConfig.BreakStartSound)
		playStartBreakSample := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {playSound("./assets/break_start.mp3")})
		startBreakWithButton := container.NewHBox(notifyStartBreak, playStartBreakSample)

		notifyEndBreak := createSoundCheckBox("Break ending      ", &soundConfig.BreakEndSound)
		playEndBreakSample := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {playSound("./assets/break_end.mp3")})
		endBreakButton := container.NewHBox(notifyEndBreak, playEndBreakSample)

		confirmButton := widget.NewButton("Confirm changes", func() {
			_, err := strconv.Atoi(workEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("Please enter a valid number"), start)
				return
			}

			_, err = strconv.Atoi(breakEntry.Text)
			if err != nil {
				dialog.ShowError(errors.New("Please enter a valid number"), start)
				return
			}
			appConfig.WorkDuration = workEntry.Text
			appConfig.BreakDuration = breakEntry.Text
			setupDone <- true
			start.Close()
		})

		cancelSess := widget.NewButton("Cancel session", func() {
			a.Quit()
		})
		cancelSess.Resize(fyne.NewSize(150, 60))

			img := canvas.NewImageFromFile("./assets/fpea.jpg")
			img.FillMode = canvas.ImageFillStretch

			formContent := container.NewVBox(
			form,beforeBreakReminder,
			startBreakWithButton, endBreakButton,
			confirmButton, cancelSess,
			)
			content := container.NewStack(img, container.NewCenter(formContent))

			start.SetContent(content)
			start.Resize(fyne.NewSize(1000, 650))
			start.SetFixedSize(true)
			start.Show()
	})
}
