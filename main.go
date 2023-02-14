package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

var usage = `Usage: maki [OPTIONS]

	Run make targets in interactive mode.

Options:
-f            The makefile name.
-h, --help    Show this message and exit.
`

func main() {
	var makefile string

	flag.StringVar(&makefile, "f", "Makefile", "The makefile name.")
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if _, err := os.Stat(makefile); err != nil {
		fmt.Printf("File named %q not found.", makefile)
		return
	}

	makeOut, err := exec.Command("make", "-f", makefile, "-pRrq", ":").Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Ignore error, make returns 1 when using -q and some targets are not up to date
			// This is irrelevant for this code's purposes
		} else {
			fmt.Println("Error executing command:", err)
			return
		}
	}

	targets, err := getTargetsFromMakeDb(makeOut)
	fmt.Print(targets)
	if err != nil {
		fmt.Println("Error executing command:", err)
	}

	prompt := promptui.Select{
		Label: "Select make target:",
		Items: targets,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	cmd := exec.Command("make", "-f", makefile, result)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing target %q command: %s", result, err)
		return
	}
}

func getTargetsFromMakeDb(makeDb []byte) ([]string, error) {
	var targets []string

	// Convert the []byte to a string
	text := string(makeDb)

	// Retrieve the makefile_list variable
	re := regexp.MustCompile(`MAKEFILE_LIST\s*:=\s*(.+)`)
	makefileList := re.FindStringSubmatch(text)[1]

	// Replace all blank lines with an empty string
	re = regexp.MustCompile(`(?m)^\s*\n`)
	textNoNewLines := re.ReplaceAllString(text, "")

	// Limit the parsing range to the block that contains the targets
	re = regexp.MustCompile(`(^|\n)# Files(\n|$)[\s\S]*?(^|\n)# Finished Make data base`)
	match := re.FindString(textNoNewLines)
	if match == "" {
		return targets, errors.New("error during the parsing of make database")
	}

	// Ignore non-targets and special targets
	re = regexp.MustCompile(`(?m)^[#|\.|:].*\n`)
	resultNoNonTargets := re.ReplaceAllString(match, "")

	re = regexp.MustCompile(`(?m)^\t.*\n`)
	// Replace all lines starting with a tab character with an empty string
	resultNoCommands := re.ReplaceAllString(resultNoNonTargets, "")

	// Remove the 'Makefile' entry
	regex := fmt.Sprintf(`(?m)^(%s).*\n`, makefileList)
	re = regexp.MustCompile(regex)
	resultNoMakefile := re.ReplaceAllString(resultNoCommands, "")

	lines := strings.Split(resultNoMakefile, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, ":")
	}

	// Remove the residual lines that are not targets
	for _, v := range lines {
		if strings.HasPrefix(v, "#") || v == "" {
			continue
		}
		targets = append(targets, v)
	}

	sort.Strings(targets)

	return targets, nil
}
