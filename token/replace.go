package token

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	"token-replace/vault"
)

type ReplacementCollection struct {
	extractionCollection *ExtractionCollection
	secretCollection     *vault.SecretCollection
	tokenRegex           *regexp.Regexp
}

func NewReplacementCollection(ec *ExtractionCollection, sc *vault.SecretCollection) *ReplacementCollection {
	return &ReplacementCollection{
		extractionCollection: ec,
		secretCollection:     sc,
		tokenRegex:           regexp.MustCompile(tokenPattern),
	}
}

func (r *ReplacementCollection) ProcessFile(filePath string) (err error) {
	// open the file for reading
	readFile, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprint("Error Opening File ", filePath, " - Error: ", err))
	}
	defer func() {
		if readErr := readFile.Close(); readErr != nil {
			err = errors.New(fmt.Sprint("error closing read file ", filePath, " error: ", readErr))
		}
	}()

	readInfo, err := readFile.Stat()
	if err != nil {
		return errors.New(fmt.Sprint("error fetching file information ", filePath, " - error: ", err))
	}

	// fetch the file's group (we're leaving the owner alone for now) to what the original was
	stat, ok := readInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New(fmt.Sprint("error fetching file group ", filePath, " - error: ", err))
	}

	// open a file for writing
	newFile := filePath + ".token-replace.tmp"
	writeFile, err := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, readInfo.Mode().Perm())
	if err != nil {
		return errors.New(fmt.Sprint("error opening file for writing ", newFile, " - error: ", err))
	}

	defer func() {
		if writeErr := writeFile.Close(); writeErr != nil {
			err = errors.New(fmt.Sprint("error closing write file ", writeFile, " error: ", err))
		}
	}()

	defer func() {
		if err != nil {
			os.Remove(newFile)
		}
	}()

	// iterate over the file contents, and replace the contents
	lineNumber := 0
	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		lineNumber++
		lineText, err := r.ProcessLine(fileScanner.Text())
		if err != nil {
			return errors.New(fmt.Sprint(err, "` in file `", filePath, "` at line ", lineNumber))
		}

		// write out the line in the new file
		if _, err := writeFile.WriteString(fmt.Sprintln(lineText)); err != nil {
			return errors.New(fmt.Sprint("error writing line ", lineNumber, " - error: ", err))
		}
	}

	// set the file ownership of the new file
	if err := os.Chown(newFile, int(stat.Uid), int(stat.Gid)); err != nil {
		return errors.New(fmt.Sprint("error setting ", newFile, " permissions - error: ", err))
	}

	// move this file into position
	if err := os.Rename(newFile, filePath); err != nil {
		return errors.New(fmt.Sprint("error renaming temporary file ", newFile, " - error: ", err))
	}

	return nil
}

func (r *ReplacementCollection) ProcessLine(line string) (string, error) {
	// regex out the tokens in this line
	matches := r.tokenRegex.FindAllString(line, -1)
	for _, match := range matches {
		// resolve this token to a match
		tokenMatch, ok := r.extractionCollection.Tokens[match]
		if !ok {
			return "", errors.New(fmt.Sprint("missing token `", match, "`"))
		}

		// get the token we fetched out of vault
		secret, ok := r.secretCollection.Secrets[tokenMatch.TokenPath]
		if !ok {
			return "", errors.New(fmt.Sprint("missing secret `", tokenMatch.TokenPath, "`"))
		}

		secretValue, err := secret.GetSubKey(tokenMatch.TokenSubKey)
		if err != nil {
			return "", errors.New(fmt.Sprint("missing secret subkey `", tokenMatch.TokenSubKey,
				"` in token `", tokenMatch.TokenPath))
		}

		// replace the value in the line
		line = strings.Replace(line, tokenMatch.TokenName, secretValue, -1)
	}

	return line, nil
}
