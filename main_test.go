package main

import (
	"bytes"
	"testing"
)

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
