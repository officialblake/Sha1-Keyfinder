package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var wordPath string
	var tag string
	var messagePath string
	fmt.Print("Path to word bank: ")
	fmt.Scan(&wordPath)

	fmt.Print("Path to message: ")
	fmt.Scan(&messagePath)

	fmt.Print("Enter desried SHA1 tag: ")
	fmt.Scan(&tag)

	file, err := os.Open(wordPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = os.Open(messagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		hexString := hex.EncodeToString([]byte(line))
		cmd := exec.Command("openssl", "dgst", "-sha1", "-mac", "HMAC", "-macopt", "hexkey:"+hexString, messagePath)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error running command:", err)
			return
		}

		outputString := strings.TrimSpace(string(output))
		parts := strings.Split(outputString, "= ")
		if len(parts) == 2 {
			hashValue := strings.TrimSpace(parts[1])

			if hashValue == tag {
				fmt.Println("Key found:", line)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
