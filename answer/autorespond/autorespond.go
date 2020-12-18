package autorespond

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"fmt"
)

var (
	Alive bool
	Input chan string
	Output chan string
	Sigs chan os.Signal
)

func init() {
	Alive = true

	Input = make(chan string)
	Output = make(chan string)

	Sigs = make(chan os.Signal)
	signal.Notify(Sigs, syscall.SIGINT, syscall.SIGTERM)

	cmd := exec.Command("python", "./answer/autorespond/message-autoresponder/main.py", "chat")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)

	fmt.Println("AAAAAAAAAAAAAAAAA")
	if err := cmd.Start(); err != nil { // TODO ignored because it is handled below, improve this
		//log.Fatal("Buffer Error:", err)
	}

	go func() {
		for {
			select {
			case msg := <-Input:
				msg = strings.ReplaceAll(msg, "\n", ".")
				log.Println("Autoresponder received:", msg)
				stdin.Write([]byte(msg+"\n"))
				reply, err := rd.ReadString('\n')
				if err != nil {
					Output <- "Autoresponder is dead, please restart tofu"
					//log.Fatal("Read Error:", err)
					return
				}
				log.Println("Autoresponder reply:", reply)
				Output <- reply
				break
			case <-Sigs:
				stdin.Close()
				stdout.Close()
				return
			}
		}
	}()

	log.Println("Autoresponder running")
}
