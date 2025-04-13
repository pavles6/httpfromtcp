package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		contents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)

			if n > 0 {
				parts := strings.Split(string(b[:8]), "\n")

				contents += parts[0]
				if len(parts) > 1 { // we reached \n
					ch <- contents
					contents = parts[1]
				}
			}

			if err != nil {
				if errors.Is(io.EOF, err) {
					break
				}

				fmt.Println("Error reading file:", err)
				return
			}
		}
		close(ch)
	}()

	return ch
}

func main() {
	messages, err := os.Open("messages.txt")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer messages.Close()

	for line := range getLinesChannel(messages) {
		fmt.Printf("read: %s\n", line)
	}
}
