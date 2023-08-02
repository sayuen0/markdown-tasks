package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

var (
	ROOT, READDIR, WRITEDIR string
	LF                      = "\n"
)

var headings = map[int]string{
	1: "# ",
	2: "## ",
	3: "### ",
	4: "#### ",
	5: "##### ",
	6: "######",
}

//go:embed promptSuggests.json
var promptSuggests []byte

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fmt.Println(homedir)
	ROOT = filepath.Join(homedir, "markdowns")
	READDIR = filepath.Join(ROOT, "private", "diary")
	WRITEDIR = filepath.Join(ROOT, "private", "diary", "dist")

	// input
	var suggests []string
	if err := json.Unmarshal(promptSuggests, &suggests); err != nil {
		panic(err)
	}
	prompt := promptui.Select{
		Label: "集計項目を選択",
		Items: suggests,
	}

	_, in, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	// summarize text
	s := headings[1] + in + LF
	s = summarize(s, in)

	// write summary
	timestamp := time.Now().Format("20060102150405")
	newFilePath := filepath.Join(WRITEDIR, timestamp+"_"+in+".md")
	fmt.Println(newFilePath)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		panic(err)
	}

	defer newFile.Close()
	if _, err := newFile.WriteString(s); err != nil {
		panic(err)
	}
}

func summarize(s string, heading string) string {
	if err := filepath.Walk(READDIR, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// マークダウンファイルかつ指定のフォーマットを満たしている場合だけ開く
		if isTargetMarkdown(info) {
			fmt.Println(info.Name(), path)
			// ターゲットのファイルなので、読みだして溜め込む
			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(f)
			collecting := false
			for scanner.Scan() {
				line := scanner.Text()
				// "今日のGJ"を抽出してファイルに書き込む ただし見出しは日付h2で置換
				// 今日のGJでないh2見出しがきたら集計停止
				if strings.HasPrefix(line, headings[1]) ||
					strings.HasPrefix(line, headings[2]) {
					collecting = false
				}
				if collecting {
					s += line + LF
				}
				if line == headings[2]+heading {
					collecting = true
					// "docs/2021/08/03.md" → "## 2021-08-03(火)"
					// "docs/2021-08-03.md" → "## 2021-08-03(火)"
					h, _ := filepath.Rel(READDIR, path)
					h = strings.TrimSuffix(h, ".md")
					week, _ := getDayOfWeekChar(h, "2006-01-02")
					h += fmt.Sprintf("(%s)", week)
					s += headings[2] + h + LF
				}
			}
		}
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	return s
}

var r = regexp.MustCompile("[0-2][0-9].md")

func isTargetMarkdown(info fs.FileInfo) bool {
	name := info.Name()
	return filepath.Ext(name) == ".md" && r.MatchString(name)
}
