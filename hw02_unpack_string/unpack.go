package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const backslash rune = 92

func Unpack(str string) (string, error) { //nolint:gocognit // complexity 46. How can we reduce it?
	var (
		err          error
		resultWord   strings.Builder
		lastRune     rune
		printNum     bool
		unpackSlash  bool
		repeatNumber int
	)
	if len(str) == 0 {
		return "", nil
	}

	// convert str to unicode runes
	for i, v := range str {
		if i == 0 && unicode.IsDigit(v) {
			return "", ErrInvalidString
		}

		// Checking printable for printable lastRune or should we unpack slash
		if !unicode.IsPrint(lastRune) && unicode.IsDigit(v) && unpackSlash {
			repeatNumber, err = strconv.Atoi(string(v))
			if err != nil {
				return "", err
			}

			resultWord.WriteString(strings.Repeat(string(backslash), repeatNumber-1))
			unpackSlash = false
			lastRune = 0
			continue
		} else if !unicode.IsPrint(lastRune) {
			lastRune = v
			continue
		}

		// Checking for double numbers in a row with printNum flag
		if unicode.IsDigit(v) && unicode.IsDigit(lastRune) && printNum {
			repeatNumber, err = strconv.Atoi(string(v))
			if err != nil {
				return "", err
			}

			resultWord.WriteString(strings.Repeat(string(lastRune), repeatNumber))
			printNum = false
			lastRune = 0
			continue
		} else if unicode.IsDigit(v) && unicode.IsDigit(lastRune) && !printNum {
			return "", ErrInvalidString
		}

		// Processing number
		if unicode.IsDigit(v) {
			printNum = false
			if lastRune == backslash {
				printNum = true
				lastRune = v
				continue
			}

			repeatNumber, err = strconv.Atoi(string(v))
			if err != nil {
				return "", err
			}

			resultWord.WriteString(strings.Repeat(string(lastRune), repeatNumber))
			lastRune = v
			continue
		}

		// Processing double backslash
		if v == backslash && lastRune == backslash {
			resultWord.WriteRune(v)
			unpackSlash = true
			lastRune = 0
			continue
		}

		// Scan into results if lastRune is not a number with printNum flag
		if !unicode.IsDigit(lastRune) || printNum {
			resultWord.WriteRune(lastRune)
		}
		lastRune = v
	}

	if lastRune != 0 {
		resultWord.WriteRune(lastRune)
	}

	return resultWord.String(), nil
}
