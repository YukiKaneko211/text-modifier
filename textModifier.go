package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	argsLength := len(args)

	if argsLength == 3 {
		// check files exist
		for i := 1; i < argsLength; i++ {
			_, err := os.Open(os.Args[i])
			if err != nil {
				fmt.Printf("no such file: %v\n", err.Error())
				os.Exit(1)
			}
		}

		// if files exist
		fmt.Println("Successfully started")

	} else { // not enough or too many arguments
		fmt.Println("Please specify files to read and write.")
	}
}
