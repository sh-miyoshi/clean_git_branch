package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func runCommand(name string, args ...string) string {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("%v\n%s", err, strings.TrimSpace(string(output))))
	}

	return string(output)
}

func main() {
	output := runCommand("git", "remote", "show", "origin")

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

	runCommand("git", "checkout", branch)

	output = runCommand("git", "branch")

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

	deleteArgs := append([]string{"branch", "-d"}, branches...)
	runCommand("git", deleteArgs...)

	fmt.Println("Deleted other branches.")
}
