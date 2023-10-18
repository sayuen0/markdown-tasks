package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// 対象のディレクトリパスを指定
var TargetDirectory = makeTargetPath()

func main() {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})\.md$`)

	err := filepath.Walk(TargetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		match := re.FindStringSubmatch(info.Name())
		if match != nil {
			year, _ := strconv.Atoi(match[1])
			month, _ := strconv.Atoi(match[2])
			day, _ := strconv.Atoi(match[3])

			t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
			weekday := t.Weekday().String()

			// 曜日の英名から漢字に変換
			var jpWeekday string
			switch weekday {
			case "Monday":
				jpWeekday = "月"
			case "Tuesday":
				jpWeekday = "火"
			case "Wednesday":
				jpWeekday = "水"
			case "Thursday":
				jpWeekday = "木"
			case "Friday":
				jpWeekday = "金"
			case "Saturday":
				jpWeekday = "土"
			case "Sunday":
				jpWeekday = "日"
			}

			newFilename := fmt.Sprintf("%d-%02d-%02d(%s).md", year, month, day, jpWeekday)
			newPath := filepath.Join(filepath.Dir(path), newFilename)

			if err := os.Rename(path, newPath); err != nil {
				return fmt.Errorf("Error renaming file %s to %s: %v", path, newPath, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error processing files:", err)
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
