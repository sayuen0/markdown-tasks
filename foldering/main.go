package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// 日記を月ごとフォルダ分けして整理する
// ./2022-01-01.md → ./2022/1/2022-01-01.md

// 対象のディレクトリパスを指定
var TargetDirectory = makeTargetPath()

func main() {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})(?:\((月|火|水|木|金|土|日)\))?\.md$`)

	// ディレクトリ内のすべてのファイルとフォルダをリスト
	files, err := os.ReadDir(makeTargetPath())
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		match := re.FindStringSubmatch(file.Name())
		if match != nil {
			year := match[1]
			month, err := strconv.Atoi(match[2])
			if err != nil {
				fmt.Printf("Error converting month for %s: %v\n", file.Name(), err)
				continue
			}

			// 新しいディレクトリパスを生成
			newDir := filepath.Join(TargetDirectory, year, strconv.Itoa(month))
			if _, err := os.Stat(newDir); os.IsNotExist(err) {
				// 年と月のディレクトリを作成
				if err := os.MkdirAll(newDir, 0755); err != nil {
					fmt.Printf("Error creating directory %s: %v\n", newDir, err)
					continue
				}
			}

			// ファイルを新しいディレクトリに移動
			oldPath := filepath.Join(TargetDirectory, file.Name())
			newPath := filepath.Join(newDir, file.Name())
			if err := os.Rename(oldPath, newPath); err != nil {
				fmt.Printf("Error moving file %s to %s: %v\n", oldPath, newPath, err)
			}
		}
	}
}

func makeTargetPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	root := filepath.Join(homedir, "markdowns")
	return filepath.Join(root, "private", "diary")
}
