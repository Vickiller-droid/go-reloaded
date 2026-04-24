package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input.txt> <output.txt>")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var processedLines []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			processedLines = append(processedLines, "")
			continue
		}

		words := strings.Fields(line)

		for i := 0; i < len(words); i++ {
			tag := words[i]

			switch {
			case tag == "(hex)":
				if i > 0 {
					val, _ := strconv.ParseInt(words[i-1], 16, 64)
					words[i-1] = fmt.Sprint(val)
					words[i] = ""
				}

			case tag == "(bin)":
				if i > 0 {
					val, _ := strconv.ParseInt(words[i-1], 2, 64)
					words[i-1] = fmt.Sprint(val)
					words[i] = ""
				}

			case tag == "(up)":
				if i > 0 {
					words[i-1] = strings.ToUpper(words[i-1])
					words[i] = ""
				}

			case tag == "(low)":
				if i > 0 {
					words[i-1] = strings.ToLower(words[i-1])
					words[i] = ""
				}

			case tag == "(cap)":
				if i > 0 {
					words[i-1] = capitalize(words[i-1])
					words[i] = ""
				}

			case strings.HasPrefix(tag, "(up,") || strings.HasPrefix(tag, "(low,") || strings.HasPrefix(tag, "(cap,"):
				numStr := strings.Trim(words[i+1], " )")
				count, _ := strconv.Atoi(numStr)

				for j := 1; j <= count; j++ {
					if i-j >= 0 {
						if strings.Contains(tag, "up") {
							words[i-j] = strings.ToUpper(words[i-j])
						} else if strings.Contains(tag, "low") {
							words[i-j] = strings.ToLower(words[i-j])
						} else if strings.Contains(tag, "cap") {
							words[i-j] = capitalize(words[i-j])
						}
					}
				}
				words[i], words[i+1] = "", ""
			}
		}

		var filtered []string
		for _, w := range words {
			if w != "" {
				filtered = append(filtered, w)
			}
		}

		for i := 0; i < len(filtered)-1; i++ {
			word := strings.ToLower(filtered[i])
			if word == "a" {
				next := strings.ToLower(filtered[i+1])
				if len(next) > 0 && strings.ContainsAny(string(next[0]), "aeiouh") {
					if filtered[i] == "A" {
						filtered[i] = "An"
					} else {
						filtered[i] = "an"
					}
				}
			}
		}

		currentLine := strings.Join(filtered, " ")

		rePunc := regexp.MustCompile(`\s+([.,!?:;]+)`)
		currentLine = rePunc.ReplaceAllString(currentLine, "$1")

		reSpace := regexp.MustCompile(`([.,!?:;]+)([a-zA-Z0-9])`)
		currentLine = reSpace.ReplaceAllString(currentLine, "$1 $2")

		currentLine = handleQuotes(currentLine)
		processedLines = append(processedLines, currentLine)
	}

	finalOutput := strings.Join(processedLines, "\n")
	os.WriteFile(os.Args[2], []byte(strings.TrimSpace(finalOutput)+"\n"), 0644)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func handleQuotes(s string) string {
	words := strings.Fields(s)
	isOp := true

	for i := 0; i < len(words); i++ {
		if words[i] == "'" {
			if isOp {
				if i+1 < len(words) {
					words[i+1] = "'" + words[i+1]
					words[i] = ""
					isOp = false
				}
			} else {
				if i-1 >= 0 {
					words[i-1] += "'"
					words[i] = ""
					isOp = true
				}
			}
		}
	}

	return strings.Join(strings.Fields(strings.Join(words, " ")), " ")
}
