package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//const READDIR = "./private/diary"

const (
	READDIR = "docs"
	WRITEDIR = "dist"
	LF      = "\n"
)


var headingPrefixes = map[string]string{
	"H1": "# ",
	"H2": "## ",
	"H3": "### ",
	"H4": "#### ",
	"H5": "##### ",
	"H6": "######",
}

func main() {
	// 指定ディレクトリ以下のファイルを開いてテキストを全部保存をする
	//buf := make([]byte, 0)
	const heading = "今日のGJ"
	s := headingPrefixes["H1"] + heading + LF
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
				if strings.HasPrefix(line, headingPrefixes["H1"]) ||
					strings.HasPrefix(line, headingPrefixes["H2"]) {
					collecting = false
				}
				if collecting {
					s += line + LF
				}
				if line == headingPrefixes["H2"] + heading {
					collecting = true
					// "docs/2021/08/03.md" → "# 2021/08/03
					h, _ := filepath.Rel(READDIR, path)
					h = strings.TrimSuffix(h, ".md")
					s += headingPrefixes["H2"] + h + LF
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	newFile, err := os.Create(filepath.Join(WRITEDIR, "fuga.md"))
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	if _, err := newFile.WriteString(s); err != nil {
		panic(err)
	}
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
