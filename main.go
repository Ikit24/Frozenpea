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

func initAudio() {
	f, _ := os.Open("./assets/before_break.mp3")
	defer f.Close()

	streamer, format, _ := mp3.Decode(f)
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func playSound(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	streamer, _, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func main () {
	initAudio()
	a := app.New()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() { <-sigs; done <- true} ()

	go func () {
		time.Sleep(100 * time.Millisecond)
		setupDone := make(chan bool)
		startupWindow(a, setupDone)
		playSound("./assets/intro.mp3")
		<-setupDone

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
					//Break end
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
					//Break in 1 min
					n = showNotification(a)
					playSound("./assets/before_break.mp3")
				case <- done:
					break mainLoop
				}
				select {
				case <- time.After(1 * time.Minute):
					//Break start
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
