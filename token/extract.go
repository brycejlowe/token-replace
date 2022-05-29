package token

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

type ExtractionCollection struct {
	Tokens map[string]*Token
	parsedFiles map[string]bool
	tokenRegex *regexp.Regexp
}

func NewExtractionCollection() *ExtractionCollection {
	return &ExtractionCollection{
		Tokens: make(map[string]*Token),
		parsedFiles: make(map[string]bool),
		tokenRegex: regexp.MustCompile(tokenPattern),
	}
}

func (t *ExtractionCollection) ParseFile(filePath string) error {
	// open the file and start parsing
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// close when we're done
	defer file.Close()

	// line tracker
	lineNumber := 0

	// loop through the lines in the file and extract the tokens, it's worth pointing out here that
	// this breaks horribly if a token were to span more than 1 line.  that's a limitation were going
	// to live with for now
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lineNumber++
		if err := t.ParseLine(fileScanner.Text()); err != nil {
			return errors.New(fmt.Sprint("unable to parse token at line ",lineNumber, " error: ", err))
		}
	}

	return nil
}


func (t *ExtractionCollection) ParseLine(fileLine string) error {
	// regex out the tokens in this line
	matches := t.tokenRegex.FindAllString(fileLine, -1)

	// no tokens found
	if len(matches) == 0 {
		return nil
	}

	// process all of the matches
	for _, match := range matches {
		// duplicates just cause more calls to vault, so eliminate them
		if _, ok := t.Tokens[match]; ok {
			continue
		}

		// append it to the final result
		if parsedToken, err := NewToken(match); err == nil {
			t.Tokens[match] = parsedToken
		} else {
			return err
		}
	}

	return nil
}
