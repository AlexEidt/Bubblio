package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseFont(filename string) map[rune]*symbol {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	pattern, _ := regexp.Compile("0x[0-9|a-f]{2}")

	symbols := make(map[rune]*symbol)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			var bitmap []string
			ascii, _ := strconv.ParseInt(line[:4], 0, 16)
			scanner.Scan()
			line = scanner.Text()
			for strings.HasPrefix(line, " ") && scanner.Scan() {
				bitmap = append(bitmap, strings.Trim(line, " "))
				line = scanner.Text()
			}
			symbols[rune(ascii)] = NewSymbol(bitmap, '@')
		}
	}
	return symbols
}
