package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type MessageFile struct {
	Sets []MessageSet
}

type MessageSet struct {
	Flag     string
	Messages []LoadingMessage
}

type LoadingMessage struct {
	Text       string
	MinSeconds int
	MaxSeconds int
}

const messageConfigFilePath = "Messages/messages.json"

func main() {
	file, err := os.ReadFile(messageConfigFilePath)
	if err != nil {
		log.Fatalf("error opening messages config: %v", err)
	}

	var messageFile MessageFile
	err = json.Unmarshal(file, &messageFile)
	if err != nil {
		log.Fatalf("error decoding message config file: %v", err)
	}

	var setFlag = flag.String("s", "default", "the set of messages to loop through")
	flag.Parse()

	var currentSet MessageSet
	for _, set := range messageFile.Sets {
		if set.Flag == *setFlag {
			currentSet = set
		}
	}

	if len(currentSet.Messages) < 1 {
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

	for {
		if messageComplete {
			nextMessage := set.Messages[rand.Intn(len(set.Messages))]
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

		time.Sleep(500 * time.Millisecond)

		select {
		case <-killChan:
			return
		default:
		}
	}
}
