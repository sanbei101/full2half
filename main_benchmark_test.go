package main

import (
	"strings"
	"testing"
)

func BenchmarkNormalizePunctuation(b *testing.B) {
	benchmarks := []struct {
		name    string
		content []byte
	}{
		{
			name:    "NoChangeASCII",
			content: []byte(strings.Repeat("plain ascii text 12345, punctuation!? []{}() <> @#%&*\n", 4096)),
		},
		{
			name:    "MixedContent",
			content: []byte(strings.Repeat("你好，“世界”！【Go】（benchmark）　ASCII mixed，symbols：＠＃％＆＊ plus plain text\n", 4096)),
		},
		{
			name:    "DenseFullWidth",
			content: []byte(strings.Repeat("“”‘’【】　，！？；：（）—～｛｝／＼｜「」『』＋－＝＜＞＠＃％＆＊\n", 4096)),
		},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name+"/Legacy", func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(benchmark.content)))
			for b.Loop() {
				result := legacyNormalizePunctuation(benchmark.content)
				if len(result) == 0 && len(benchmark.content) > 0 {
					b.Fatal("unexpected empty result")
				}
			}
		})

		b.Run(benchmark.name+"/Optimized", func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(benchmark.content)))
			for b.Loop() {
				result, _ := normalizePunctuation(benchmark.content)
				if len(result) == 0 && len(benchmark.content) > 0 {
					b.Fatal("unexpected empty result")
				}
			}
		})
	}
}

func legacyNormalizePunctuation(content []byte) []byte {
	return []byte(strings.Map(legacyConvertPunctuation, string(content)))
}

func legacyConvertPunctuation(r rune) rune {
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
	default:
		return r
	}
}
