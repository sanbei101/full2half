package main

import (
	"bytes"
	"testing"
)

func TestIsBinaryExtension(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "main.go", want: false},
		{path: "README.TXT", want: false},
		{path: "image.PNG", want: true},
		{path: "archive.tar.gz", want: true},
		{path: "Makefile", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isBinaryExtension(tt.path); got != tt.want {
				t.Fatalf("isBinaryExtension(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsTextContent(t *testing.T) {
	if !isTextContent([]byte("代码，“文本”")) {
		t.Fatal("expected UTF-8 text to be accepted")
	}
	if isTextContent([]byte("text\x00with NUL")) {
		t.Fatal("expected content with NUL to be rejected")
	}
	if isTextContent([]byte{0xff, 0xfe, 0xfd}) {
		t.Fatal("expected invalid UTF-8 to be rejected")
	}
}

func TestNormalizePunctuation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        string
		wantChanged bool
	}{
		{
			name:        "no changes",
			input:       "plain ascii text 12345, punctuation!? []{}() <> @#%&*",
			want:        "plain ascii text 12345, punctuation!? []{}() <> @#%&*",
			wantChanged: false,
		},
		{
			name:        "mixed punctuation",
			input:       "你好，“世界”！【Go】（test）　symbols：＠＃％＆＊",
			want:        "你好,\"世界\"![Go](test) symbols:@#%&*",
			wantChanged: true,
		},
		{
			name:        "dense full width",
			input:       "“”‘’【】　，！？；：（）—～｛｝／＼｜「」『』＋－＝＜＞＠＃％＆＊",
			want:        `""''[] ,!?;:()-~{}/\|''""+-=<>@#%&*`,
			wantChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.input)
			got, changed := normalizePunctuation(input)

			if string(got) != tt.want {
				t.Fatalf("unexpected output: got %q want %q", string(got), tt.want)
			}

			if changed != tt.wantChanged {
				t.Fatalf("unexpected changed flag: got %v want %v", changed, tt.wantChanged)
			}

			legacy := legacyNormalizePunctuation(input)
			if !bytes.Equal(got, legacy) {
				t.Fatalf("optimized output differs from legacy output: got %q want %q", string(got), string(legacy))
			}
		})
	}
}
