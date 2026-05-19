package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/pflag"
)

func runCommand(name string, args ...string) string {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("%v\n%s", err, strings.TrimSpace(string(output))))
	}

	return string(output)
}

func main() {
	var remoteName string
	var interactive, forceDelete bool
	pflag.StringVarP(&remoteName, "name", "n", "origin", "remote name")
	pflag.BoolVarP(&interactive, "interactive", "i", false, "interactive mode, ask for confirmation before deleting branchs")
	pflag.BoolVarP(&forceDelete, "force", "D", false, "force delete branchs")

	pflag.Parse()

	output := runCommand("git", "remote", "show", remoteName)

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

	if interactive {
		var confirm string
		fmt.Print("Are you sure you want to delete these branches? (y/N): ")
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Aborting.")
			return
		}
	}

	del := "-d"
	if forceDelete {
		del = "-D"
	}

	deleteArgs := append([]string{"branch", del}, branches...)
	runCommand("git", deleteArgs...)

	fmt.Println("Deleted other branches.")
}
