package main

import (
    "bytes"
    "encoding/json"
    "io"
    // "log"
    "mime/multipart"
    "os"
    "os/exec"
    "strings"
    "testing"
)

func LoadCorpus(f *testing.F) {
        files, _ := os.ReadDir("corpus/")
        for _, file := range files {
                if data, err := os.ReadFile("corpus/" + file.Name()); err == nil {
                        f.Add(data, []byte{})
                }
        }
}

func FuzzMultipartParser(f *testing.F) {
    f.Add([]byte("--RubyBoundary\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n--RubyBoundary--\r\n"))
    LoadCorpus(f)
    f.Fuzz(func(t *testing.T, data []byte) {
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
                if strings.Contains(err.Error(), "NextPart: EOF") {
                    t.Logf("poopoo")
                    break
                }
                // t.Errorf("err: %s", err.Error())
                return // ignore on parse error
            }

            name := part.FormName()
            filename := part.FileName()
            content, _ := io.ReadAll(part)
            t.Logf("Filename: %s\n", filename)
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
        // t.Errorf("feewffew")
        // 2. Write to file
	
	if err := os.WriteFile("oof.bin", data, 0644); err != nil {
            t.Fatalf("write error: %v", err)
        }

        // 3. Run Ruby subprocess
        out, err := exec.Command("ruby", "rack_parse.rb").Output()
        if err != nil {
            return
            // t.Fatalf("ruby err: %v", err)
        }

        var rubyOutput map[string]any
        if err := json.Unmarshal(out, &rubyOutput); err != nil {
            t.Logf("Here is the thing: %s\n", out)
            t.Fatalf("ruby json error: %v", err)
        }
        // panic("fuck")
        // 4. Compare
        // t.Errorf("RUBY RAW: %s", string(out))

        // t.Errorf("Here is the stuff: %s\n", rubyOutput)
        if rubyFiles, ok := rubyOutput["files"].(map[string]interface{}); ok {
            // t.Errorf("Here is the thing: %s\n", rubyOutput)
            // panic("fe")
            for k, fileInfo := range rubyFiles {
                fileMap := fileInfo.(map[string]interface{})
                goFile, ok := goFiles[k]
                if !ok {
                    t.Fatalf("Go missing file for key %q", k)
                }

                // Compare filename
                if goFile.Filename != fileMap["filename"] {
                    t.Fatalf("Filename mismatch for %q: Go=%q Ruby=%q", k, goFile.Filename, fileMap["filename"])
                }

                // Compare content
                if goFile.Content != fileMap["content"] {
                    t.Fatalf("Content mismatch for %q:\nGo: %q\nRuby: %q", k, goFile.Content, fileMap["content"])
                }
            }
        }

        // 5. Check that the goFiles do not have extra files not in rubyfiles.
        for k := range goFiles {
            if _, ok := rubyOutput["files"].(map[string]interface{})[k]; !ok {
                t.Fatalf("Ruby missing file for key %q", k)
            }
        }
	
    })
}
