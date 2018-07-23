package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cyrus-and/gdb"
	"github.com/sysu-go-online/gdb_service/types"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	// *************read configure file*****************
	cwd, err := os.Getwd()
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}
	userProjectConfPath := filepath.Join(cwd, "go-online.yml")
	userProjectConf := types.UserConf{}
	if _, err := os.Stat(userProjectConfPath); os.IsExist(err) {
		userProjectConfData, err := ioutil.ReadFile(userProjectConfPath)
		if err != nil {
			PrintError(err)
		}
		if err = yaml.Unmarshal(userProjectConfData, &userProjectConf); err != nil {
			PrintError(err)
		}
	}
	userProjectConf.SetDefault()
	// *************************************************

	// receive message from stdin
	inputChan := make(chan string, 0)
	go ReadMessage(inputChan)

	// compile and start gdb
	Compile(userProjectConf.ProjectName)
	gdb, err := gdb.New(func(notification map[string]interface{}) {
		fmt.Println(notification)
	})
	// Read out put and send into stdout
	go func() {
		for {
			msg := make([]byte, 5)
			n, err := gdb.Read(msg)
			if err != nil {
				PrintError(err)
				os.Exit(1)
			}
			if n != 0 {
				fmt.Print(string(msg))
			}
		}
	}()
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}
	ret, err := gdb.CheckedSend("file-exec-and-symbols", "Debug/"+"main")
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}
	fmt.Println(ret)

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
			PrintError(err)
		}
		fmt.Println(ret)
	}
}

// ReadMessage read message from stdin
func ReadMessage(input chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			PrintError(err)
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
func Compile(pn string) {
	// create Debug/temp folder if not exists
	err := os.MkdirAll("Debug/temp", os.ModePerm)
	if err != nil {
		PrintError(err)
	}

	// generate runable file
	cmd := exec.Command("make", "-f", "Makefile")
	cmdout, err := cmd.StderrPipe()
	if err != nil {
		PrintError(err)
	}
	go io.Copy(os.Stderr, cmdout)
	err = cmd.Run()
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}
}

// PrintError print error as struct
func PrintError(err error) {
	type Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	retError := Error{"error", err.Error()}
	byteError, _ := json.Marshal(retError)
	fmt.Fprintln(os.Stderr, byteError)
}
