package main

import (
	"fmt"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/container"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func showBreakWindow(a fyne.App) fyne.Window {
	var w fyne.Window
	fyne.DoAndWait(func() {
		w = a.NewWindow("Break time!")

		img := canvas.NewImageFromFile("./peakpx.jpg")
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
		img := canvas.NewImageFromFile("./notify.jpg")
		img.FillMode = canvas.ImageFillStretch

		topText := canvas.NewText("Break starts in 1 minute! Save your work.", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		topText.Alignment = fyne.TextAlignCenter
		topText.TextSize = 24
		topText.TextStyle = fyne.TextStyle{Bold: true}

		textContainer := container.NewVBox(
			topText,
			widget.NewLabel(""),
		)
		
		content := container.NewStack(img, container.NewBorder(textContainer, nil, nil, nil))

		n.SetContent(content)
		n.Resize(fyne.NewSize(1265, 650))
		n.SetFixedSize(true)
		n.SetCloseIntercept(func() {
		})
		n.Show()
	})
	return n
}

func startupWindow(a fyne.App, setupDone chan bool) {
	var start fyne.Window
	
	fyne.DoAndWait(func () {
		start = a.NewWindow("Setup")
		workEntry := widget.NewEntry()
		breakEntry := widget.NewEntry()
		workDur := widget.NewFormItem("Session duration", workEntry)
		breakDur := widget.NewFormItem("Break duration", breakEntry)

		form := widget.NewForm(workDur, breakDur)
		
		confirmButton := widget.NewButton("Confirm changes", func() {
			appConfig.WorkDuration = workEntry.Text
			appConfig.BreakDuration = breakEntry.Text
			setupDone <- true
			start.Close()
		})

			content := container.NewVBox(form, confirmButton)
			start.SetContent(content)
			start.Show()
	})
}
