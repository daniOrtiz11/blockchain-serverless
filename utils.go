package main

import (
	"fmt"
	"os"
)

func bytestofile(b []byte) {
	f, err := os.Create(localfile)
	check(err)
	defer f.Close()
	n2, err := f.Write(b)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
