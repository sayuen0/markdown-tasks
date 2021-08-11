package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
)

//const READDIR = "./private/diary"

const (
	READDIR  = "docs"
	WRITEDIR = "dist"
	LF       = "\n"
)

var headings = map[int]string{
	1: "# ",
	2: "## ",
	3: "### ",
	4: "#### ",
	5: "##### ",
	6: "######",
}

var promptSuggests = []string{
	"今日のGJ",
	"今日の伸びしろ",
}

func contains(s string, slice []string) bool {
	for _, a := range slice {
		if strings.TrimSpace(s) == a {
			return true
		}
	}
	return false
}

func completer(in prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0, len(promptSuggests))
	for _, suggest := range promptSuggests {
		s = append(s, prompt.Suggest{Text: suggest})
	}
	return prompt.FilterHasSuffix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	// input
	in := prompt.Input(">>>", completer, prompt.OptionTitle("集計項目"))
	fmt.Println(in)
	if !contains(in, promptSuggests) {
		fmt.Println("invalid input")
		return
	}

	// summarize text
	s := headings[1] + in + LF
	s = summarize(s, in)

	// write summary
	newFile, err := os.Create(filepath.Join(WRITEDIR, in + ".md"))
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
				// "今日のGJ"を抽出してファイルに書き込む ただし見出しの代わりに日付のh2が欲しい
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
					// "docs/2021/08/03.md" → "# 2021/08/03
					h, _ := filepath.Rel(READDIR, path)
					h = strings.TrimSuffix(h, ".md")
					s += headings[2] + h + LF
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return s
}

var r = regexp.MustCompile("[0-2][0-9].md")

func isTargetMarkdown(info fs.FileInfo) bool {
	name := info.Name()
	return filepath.Ext(name) == ".md" && r.MatchString(name)
}

func copyFile() {
	f, err := os.Open(filepath.Join(READDIR, "hoge.md"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 書き込み対象ファイル
	newFile, err := os.Create(filepath.Join(READDIR, "fuga.md"))
	defer newFile.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.Replace(s, "2", "4", -1)
		if _, err := newFile.WriteString(s + "\n"); err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
