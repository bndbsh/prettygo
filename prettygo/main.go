package main 

import (
	"fmt"
	"os"
	"bytes"
	"io"
	"strings"
	pattern "github.com/jnwhiteh/go-luapatterns/pattern"
)

const (
	Black = "\x1b[30m"
	Red = "\x1b[31m"
	Green = "\x1b[32m"
	Yellow = "\x1b[33m"
	Blue = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan = "\x1b[36m"
	White = "\x1b[37m"
)

func main() {
	rules := map[string]string{
		"^\n" : "",
		"^rm.*" : "",
		"^cp.*" : "",
		"^make.*" : "",
		"^gotest.*" : "",
		"^gopack.*" : "",
		"^[856]l.*" : "",
		"^[856]g %-o _.+_%.[856] (.+)" : Green + "Compiling:" + White + " %1",
		"^%w+%.go:%d+:.+" : Red + "%0" + White,
	}
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Printf("Init failed: %s", err)
		return
	}
	pid, err := os.ForkExec("/usr/bin/env", []string{"/usr/bin/env", "gomake"}, os.Environ(), "", []*os.File{nil, w, w})

	if err != nil {
		fmt.Printf("Fork failed: %s\n", err)
		return
	}
	w.Close()

	var b bytes.Buffer
	io.Copy(&b, r)
	os.Wait(pid, 0)

	split := strings.Split(b.String(), "\n", -1)

	for _, line := range split {
		for pat, sub := range rules {
			if b, _ := pattern.Match(line, pat); b {
				line, _ = pattern.Replace(line, pat, sub, -1)
				break
			}
		}
		if len(line) > 0 {
			fmt.Printf("%s\n", line)
		}
	}
	fmt.Printf("Done!\n")
}

