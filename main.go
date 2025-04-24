package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Config struct {
	UnrealPath    string
	UnrealVersion string
	ProjectName   string
	ProjectPath   string
}

func main() {
	// Read config.toml for required info about project
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatal("Error reading config.toml file")
	}

	// Prepare batch file
	path := config.UnrealPath + "UE_" + config.UnrealVersion + "/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.exe"
	path, err = exec.LookPath(path)
	if err != nil {
		log.Fatal("Error finding unreal path")
	}

	// Arguments to pass to unreal build tool
	var args [7]string
	args[0] = "-mode=GenerateClangDatabase"
	args[1] = `-project=` + config.ProjectPath + config.ProjectName + `.uproject`
	args[2] = "-game"
	args[3] = "-engine"
	args[4] = config.ProjectName + "Editor"
	args[5] = "Development"
	args[6] = "Win64"
	bat := exec.Command(path, args[0], args[1], args[2], args[3], args[4], args[5], args[6])

	// Prepare Error/Output Messages
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	bat.Stderr = &stderr
	bat.Stdout = &stdout

	// Run batch file
	err = bat.Run()
	if err != nil {
		fmt.Println("Error running unreal build tool")
		fmt.Println(stdout.String())
		log.Fatal(fmt.Sprint(err) + ": " + fmt.Sprint(stderr.String()))
	}

	// Open compile_commands.json to copy
	inputFile, err := os.Open(config.UnrealPath + "UE_" + config.UnrealVersion + "/compile_commands.json")
	if err != nil {
		log.Fatal("Error opening old compile commands: " + fmt.Sprint(err))
	}

	// Open new compile_commands.json
	outputFile, err := os.Create(config.ProjectPath + "compile_commands.json")
	if err != nil {
		log.Fatal("Error opening new compile commands: " + fmt.Sprint(err))
	}
	defer outputFile.Close()

	// Copy contents over
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		log.Fatal("Error copying compile commands: " + fmt.Sprint(err))
	}
	inputFile.Close()

	// Delete old compile_commands.json
	err = os.Remove(config.UnrealPath + "UE_" + config.UnrealVersion + "/compile_commands.json")
	if err != nil {
		log.Fatal("Error deleting old compile commands " + fmt.Sprint(err))
	}
}
