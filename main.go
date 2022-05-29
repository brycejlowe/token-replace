package main

import (
	"log"
	"os"
	"path/filepath"
	"token-replace/token"
	"token-replace/vault"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Files []string `long:"file" short:"f" description:"File to Convert" required:"true"`

	Server string `long:"server" short:"s" description:"Vault Server" required:"true"`
	Token  string `long:"token" short:"t" description:"Authentication Token" required:"true"`

	Implode bool `long:"implode" short:"i" description:"Delete Self"`
}

func init() {
	// log to std out
	log.SetOutput(os.Stdout)
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalln("Error Parsing Flags")
	}

	// extract the tokens from the token requested
	extractedTokens := token.NewExtractionCollection()
	for i, file := range opts.Files {
		if file, err := filepath.Abs(file); err != nil {
			log.Fatalln("Unable to Determine Absolute File Path for", file)
		}

		opts.Files[i] = file

		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalln("Invalid File Specified", file)
		}

		// extract the tokens from the file
		log.Println("Extracting Tokens from File", file)
		if err := extractedTokens.ParseFile(file); err != nil {
			log.Fatalln("Error Extracting Tokens from", file, ":", err)
		}
	}

	// init the connection to vault
	vaultClient := vault.NewClient(opts.Server, opts.Token)

	// use the vault api to resolve fetch the secrets, we're de-duping calls here as vault stores secrets in paths
	// and we are storing per-package secrets in the same path.
	secretCollection := vault.NewSecretCollection()
	for tokenKey := range extractedTokens.Tokens {
		tokenPath := &extractedTokens.Tokens[tokenKey].TokenPath
		if secretCollection.HasSecret(*tokenPath) {
			// we already have this path fetched from vault, don't do it again
			continue
		}

		// do the fetch
		log.Println("Fetching Vault Secret:", *tokenPath)
		secretContents, err := vault.SecretFetch(*tokenPath, vaultClient)
		if err != nil {
			log.Fatalln("Error Fetching Secret:", err)
		}

		// log the request id, for reasons...
		log.Println("Successfully Fetch Secret with Request ID", secretContents.RequestId)
		secretCollection.AddSecret(*tokenPath, secretContents)
	}

	replacementCollection := token.NewReplacementCollection(extractedTokens, secretCollection)
	for _, file := range opts.Files {
		log.Println("Replacing Tokens in File", file)
		if err := replacementCollection.ProcessFile(file); err != nil {
			log.Fatalln("Error Processing File:", file, "- Error:", err)
		}
	}

	defer func() {
		if opts.Implode {
			log.Println("Imploding Executable")
			implode()
		}
	}()

	os.Exit(0)
}

func implode() {
	exePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalln("Error Detecting Executable Path", err)
	}

	if err := os.Remove(exePath); err != nil {
		log.Fatalln("Error Removing Executable", err)
	}
}
