package main

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test_makeFilepathFromDate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1桁の数値",
			args: args{
				date: time.Date(2021, 2, 1, 0, 0, 0, 0, loc),
			},
			want: "2021/2/1.md",
		},
		{
			name: "2桁の数値",
			args: args{
				date: time.Date(2021, 12, 10, 0, 0, 0, 0, loc),
			},
			want: "2021/12/10.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeFilepathFromDate(tt.args.date); got != tt.want {
				t.Errorf("makeFilepathFromDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLastDateFile(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "",
			want: "2021/12/10.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLastDateFile(); got != tt.want {
				t.Errorf("getLastDateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDoneTasks(t *testing.T) {
	type args struct {
		tasks string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "多段階の階層に対応",
			args: args{
				tasks: `- [ ] 手を洗う
  - [x] 石鹸を使う
- [x] 歯を磨く
  - [x] 歯磨き粉を使う
- [ ] 顔を洗う
  - [x] 洗顔フォームを使う
  - [ ] タオルで拭く
    - [x] いいタオルだ
    - [ ] 匂いを嗅ぐ
`,
			},
			want: `- [ ] 手を洗う
- [ ] 顔を洗う
  - [ ] タオルで拭く
    - [ ] 匂いを嗅ぐ
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDoneTasks(tt.args.tasks); got != tt.want {
				t.Errorf("RemoveDoneTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTask(t *testing.T) {
	type args struct {
		report string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "先頭に余計なものが付いてたり、タスクの内容取れるか",
			args: args{
				report: `本日は晴天なり
# 今日のタスク
moge
piyo
# 作業ログ
# 呟き場
`,
			},
			want: `moge
piyo`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ExtractTask(tt.args.report); got != tt.want {
				t.Errorf("ExtractTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
