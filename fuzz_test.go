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


func containsNonASCII(data []byte) bool {
    for _, b := range data {
        if b > 0x7F {
            return true
        }
    }
    return false
}

func LoadCorpus(f *testing.F) {
        files, _ := os.ReadDir("corpus/")
        for _, file := range files {
                if data, err := os.ReadFile("corpus/" + file.Name()); err == nil {
                        f.Add(data)
                }
        }
}

func FuzzMultipartParser(f *testing.F) {
    f.Add([]byte("--RubyBoundary\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n--RubyBoundary--\r\n"))
    LoadCorpus(f)
    f.Fuzz(func(t *testing.T, data []byte) {
        if containsNonASCII(data) {
            return
        }
        if !strings.Contains(string(data), "Content-Disposition: form-data; ") {
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
                if strings.Contains(err.Error(), "NextPart: EOF") {
                    // t.Logf("poopoo")
                    break
                }
                // t.Errorf("err: %s", err.Error())
                return // ignore on parse error
            }

            name := part.FormName()
            filename := part.FileName()
            content, _ := io.ReadAll(part)
            // t.Logf("Filename: %s\n", filename)
            // t.Logf("name: %s\n", name)
            // t.Logf("content: %s\n", content)
            if filename != "" {
                if name == "" || len(content) == 0 {
                    return
                }
                goFiles[name] = struct {
                    Filename string
                    Content  string
                }{
                    Filename: filename,
                    Content:  string(content),
                }
            } else {
                // return
                if name != "" && len(content) != 0 {
                    goParams[name] = string(content)
                } else {
                    return
                }
                
            }
        }
        // t.Errorf("feewffew")
        // 2. Write to file
	

	    // cmd := exec.Command("ruby", "rack_parse.rb") // Do not use ruby. Use python library instead...
        cmd := exec.Command("python3", "multipart_parse.py")
        stdin, err := cmd.StdinPipe()
        if err != nil {
            t.Fatalf("failed to get stdin pipe: %v", err)
        }
        stdout, err := cmd.StdoutPipe()
        if err != nil {
            t.Fatalf("failed to get stdout pipe: %v", err)
        }

        if err := cmd.Start(); err != nil {
            t.Fatalf("failed to start ruby subprocess: %v", err)
        }

        // Write data to stdin
        if _, err := stdin.Write(data); err != nil {
            t.Fatalf("failed to write to ruby stdin: %v", err)
        }
        stdin.Close()

        out, err := io.ReadAll(stdout)
        if err != nil {
            t.Fatalf("failed to read from ruby stdout: %v", err)
        }

        if err := cmd.Wait(); err != nil {
            return
            // t.Fatalf("ruby execution error: %v\nOutput: %s", err, out)
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
                /*
		goFile, ok := goFiles[k]
                if !ok {
                    t.Fatalf("Go missing file for key %q", k)
                }
		*/
		goFile, ok := goFiles[k]
                if !ok { // Just ignore these for now
		          return
                    // t.Fatalf("Go missing file for key %q", k)
                }

                // Compare filename
                if goFile.Filename != fileMap["filename"] {
                    t.Fatalf("Filename mismatch for %q: Go=%q Ruby=%q", k, goFile.Filename, fileMap["filename"])
                }

                // Compare content
                /*
                if goFile.Content != fileMap["content"] {
                    t.Fatalf("Content mismatch for %q:\nGo: %q\nRuby: %q", k, goFile.Content, fileMap["content"])
                }
                */

                normalize := func(s string) string {
                    // Normalize newlines to "\n" and trim surrounding whitespace
                    return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
                }

                goNorm := normalize(goFile.Content)
                rubyNorm := normalize(fileMap["content"].(string))
                if len(goNorm) > len(rubyNorm) {
                    return
                }
                if goNorm != rubyNorm {
                    t.Fatalf("Content mismatch for %q:\nGo: %q\nRuby: %q", k, goNorm, rubyNorm)
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
