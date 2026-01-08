package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
