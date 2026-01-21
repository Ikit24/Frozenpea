package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Config struct {
	WorkDuration  string
	BreakDuration string
}

var appConfig Config

speaker.Init(sampleRate, sampleRate.N(time.Second/10))

func playSound(filename string) {
	f, _ := os.Open(filename)
	defer f.Close()

	streamer, format, _ := mp3.Decode(f)
	defer streamer.Close()

	speaker.Play(streamer)
}

func main () {
	a := app.New()

	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() { <-sigs; done <- true} ()

	go func () {
		time.Sleep(100 * time.Millisecond)
		setupDone := make(chan bool)
		startupWindow(a, setupDone)
		<-setupDone

		fmt.Println("Session is running.")
		workMins, _ := strconv.Atoi(appConfig.WorkDuration)
		workDur := time.Duration(workMins) * time.Minute
		breakMins, _ := strconv.Atoi(appConfig.BreakDuration)
		breakDur := time.Duration(breakMins) * time.Minute

		isBreak := false
		var w fyne.Window
		mainLoop: for {
			if isBreak {
				select {
				case <- time.After(breakDur):
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
				case <- time.After(workDur - 1 * time.Minute):
					playSound(before_break)
					fmt.Println("Mandatory break starts in 1 minute! Make sure you save your work.")
					n = showNotification(a)
				case <- done:
					break mainLoop
				}
				select {
				case <- time.After(1 * time.Minute):
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
