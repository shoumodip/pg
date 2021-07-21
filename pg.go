package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"strconv"
	"strings"

	"math/rand"
	"time"
)

// The minimum and maximum number of words possible in a password to
// be used if the user didn't provide an absolute value
const (
	MIN_WORDS = 3
	MAX_WORDS = 7
)

// Write error messages to stderr with a leading "error: " string.
// @param condition bool The condition to check
// @param format string The format string
// @param a interface the rest of the arguments to supply to Fprintf()
func displayError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: ")
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}

// Display invalid usage information.
// @param condition bool The condition to check
// @param format string The format string
// @param a interface the rest of the arguments to supply to Fprintf()
func invalidUsage(condition bool, format string, a ...interface{}) {
	if condition {
		displayError(format, a...)
		fmt.Fprintln(os.Stderr, "Usage: pg FILE [WORDS]")
		os.Exit(1)
	}
}

// Check to see if a condition holds true, with accompanying error
// message. Like assert() in C.
// @param condition bool The condition to check
// @param format string The format string
// @param a interface the rest of the arguments to supply to Fprintf()
func checkCond(condition bool, format string, a ...interface{}) {
	if !condition {
		displayError(format, a...)
		os.Exit(1)
	}
}

// Remove empty strings from a slice of words.
// @param words []string The slice of words to filter
// @return []string the filtered slice of words
func filterWords(words []string) []string {
	result := []string{}
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}
	return result
}

// Read a file into a slice of strings separated by newlines.
// @param filePath string the path of the file to read
// @return ([]string, int) the slice of strings, number of lines
func readFile(filePath string) ([]string, int) {
	fileBytes, err := ioutil.ReadFile(filePath)
	checkCond(err == nil, "could not read file: '%s'", filePath)

	fileLines := filterWords(strings.Split(string(fileBytes), "\n"))
	fileSize := len(fileLines)

	checkCond(fileSize != 0, "no words defined in file '%s'", filePath)
	return fileLines, fileSize
}

// Get the number of words to be used in the password, either from
// command line arguments or from psuedo-random generators.
// @param argsCount int The number of command line arguments
// @return int the number of words
func getWordsCount(argsCount int) int {
	if argsCount == 3 {
		wordsCount, err := strconv.Atoi(os.Args[2])
		checkCond(err == nil && wordsCount > 0, "invalid word count: '%s'", os.Args[2])

		return wordsCount
	} else {
		return getRandom(MIN_WORDS, MAX_WORDS)
	}
}

// Generate a pseudo-random number between the two arguments provided
// @return int a random number
func getRandom(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	argsCount := len(os.Args)
	invalidUsage(argsCount < 2, "too few arguments to pg: '%d'", argsCount)
	invalidUsage(argsCount > 3, "too many arguments to pg: '%d'", argsCount)

	filePath := os.Args[1]
	fileLines, fileSize := readFile(filePath)

	// Seed the randomizer
	rand.Seed(time.Now().UnixNano())

	wordsCount := getWordsCount(argsCount)
	password := []string{}

	// Generate the password
	for i := 0; i < wordsCount; i++ {
		password = append(password, fileLines[getRandom(0, fileSize)])
	}

	// Title-case the first word to further increase entropy.
	// No check is being done if the first word even exists, because
	// it is already handled in readFile().
	strings.Title(password[0])

	fmt.Print(strings.Join(password, "_"))

	// Print a symbol ('@') and a random number because why not?
	fmt.Println("@" + strconv.Itoa(getRandom(69, 420)))
}
