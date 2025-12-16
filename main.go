package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/driver"

	//"github.com/BurntSushi/xgb"
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
		var w fyne.Window
		mainLoop: for {
			if isBreak {
				select {
				case <- time.After(10 * time.Second):
					fmt.Println("Break is over, you can continue!")
					fyne.Do(func() {
					w.Close()
					})
				case <- done:
				break mainLoop
				}
				isBreak = false
			} else {
				select {
				case <- time.After(20 * time.Second):
					fmt.Println("Starting mandatory break.")
					w = showBreakWindow(a)
				case <- done:
				break mainLoop
				}
			isBreak = true
			}
		}
	} ()
	a.Run()
}

func showBreakWindow(a fyne.App) fyne.Window {
	var w fyne.Window
	fyne.DoAndWait(func() {
		w = a.NewWindow("Break time!")
		w.SetContent(widget.NewLabel("Take a break!"))
		w.SetFullScreen(true)
		w.SetCloseIntercept(func() {

		})
		w.Show()
		if nativeWin, ok := w.(driver.NativeWindow); ok {
			nativeWin.RunNative(func(ctx any) {
				if x11, ok := ctx.(*driver.X11WindowContext); ok {
					fmt.Printf("Window ID: %d\n", x11.WindowHandle)
				}
			})
		}
	})
	return w
}
