package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("❌ 参数错误!")
		fmt.Println("👉 用法: full2half <文件路径>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		os.Exit(1)
	}

	processedContent, changed := normalizePunctuation(content)
	if !changed {
		fmt.Println("✅ 文件无需替换")
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("❌ 无法访问文件: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(filePath, processedContent, info.Mode())
	if err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 文件标点符号替换完成！")
}

func normalizePunctuation(content []byte) ([]byte, bool) {
	for index := 0; index < len(content); {
		if content[index] < utf8.RuneSelf {
			index++
			continue
		}

		r, size := utf8.DecodeRune(content[index:])
		converted, ok := convertPunctuation(r)
		if !ok {
			index += size
			continue
		}

		processed := make([]byte, 0, len(content))
		processed = append(processed, content[:index]...)
		processed = appendRune(processed, converted)
		index += size

		for index < len(content) {
			if content[index] < utf8.RuneSelf {
				processed = append(processed, content[index])
				index++
				continue
			}

			r, size = utf8.DecodeRune(content[index:])
			if converted, ok = convertPunctuation(r); ok {
				processed = appendRune(processed, converted)
			} else {
				processed = append(processed, content[index:index+size]...)
			}
			index += size
		}

		return processed, true
	}

	return content, false
}

func appendRune(dst []byte, r rune) []byte {
	var encoded [utf8.UTFMax]byte
	n := utf8.EncodeRune(encoded[:], r)
	return append(dst, encoded[:n]...)
}

func convertPunctuation(r rune) (rune, bool) {

	switch r {
	case '“', '”':
		return '"', true
	case '‘', '’':
		return '\'', true
	case '【':
		return '[', true
	case '】':
		return ']', true
	case '　':
		return ' ', true
	case '，':
		return ',', true
	case '！':
		return '!', true
	case '？':
		return '?', true
	case '；':
		return ';', true
	case '：':
		return ':', true
	case '（':
		return '(', true
	case '）':
		return ')', true
	case '—':
		return '-', true
	case '～':
		return '~', true
	case '｛':
		return '{', true
	case '｝':
		return '}', true
	case '／':
		return '/', true
	case '＼':
		return '\\', true
	case '｜':
		return '|', true
	case '「', '」':
		return '\'', true
	case '『', '』':
		return '"', true
	case '＋':
		return '+', true
	case '－':
		return '-', true
	case '＝':
		return '=', true
	case '＜':
		return '<', true
	case '＞':
		return '>', true
	case '＠':
		return '@', true
	case '＃':
		return '#', true
	case '％':
		return '%', true
	case '＆':
		return '&', true
	case '＊':
		return '*', true
	}
	return r, false
}
