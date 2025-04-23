package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

	// Read Paths.toml for required info about project
	file, err := os.Open("./paths.toml")
	if err != nil {
		log.Fatal("Error reading paths.toml file")
	}
	defer file.Close()

	data := make([]byte, 1000)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes: %q\n", count, data[:count])

	// Manipulate the data
	seperator := []byte("\n")
	lines := bytes.SplitAfter(data, seperator)

	unrealEnginePath := string(lines[1])
	unrealEngineVersion := string(lines[2])

	unrealEnginePath = strings.Split(unrealEnginePath, "'")[1]
	unrealEngineVersion = strings.Split(unrealEngineVersion, "'")[1]

	fmt.Println(unrealEnginePath + unrealEngineVersion + "/")
}
