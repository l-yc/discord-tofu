package brain

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
)

var (
	Input chan BrainInput
	Output chan BrainOutput
	Sigs chan os.Signal
	State chan BrainState
)

const (
	IDLE_DURATION = 15 * time.Second
	POLL_INTERVAL = 15 * time.Second
)

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
	Input = make(chan BrainInput)
	Output = make(chan BrainOutput)

	Sigs = make(chan os.Signal)
	signal.Notify(Sigs, syscall.SIGINT, syscall.SIGTERM)

	cmd := exec.Command("python", "./brain/tofu-ai/main.py", "jsonio")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)

	if err := cmd.Start(); err != nil { // TODO ignored because it is handled below, improve this
		//log.Fatal("Buffer Error:", err)
	}

	State = make(chan BrainState, 1) // must be buffered, because receiver might not be ready
	mood := "Just woke up!"
	status := BrainStateStatusOnline

	State <- BrainState{ Mood: mood, Status: status }
	idleTimer := time.AfterFunc(IDLE_DURATION, func() {
		log.Println("Setting idle status")
		status = BrainStateStatusIdle
		State <- BrainState{ Mood: mood, Status: status }
	})

	var mutex sync.Mutex
	go func() {
		for {
			select {
			case input := <-Input:
				mutex.Lock()

				msg := strings.ReplaceAll(input.Content, "\n", ".")
				log.Println("Autoresponder received:", msg)

				// build the json message
				bytes, err := json.Marshal(AIMessage{
					Type: []string{
							AIMessageTypeGroup,
							AIMessageTypeDirect,
							AIMessageTypeStatus,
						}[input.Type],
					Contents: input.Content,
				})
				log.Println("Send to autoresponder:", string(bytes[:]))
				if err != nil {
					log.Println("Failed to marshal message", err)
					mutex.Unlock()
					break
				}

				stdin.Write(append(bytes, byte('\n')))

				// read the json reply 
				bytes, err = rd.ReadBytes(byte('\n'))
				if err != nil {
					Output <- BrainOutput{ Error: err, Content: "Autoresponder is dead, please restart tofu" }
					log.Println("Autoresponder error:", err)
					mutex.Unlock()
					return
				}

				// parse the json reply
				if reply, err := parseReply(bytes); err == nil {
					if input.Type == BrainInputTypeStatus {
						log.Println("Autoresponder state:", reply.StatusMessage)
						mood = reply.StatusMessage
						State <- BrainState{ Mood: mood, Status: status }
					} else {
						log.Println("Autoresponder reply:", reply.Response)
						Output <- BrainOutput{ Error: nil, Content: reply.Response }
						if reply.Response != "" {
							log.Println("Resetting idle timer")
							State <- BrainState{ Mood: mood, Status: BrainStateStatusOnline }
							idleTimer.Stop()
							idleTimer.Reset(IDLE_DURATION)
						}
					}
				} else {
					Output <- BrainOutput{ Error: err, Content: err.Error() }
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

	ticker := time.NewTicker(POLL_INTERVAL)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				Input <- BrainInput{
					Type: BrainInputTypeStatus,
				}
				break
			case <-Sigs:
				return
			}
		}
	}()

	log.Println("Autoresponder running")
}
