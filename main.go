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

	mainLoop: for {
		select {
		case <- time.After(10 * time.Second):
			fmt.Println("You are still on break")
		case <- done:
		break mainLoop
		}
	}
	fmt.Println()
	fmt.Println("Break is over, you can continue! Happy coding!")
}
