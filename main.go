package main

import (
	"bufio"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	userProjectConfPath := filepath.Join(cwd, "go-online.yml")
	userProjectConf := types.UserConf{}
	if _, err := os.Stat(userProjectConfPath); os.IsExist(err) {
		userProjectConfData, err := ioutil.ReadFile(userProjectConfPath)
		if err != nil {
			fmt.Println(err)
		}
		if err = yaml.Unmarshal(userProjectConfData, &userProjectConf); err != nil {
			fmt.Println(err)
		}
	}
	userProjectConf.SetDefault()
	// *************************************************

	// receive message from stdin
	inputChan := make(chan string, 10)
	go ReadMessage(inputChan)

	// compile and start gdb
	Compile(userProjectConf.ProjectName)
	gdb, err := gdb.New(func(notification map[string]interface{}) {
		fmt.Println(notification)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	gdb.Send("file-exec-and-symbols", "Debug/"+userProjectConf.ProjectName)

	// read stdin message and send to gdb
	for msg := range inputChan {
		if msg == "quit" {
			gdb.Exit()
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			os.Exit(0)
		}
		ret, err := gdb.Send(msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
			fmt.Fprintln(os.Stderr, err)
			close(input)
			return
		}
		text = text[:len(text)-1]
		input <- text
	}
}

// Compile read makefile and compile
func Compile(pn string) {
	// create Debug/temp folder if not exists
	err := os.MkdirAll("Debug/temp", os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// generate runable file
	cmd := exec.Command("make", "-f", "Makefile")
	cmdout, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go io.Copy(os.Stderr, cmdout)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
