package autorespond

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"fmt"
)

var (
	Alive bool
	Input chan string
	Output chan string
	Sigs chan os.Signal
)

type AIMessage struct {
	Type string `json:"type"`
	Contents string `json:"contents,omitempty"`
}

type AIReply struct {
	StatusMessage string `json:"statusMessage"`
	PrimaryMood float32 `json:"primaryMood"`
	MoodStability float32 `json:"moodStability"`
	ExposedPositivity float32 `json:"exposedPositivity"`
	PositivityOverload bool `json:"positivityOverload"`
	Response string `json:"response"`
}

type AIError struct {
	ErrorMessage string `json:"error"`
	Response *string `json:"response"` // can be null
}

func parseReply(bytes []byte) (AIReply, error) {
	var reply AIReply
	if err := json.Unmarshal(bytes, &reply); err != nil {
		log.Println("Autoresponder reply is invalid:", err)
		// might be an error instead
		var replyErr AIError
		if err := json.Unmarshal(bytes, &replyErr); err != nil {
			log.Println("Autoresponder reply is invalid:", err)
			return reply, errors.New("Autoresponder returned an unknown response")
		} else {
			response := replyErr.ErrorMessage
			if replyErr.Response != nil {
				response += " | " + *replyErr.Response
			}
			return reply, errors.New(response)
		}
	} else {
		return reply, nil
	}
}

func init() {
	Alive = true

	Input = make(chan string)
	Output = make(chan string)

	Sigs = make(chan os.Signal)
	signal.Notify(Sigs, syscall.SIGINT, syscall.SIGTERM)

	//cmd := exec.Command("python", "./answer/autorespond/tofu-ai/main.py", "chat")
	cmd := exec.Command("python", "./answer/autorespond/tofu-ai/main.py", "jsonio")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)

	if err := cmd.Start(); err != nil { // TODO ignored because it is handled below, improve this
		//log.Fatal("Buffer Error:", err)
	}

	var mutex sync.Mutex

	go func() {
		for {
			select {
			case msg := <-Input:
				mutex.Lock()

				msg = strings.ReplaceAll(msg, "\n", ".")
				log.Println("Autoresponder received:", msg)

				// build the json message
				bytes, err := json.Marshal(AIMessage{
					Type: "group message",
					Contents: msg,
				})
				log.Println(string(bytes[:]))
				if err != nil {
					log.Println("Failed to marshal message", err)
					break
				}

				stdin.Write(append(bytes, byte('\n')))

				// read the json reply 
				bytes, err = rd.ReadBytes(byte('\n'))
				if err != nil {
					Output <- "Autoresponder is dead, please restart tofu"
					log.Println("Autoresponder error:", err)
					return
				}

				// parse the json reply
				if reply, err := parseReply(bytes); err == nil {
					log.Println("Autoresponder reply:", reply.Response)
					Output <- reply.Response
				} else {
					Output <- err.Error()
				}

				mutex.Unlock()
				break
			case <-Sigs:
				stdin.Close()
				stdout.Close()
				return
			}
		}
	}()

	//ticker := time.NewTicker(time.Second)
	//go func() {
	//	defer ticker.Stop()
	//	for {
	//		select {
	//		case t := <-ticker.C:
	//			mutex.Lock()

	//			log.Println("Current time: ", t)

	//			mutex.Unlock()
	//		case <-Sigs:
	//			return
	//		}
	//	}
	//}()

	log.Println("Autoresponder running")
	fmt.Println("Autoresponder running")
}
