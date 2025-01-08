package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type MessageFile struct {
	Sets []MessageSet
}

type MessageSet struct {
	SetName  string
	Messages []LoadingMessage
}

type LoadingMessage struct {
	Text       string
	MinSeconds int
	MaxSeconds int
}

const messageConfigFilePath = "Messages/messages.json"

func main() {
	var setFlag = flag.String("s", "default", "the set of messages to loop through")
	flag.Parse()

	file, err := os.ReadFile(messageConfigFilePath)
	if err != nil {
		log.Fatalf("error opening messages config: %v", err)
	}

	currentSet, err := parseConfig(file, *setFlag)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if currentSet == nil || len(currentSet.Messages) < 1 {
		log.Fatalf("no messages in current set %+v", currentSet)
	}

	reader := bufio.NewReader(os.Stdin)

	killChan := make(chan byte)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		currentSet.PrintMessages(killChan)
	}()

	// listen for an enter to end the program
	_, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error from reader: %v", err)
	}

	killChan <- '1'

	wg.Wait()
}

func (set *MessageSet) PrintMessages(killChan chan byte) {
	messageComplete := true
	var messageCompletionTime time.Time

	freshMessages := set.Messages
	var staleMessages []LoadingMessage

	clearScreen()

	for {
		if messageComplete {
			var nextMessage LoadingMessage

			if len(staleMessages) == 0 {
				clearScreen()
			}

			if len(freshMessages) > 1 {
				nextMessageIndex := rand.Intn(len(freshMessages))
				nextMessage = freshMessages[nextMessageIndex]

				freshMessages = append(freshMessages[:nextMessageIndex], freshMessages[nextMessageIndex+1:]...)
				staleMessages = append(staleMessages, nextMessage)
			} else {
				nextMessage = freshMessages[0]
				staleMessages = append(staleMessages, nextMessage)
				freshMessages = staleMessages
				staleMessages = []LoadingMessage{}
			}

			fmt.Print(nextMessage.Text)
			messageComplete = false

			messageDuration := rand.Intn(nextMessage.MaxSeconds) + nextMessage.MinSeconds

			messageCompletionTime = time.Now().Add(time.Duration(messageDuration) * time.Second)
			continue
		} else {
			fmt.Print(".")

			if time.Now().After(messageCompletionTime) {
				fmt.Print("\n")
				messageComplete = true
			}
		}

		randomSleep := rand.Intn(500) + 100
		time.Sleep(time.Duration(randomSleep) * time.Millisecond)

		select {
		case <-killChan:
			clearScreen()
			return
		default:
		}
	}
}

func parseConfig(configData []byte, messageSetName string) (*MessageSet, error) {

	var messageFile MessageFile
	err := json.Unmarshal(configData, &messageFile)
	if err != nil {
		return nil, fmt.Errorf("error decoding message config file: %v", err)
	}

	var currentSet MessageSet
	for _, set := range messageFile.Sets {
		if set.SetName == messageSetName {
			currentSet = set
		}
	}

	return &currentSet, nil
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
		break
	case "windows":
		cmd = exec.Command("cls")
	}
	if cmd == nil {
		log.Fatalf("unsupported platform: %v", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("error clearing screen: %v", err)
	}
}
