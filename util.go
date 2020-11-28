package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Reads an integer from standard input. On error, the returned int will be 0,
// and the error will be returned the caller along with the actual read string.
func readInt() (int, error, string) {
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return 0, err, in
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if err != nil {
		return 0, err, in
	}
	num, err := strconv.Atoi(in)
	if err != nil {
		return num, err, in
	}
	return num, nil, in
}

// Displays the given promp and calls readInt() to read an integer from standard
// input. On error, the returned int is 0 and the error will be returned the
// caller along with the actual read string.
func promptInt(prompt string) (int, error, string) {
	fmt.Print(prompt)
	return readInt()
}
