package main

import (
	"log"
	"time"
)

func main() {
	commandch := make(chan command, 0) // unbuffered
	go commandRunner(commandch)

	// we can only run one command at a time, but don't want to block on running
	// it. we need to error if we try and run a command twice.
	for i := 0; i < 10; i++ {
		time.Sleep(100*time.Millisecond)
		command := command{t: b}
		select {
		case commandch <- command:
		default:
			log.Println("could not start command")
		}
	}

	log.Println("waiting for commands to finish")
	time.Sleep(10 * time.Second) // in reality we run infinitely
}

type commandType int

const (
	a commandType = iota
	b
	c
)

type command struct{
	t commandType
}

func commandRunner(commandch <-chan command) {
	for {
		cmd := <-commandch // pull command off the chan

		// run the command; commandch will be blocked until we loop around again
		log.Printf("running command %v", cmd)
		switch cmd.t {
		case a:
			runCommandA()
		case b:
			runCommandB()
		case c:
			runCommandC()
		}
		log.Printf("finished command %v", cmd)
	}
}

func runCommandA() {
	time.Sleep(2 * time.Second)
}

func runCommandB() {
	time.Sleep(4 * time.Second)
}

func runCommandC() {
	time.Sleep(6 * time.Second)
}
