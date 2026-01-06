package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func main () {
	a := app.New()
	fmt.Println("Session is running.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() { <-sigs; done <- true} ()

	isBreak := false
	go func () {
		time.Sleep(1 * time.Second)
		var w fyne.Window
		mainLoop: for {
			if isBreak {
				select {
				case <- time.After(10 * time.Minute):
					fmt.Println("Break is over, you can continue!")
					fyne.Do(func() {
					w.Close()
					})
				case <- done:
				break mainLoop
				}
				isBreak = false
			} else {
				var n fyne.Window
				select {
				case <- time.After(49 * time.Minute):
					fmt.Println("Mandatory break starts in 1 minute! Make sure you save your work.")
					n = showNotification(a)
				case <- done:
					break mainLoop
				}
				select {
				case <- time.After(10 * time.Minute):
					fmt.Println("Starting mandatory break.")
					fyne.Do(func() { n.Close() })
					w = showBreakWindow(a)
				case <- done:
					break mainLoop
				}
			}

			isBreak = true
		}
	} ()
	dummy := a.NewWindow("")
	dummy.Resize(fyne.NewSize(1, 1))
	dummy.Hide()
	a.Run()
}

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
