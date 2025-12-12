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
		mainLoop: for {
			if isBreak {
				select {
				case <- time.After(10 * time.Second):
					fmt.Println("Break is over, you can continue!")
				case <- done:
				break mainLoop
				}
				isBreak = false
			} else {
				select {
				case <- time.After(20 * time.Second):
					fmt.Println("Starting mandatory break.")
					fmt.Println("About to show window")
					showBreakWindow(a)
					fmt.Println("Window shown")
				case <- done:
				break mainLoop
				}
			isBreak = true
			}
		}
	} ()
	a.Run()
}

func showBreakWindow(a fyne.App) {
	fyne.DoAndWait(func() {
		w := a.NewWindow("Break time!")

		w.SetContent(widget.NewLabel("Take a break!"))
		w.SetFullScreen(true)
		w.Show()

		time.AfterFunc(10 * time.Second, func() {
			fyne.Do(func() {
			w.Close()
			})
		})
	})
}
