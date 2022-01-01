package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var weekdayja = strings.NewReplacer(
	"Sun", "日",
	"Mon", "月",
	"Tue", "火",
	"Wed", "水",
	"Thu", "木",
	"Fri", "金",
	"Sat", "土",
)

var BASE_PATH = filepath.Join(GetHomeDir(), "markdowns/private/diary")
var TASK_HEADING = "# 今日のタスク"

func GetHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dirname
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Recover: ", err, string(debug.Stack()))
		}
	}()

	if len(os.Args) > 1 { // 1st arg is self binary name
		flag.Parse()
		basePath := flag.Args()[0]
		if basePath != "" {
			BASE_PATH = basePath
		}
	}

	// read last
	lastPath := getLastDateFile(BASE_PATH)
	if lastPath == "" {
		panic(errors.New("no files"))
	}
	fmt.Printf("last file found: %s\n", lastPath)
	lastContent, err := ioutil.ReadFile(filepath.Join(BASE_PATH, lastPath))
	if err != nil {
		panic(err)
	}

	// mkdir if not exists
	year := time.Now().Year()
	month := time.Now().Month()
	dir := filepath.Join(BASE_PATH, strconv.Itoa(year), strconv.Itoa(int(month)))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	// check today
	today := time.Now()
	filePath := makeFilepathFromDate(today)
	todayPath := filepath.Join(BASE_PATH, filePath)
	exists, err := Exists(todayPath)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Printf("today file exists: %s\n", todayPath)
		return
	}

	// create today's content
	lc := string(lastContent)
	taskStr, err := ExtractTask(lc)
	if err != nil {
		panic(err)
	}
	undoneTasks := RemoveDoneTasks(taskStr)
	todayContent := fmt.Sprintf(DiaryTemplate(today), undoneTasks)
	if err := ioutil.WriteFile(todayPath, []byte(todayContent), os.ModePerm); err != nil {
		panic(err)
	}
	fmt.Printf("created new file: %s \n", todayPath)
}

func makeFilepathFromDate(date time.Time) string {
	year := date.Year()
	month := date.Month()

	day := date.Day()
	return fmt.Sprintf("%s.md", filepath.Join(strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day)))
}

var fileRegex = regexp.MustCompile(`\d{4}/\d{1,2}/\d{1,2}\.md`)

func getLastDateFile(basePath string) string {
	lastPath := ""
	lastDate := time.Unix(0, 0).UTC()
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		// skip invalid format file name
		if !fileRegex.MatchString(rel) {
			return nil
		}

		stem := strings.TrimSuffix(rel, filepath.Ext(rel))
		fileDate, err := time.Parse("2006/1/2", stem)
		if err != nil {
			return err
		}

		if fileDate.IsZero() {
			return nil
		}

		if fileDate.After(lastDate) {
			lastDate = fileDate
			lastPath = stem
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s.md", lastPath)
}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func RemoveDoneTasks(tasks string) string {
	doneTaskSymbol := "- [x]"
	tasksArr := strings.Split(tasks, "\n")
	undoneTasks := make([]string, 0, len(tasksArr))
	for _, t := range tasksArr {
		if strings.Index(strings.TrimSpace(t), doneTaskSymbol) == -1 {
			undoneTasks = append(undoneTasks, t)
		}
	}
	return strings.Join(undoneTasks, "\n")
}

func ExtractTask(report string) (string, error) {
	tasks := make([]string, 0)
	reading := false
	for _, line := range strings.Split(strings.TrimSuffix(report, "\n"), "\n") {
		if reading {
			if strings.Index(strings.TrimSpace(line), "# ") == 0 {
				break
			}
			tasks = append(tasks, line)
		}
		if !reading && strings.TrimSpace(line) == TASK_HEADING {
			reading = true
		}
	}
	if !reading {
		return "", errors.New("no tasks found")
	}
	return strings.Join(tasks, "\n"), nil
}

func DiaryTemplate(date time.Time) string {
	dateF := weekdayja.Replace(date.Format("2006年01月02日(Mon)"))

	return fmt.Sprintf(`%s
%s
%s`, dateF, TASK_HEADING, `%s
# 作業ログ
# 呟き場
## 今日の発見
## 今日の感謝
## 今日のGJ
## 今日の伸びしろ 
`)
}
