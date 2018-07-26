package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/cyrus-and/gdb"
	"github.com/sysu-go-online/gdb_service/types"
)

func main() {
	// receive message from stdin
	inputChan := make(chan string, 0)
	go ReadMessage(inputChan)

	// compile and start gdb
	Compile()
	gdb, err := gdb.New(func(notification map[string]interface{}) {
		PrintMessage(1, notification)
	})
	if err != nil {
		PrintMessage(3, map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	// Read user process output and send into stdout
	go func() {
		for {
			// use multiplicy of 3 as it use 3 bytes to store chinese character
			msg := make([]byte, 90)
			_, err := gdb.Read(msg)
			if err != nil {
				// PrintMessage(3, map[string]interface{}{"error": err.Error()})
				os.Exit(1)
			}
			msg = bytes.Trim(msg, "\x00")
			if len(msg) != 0 {
				PrintMessage(2, map[string]interface{}{"msg": string(msg)})
			}
		}
	}()
	ret, err := gdb.CheckedSend("file-exec-and-symbols", "Debug/"+"main")
	if err != nil {
		PrintMessage(3, map[string]interface{}{"error": err.Error(), "command": "file-exec-and-symbols"})
		os.Exit(1)
	}
	ret["command"] = "file-exec-and-symbols"
	PrintMessage(1, ret)

	// read stdin message and send to gdb
	for msg := range inputChan {
		if msg == "quit" {
			gdb.Exit()
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			os.Exit(0)
		}
		ret, err := gdb.CheckedSend(msg)
		if err != nil {
			PrintMessage(3, map[string]interface{}{"error": err.Error(), "command": msg})
			continue
		}
		ret["command"] = msg
		PrintMessage(1, ret)
	}
}

// ReadMessage read message from stdin
func ReadMessage(input chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			// PrintMessage(3, map[string]interface{}{"error": err.Error()})
			close(input)
			return
		}
		text = text[:len(text)-1]
		if len(text) >= 1 {
			input <- text
		}
	}
}

// Compile read makefile and compile
func Compile() {
	// create Debug/temp folder if not exists
	err := os.MkdirAll("Debug/temp", os.ModePerm)
	if err != nil {
		PrintMessage(3, map[string]interface{}{"error": err.Error()})
	}

	// generate runable file
	cmd := exec.Command("make", "-f", "Makefile")
	cmdout, err := cmd.StderrPipe()
	if err != nil {
		PrintMessage(3, map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	go io.Copy(os.Stderr, cmdout)
	err = cmd.Run()
	if err != nil {
		PrintMessage(3, map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}

// PrintMessage print gdb,error and user process output data
func PrintMessage(msgType int, msg map[string]interface{}) {
	Type := "gdb"
	if msgType == 2 {
		Type = "output"
	} else if msgType == 3 {
		Type = "error"
	}
	retMsg := types.ResponseData{Type, msg}
	byteMsg, _ := json.Marshal(retMsg)
	fmt.Println(string(byteMsg))
}
