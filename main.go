package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type iPInfo struct {
	Data []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need a file name")
		return
	}
	inputFile, inputError := os.Open(os.Args[1]) //变量指向os.Open打开的文件时生成的文件句柄
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n")
		return
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	lineCounter := 0
	for {
		inputString, readerError := inputReader.ReadString('\n')
		//inputString, readerError := inputReader.ReadBytes('\n')
		if readerError == io.EOF {
			return
		}
		lineCounter++
		fmt.Printf("%d : %s\r\n", lineCounter, trans(strings.TrimSpace(inputString)))
	}
}

func trans(line string) string {
	reg, _ := regexp.Compile("\\d+\\.\\d+\\.\\d+\\.\\d+")
	result := reg.FindAll([]byte(line), -1)

	var iPDesc string
	for i := 0; i < len(result); i++ {
		iPDesc += getIPInfo(string(result[i])) + "\t"
	}
	return fmt.Sprintf("%s\t\t%s", line, iPDesc)
}

func getIPInfo(ip string) string {
	resp, err := http.Get("http://192.168.202.81:12101/find?ip=" + ip)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	var iPJson iPInfo
	err = json.Unmarshal(body, &iPJson)
	if err != nil {
		fmt.Println("Decode json fail")
	}

	if len(iPJson.Data) > 0 {
		return iPJson.Data[0]
	}

	return ""
}
