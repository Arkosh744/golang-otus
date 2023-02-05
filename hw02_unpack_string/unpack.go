package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const backslash rune = 92

func Unpack(str string) (string, error) { //nolint:gocognit // complexity 46. How can we reduce it?
	var (
		err         error
		resultWord  strings.Builder
		lastRune    rune
		printNum    bool
		unpackSlash bool
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
			_, err = fmt.Fprintf(&resultWord, "%s", strings.Repeat(string(backslash), int(v-'0')-1))
			if err != nil {
				return resultWord.String(), err
			}
			unpackSlash = false
			lastRune = 0
			continue
		} else if !unicode.IsPrint(lastRune) {
			lastRune = v
			continue
		}

		// Checking for double numbers in a row with printNum flag
		if unicode.IsDigit(v) && unicode.IsDigit(lastRune) && printNum {
			_, err = fmt.Fprintf(&resultWord, "%s", strings.Repeat(string(lastRune), int(v-'0')))
			if err != nil {
				return resultWord.String(), err
			}
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

			// Subtract 48 ('0') from the rune to get the int number
			_, err = fmt.Fprintf(&resultWord, "%s", strings.Repeat(string(lastRune), int(v-'0')))
			if err != nil {
				return resultWord.String(), err
			}
			lastRune = v
			continue
		}

		// Processing double backslash
		if v == backslash && lastRune == backslash {
			_, err = fmt.Fprintf(&resultWord, "%s", string(v))
			if err != nil {
				return resultWord.String(), err
			}
			unpackSlash = true
			lastRune = 0
			continue
		}

		// Scan into results if lastRune is not a number with printNum flag
		if !unicode.IsDigit(lastRune) || printNum {
			_, err = fmt.Fprintf(&resultWord, "%s", string(lastRune))
			if err != nil {
				return resultWord.String(), err
			}
		}
		lastRune = v
	}

	if lastRune != 0 {
		_, err = fmt.Fprintf(&resultWord, "%s", string(lastRune))
		if err != nil {
			return resultWord.String(), err
		}
	}

	return resultWord.String(), nil
}
