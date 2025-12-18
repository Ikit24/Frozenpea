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
		fmt.Println("Window shwon, attempting native access...")
		if nativeWin, ok := w.(driver.NativeWindow); ok {
			fmt.Println("Got natinve window interface")
			nativeWin.RunNative(func(ctx any) {
				fmt.Println("Inside RunNAtive callback")
				if x11ctx, ok := ctx.(*driver.X11WindowContext); ok {
					fmt.Println("Got X11 context, window:", x11ctx.WindowHandle)
					conn, err := xgb.NewConn()
					if err != nil {
						fmt.Println("Connection failed:", err)
						return
					}
					fmt.Println("Created X11 connection")

					cookie := xproto.GrabKeyboard(conn, false, xproto.Window(x11ctx.WindowHandle), 0, xproto.GrabModeAsync, xproto.GrabModeAsync)
					reply, err := cookie.Reply()
					if err != nil {
						fmt.Println("Cannot grab keyboard:", err)
						return
					}
					fmt.Println("Grab status:", reply.Status)
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
