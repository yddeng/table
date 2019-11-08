package main

import (
	"fmt"
	"github.com/yddeng/table"
	"os"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Printf("usage config\n")
		return
	}

	table.Start(os.Args[1])

	stop := make(chan struct{}, 1)
	select {
	case <-stop:
	}
}
