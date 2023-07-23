package main

import "testing"

func Test_getOutputName(t *testing.T) {
	type args struct {
		path   string
		suffix string
		output string
		batch  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"single-no-output-default-mp3", args{"test.amr", ".mp3", "", false}, "test.mp3"},
		{"single-no-output-default-pcm", args{"test.amr", ".pcm", "", false}, "test.pcm"},
		{"single-no-output-no-ext", args{"test", ".pcm", "", false}, "test.pcm"},
		{"single-with-output", args{"test.amr", ".pcm", "some.bit", false}, "some.bit"},
		{"batch", args{"voice/test.amr", ".mp3", "", true}, "voice/test.mp3"},
		{"batch-pcm", args{"voice/test.amr", ".pcm", "", true}, "voice/test.pcm"},
		{"batch-output-ext", args{"voice/test.amr", "", "mp3", true}, "voice/test.mp3"},
		{"batch-output-dotext", args{"voice/test.amr", "", ".pcm", true}, "voice/test.pcm"},
		{"batch-output-no-ext", args{"voice/test", "", ".pcm", true}, "voice/test.pcm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOutputName(tt.args.path, tt.args.suffix, tt.args.output, tt.args.batch); got != tt.want {
				t.Errorf("getOutputName() = %v, want %v", got, tt.want)
			}
		})
	}
}
