package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("❌ 参数错误!")
		fmt.Println("👉 用法: full2half <文件路径>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("❌ 无法访问文件: %v\n", err)
		os.Exit(1)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		os.Exit(1)
	}

	processedText := strings.Map(convertPunctuation, string(content))

	err = os.WriteFile(filePath, []byte(processedText), info.Mode())
	if err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 文件标点符号替换完成！")
}

func convertPunctuation(r rune) rune {
	switch r {
	case '。', '、', '《', '》':
		return r
	}

	switch r {
	case '“', '”':
		return '"'
	case '‘', '’':
		return '\''
	case '【':
		return '['
	case '】':
		return ']'
	case '　':
		return ' '
	case '，':
		return ','
	case '！':
		return '!'
	case '？':
		return '?'
	case '；':
		return ';'
	case '：':
		return ':'
	case '（':
		return '('
	case '）':
		return ')'
	case '—':
		return '-'
	case '～':
		return '~'
	case '｛':
		return '{'
	case '｝':
		return '}'
	case '／':
		return '/'
	case '＼':
		return '\\'
	case '｜':
		return '|'
	case '「', '」':
		return '\''
	case '『', '』':
		return '"'
	case '＋':
		return '+'
	case '－':
		return '-'
	case '＝':
		return '='
	case '＜':
		return '<'
	case '＞':
		return '>'
	case '＠':
		return '@'
	case '＃':
		return '#'
	case '％':
		return '%'
	case '＆':
		return '&'
	case '＊':
		return '*'
	}
	return r
}
