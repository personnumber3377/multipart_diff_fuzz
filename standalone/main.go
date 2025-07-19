package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	// "strings"
)

func containsNonASCII(data []byte) bool {
	for _, b := range data {
		if b > 0x7F {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./standalone_parser <input_file>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	if containsNonASCII(data) {
		fmt.Println("Input contains non-ASCII characters, skipping.")
		return
	}

	// 1. Parse with Go
	reader := multipart.NewReader(bytes.NewReader(data), "RubyBoundary")
	goParams := map[string]string{}
	goFiles := map[string]struct {
		Filename string
		Content  string
	}{}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Go parser error: %v\n", err)
			os.Exit(1)
		}

		name := part.FormName()
		filename := part.FileName()
		content, _ := io.ReadAll(part)
		fmt.Printf("Here is the name: %s\n", name)
		fmt.Printf("Here is the filename: %s\n", filename)
		fmt.Printf("content: %s\n", content)
		if filename != "" {
			goFiles[name] = struct {
				Filename string
				Content  string
			}{
				Filename: filename,
				Content:  string(content),
			}
		} else {
			goParams[name] = string(content)
		}
	}
	fmt.Printf("Here is the thing: %s", goFiles)
	// 2. Run Ruby subprocess with input from stdin
	rubyCmd := exec.Command("ruby", "rack_parse.rb")
	stdin, err := rubyCmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get stdin: %v\n", err)
		os.Exit(1)
	}
	stdout, err := rubyCmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get stdout: %v\n", err)
		os.Exit(1)
	}
	if err := rubyCmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Ruby start error: %v\n", err)
		os.Exit(1)
	}

	stdin.Write(data)
	stdin.Close()
	out, err := io.ReadAll(stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read Ruby output: %v\n", err)
		os.Exit(1)
	}
	rubyCmd.Wait()

	// 3. Parse and compare
	var rubyOutput map[string]any
	if err := json.Unmarshal(out, &rubyOutput); err != nil {
		fmt.Fprintf(os.Stderr, "JSON parse error: %v\nOutput: %s\n", err, out)
		os.Exit(1)
	}

	// Compare Go vs Ruby
	if rubyFiles, ok := rubyOutput["files"].(map[string]interface{}); ok {
		for k, fileInfo := range rubyFiles {
			fileMap := fileInfo.(map[string]interface{})
			goFile, ok := goFiles[k]
			if !ok {
				fmt.Printf("Go missing file for key %q\n", k)
				os.Exit(1)
			}
			if goFile.Filename != fileMap["filename"] {
				fmt.Printf("Filename mismatch for %q: Go=%q Ruby=%q\n", k, goFile.Filename, fileMap["filename"])
				os.Exit(1)
			}
			if goFile.Content != fileMap["content"] {
				fmt.Printf("Content mismatch for %q:\nGo: %q\nRuby: %q\n", k, goFile.Content, fileMap["content"])
				os.Exit(1)
			}
		}
	}
	for k := range goFiles {
		if _, ok := rubyOutput["files"].(map[string]interface{})[k]; !ok {
			fmt.Printf("Ruby missing file for key %q\n", k)
			os.Exit(1)
		}
	}

	fmt.Println("✅ Go and Ruby multipart parsers agree!")
}


