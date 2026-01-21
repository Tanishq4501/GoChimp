package main

import (
	"sync"
)

type Recipient struct {
	Name  string
	Email string
}

func main() {
	recipientChannnel := make(chan Recipient)

	go func() {
		loadRecipient("./emails.csv", recipientChannnel)
	}()

	var wg sync.WaitGroup

	workerCount := 5

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go emailWorker(i, recipientChannnel, &wg)
	}

	wg.Wait()

}
