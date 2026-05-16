package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "remote", "show", "origin")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var branch string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "HEAD branch: ") {
			tmp := strings.TrimSpace(line)
			branch = strings.TrimPrefix(tmp, "HEAD branch: ")
			break
		}
	}

	fmt.Printf("Default branch: %s\n", branch)
	fmt.Println("Checkout to default branch")

	cmd = exec.Command("git", "checkout", branch)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("git", "branch")
	output, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	var branches []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "*") {
			branches = append(branches, line)
		}
	}

	if len(branches) == 0 {
		fmt.Println("No other branches found.")
		return
	}

	fmt.Printf("Other branches: %s\n", strings.Join(branches, ", "))

	cmd = exec.Command("git", "branch", "-d", strings.Join(branches, ", "))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	fmt.Println("Deleted other branches.")
}
