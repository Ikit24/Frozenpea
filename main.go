package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func main () {
	fmt.Println("Session is running.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() { <-sigs; done <- true} ()

	isBreak := false
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
			case <- time.After(50 * time.Second):
				fmt.Println("Starting mandatory break.")
			case <- done:
			break mainLoop
			}
		isBreak = true
		}
	}
}
