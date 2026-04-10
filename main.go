package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func main() {
	args := os.Args[1:]

	if len(args) == 2 && args[0] == "-r" {
		dirPath := args[1]
		info, err := os.Stat(dirPath)
		if err != nil {
			fmt.Printf("❌ 无法访问目录: %v\n", err)
			os.Exit(1)
		}
		if !info.IsDir() {
			fmt.Printf("❌ 路径不是目录: %s\n", dirPath)
			os.Exit(1)
		}
		err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("❌ 访问路径失败 %s: %v\n", path, err)
				return nil
			}
			if d.IsDir() {
				return nil
			}
			processFile(path)
			return nil
		})
		if err != nil {
			fmt.Printf("❌ 遍历目录失败: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if len(args) == 1 {
		processFile(args[0])
		return
	}

	fmt.Println("❌ 参数错误!")
	fmt.Println("👉 用法: full2half <文件路径>")
	fmt.Println("👉 用法: full2half -r <目录路径>")
	os.Exit(1)
}

func processFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		return
	}

	processedContent, changed := normalizePunctuation(content)
	if !changed {
		fmt.Printf("✅ 文件无需替换: %s\n", filePath)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("❌ 无法访问文件: %v\n", err)
		return
	}

	err = os.WriteFile(filePath, processedContent, info.Mode())
	if err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 文件标点符号替换完成: %s\n", filePath)
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
