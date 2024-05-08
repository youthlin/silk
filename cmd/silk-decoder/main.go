package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/youthlin/go-lame"
	"github.com/youthlin/silk"
	"github.com/youthlin/silk/internal"
	"github.com/youthlin/t"
)

//go:embed *.po
var poFiles embed.FS
var (
	input      = flag.String("i", "", "")
	dir        = flag.String("d", "", "")
	sampleRate = flag.Int("sampleRate", 24000, "")
	mp3        = flag.Bool("mp3", true, "")
	verboe     = flag.Bool("verbose", false, "")
	output     = flag.String("o", "", "")
	lang       = flag.String("l", "", "")
	pattern    *regexp.Regexp
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *lang == "" {
		t.LoadFS(poFiles)
	} else {
		t.Load(*lang)
	}
	if *verboe {
		internal.Verbose = true
	}

	if *input == "" {
		printUsage()
		fmt.Println(t.T("[Error] input file are required.\n"))
		os.Exit(1)
	}

	if *dir == "" { // input file
		if err := decodeOneFile(*input, false); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return
	}

	// input dir
	exp, err := regexp.Compile(*dir)
	if err != nil {
		fmt.Println(t.T("[Error] input file name pattern %s are invalid: %+v", *dir, err))
		os.Exit(1)
	}
	pattern = exp

	filepath.WalkDir(*input, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !pattern.MatchString(d.Name()) {
			return nil // ignore
		}
		if err = decodeOneFile(path, true); err != nil {
			fmt.Println(err.Error()) // ignore error
		}
		return nil
	})

}

func decodeOneFile(path string, batch bool) error {
	in, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(t.T("failed to open input file %q: %w"), path, err)
	}
	defer in.Close()

	buf, err := silk.Decode(in, silk.WithSampleRate(*sampleRate))
	if err != nil {
		return fmt.Errorf(t.T("failed to decode input file %q: %w"), path, err)
	}

	var suffix = ".pcm"
	if *mp3 {
		suffix = ".mp3"
		var out bytes.Buffer
		wr, err := lame.NewWriter(&out)
		if err != nil {
			return fmt.Errorf(t.T("can not create mp3-encoder: %w"), err)
		}
		wr.InSampleRate = *sampleRate
		wr.InNumChannels = 1
		_, err = wr.Write(buf)
		if err != nil {
			return fmt.Errorf(t.T("failed to encode input file %q to mp3: %w"), path, err)
		}
		wr.Close()
		buf = out.Bytes()
	}

	var outputName = getOutputName(path, suffix, *output, batch)

	out, err := os.OpenFile(outputName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf(t.T("failed to open/create output file %q: %w"), outputName, err)

	}
	defer out.Close()

	_, err = out.Write(buf)
	if err != nil {
		return fmt.Errorf(t.T("failed to write output file %q: %w"), outputName, err)
	}
	return nil
}

func getOutputName(path, suffix, output string, batch bool) string {
	var outputName string
	if batch { // 批量
		if output != "" { // 指定输出后缀名
			if strings.HasPrefix(output, ".") {
				suffix = output
			} else {
				suffix = "." + output
			}
		}
		if i := strings.LastIndex(path, "."); i > 0 {
			outputName = path[:i] + suffix
		} else {
			outputName = path + suffix
		}
	} else { // 单个
		if output != "" { // 指定了输出名
			outputName = output
		} else { // 没指定输出文件名，默认为 pcm 或 mp3
			if i := strings.LastIndex(path, "."); i > 0 {
				outputName = path[:i] + suffix
			} else {
				outputName = path + suffix
			}
		}
	}
	return outputName
}

func printUsage() {
	var name = os.Args[0]
	fmt.Println()
	fmt.Println(t.T("Silk decoder, Go version, based on v1.0.9 of C version"))
	fmt.Println(t.T("Decode silk v3 file to pcm or mp3, by youthlin"))
	fmt.Println(t.T("GitHub: https://github.comyouthlin/silk"))
	fmt.Println()
	fmt.Println(t.T("Usage: %s -i <input file> [settings]", name))
	fmt.Println(t.T("  -i <input file>\tInput file or input folder(should with -d settings)"))
	fmt.Println(t.T("  [settings]"))
	fmt.Println(t.T("    -d <pattern>\tInput is a dir, and use the regexp <pattern> to test input file"))
	fmt.Println(t.T("    -sampleRate <hz>\tSample rate in Hz, default 24000"))
	fmt.Println(t.T("    -mp3[=false]\tOutput as mp3 file, default true, set false to output as pcm file"))
	fmt.Println(t.T("    -o <output file>\tOutput file name, or output file extension name when input is folder.\n\t\t\tIf not provide, output name is <input>.mp3 or <input>.pcm(when -mp3=false)"))
	fmt.Println(t.T("    -l <language>\tLanguage path(pointer to po file/dir)"))
	fmt.Println(t.T("    -verbose\t\tprint verbose log(default false)"))
	fmt.Println()
	fmt.Println(t.T("Example:"))
	fmt.Println(t.X("cmd-example", "%s -i a.amr\n\tdecode a.amr to a.mp3", name))
	fmt.Println(t.X("cmd-example", "%s -i amr.1\n\tdecode amr.1 to amr.mp3", name))
	fmt.Println(t.X("cmd-example", "%s -i file\n\tdecode file to file.mp3", name))
	fmt.Println(t.X("cmd-example", "%s -i a.amr -o b.mp3\n\tdecode a.amr to b.mp3", name))
	fmt.Println(t.X("cmd-example", "%s -i a.amr -mp3=false\n\tdecode a.amr to a.pcm", name))
	fmt.Println(t.X("cmd-example", "%s -i a.amr -mp3=false -o b.pcm\n\tdecode a.amr to b.pcm", name))
	fmt.Println(t.X("cmd-example", "%s -i voice -d \".*\\.amr\"\n\tdecode files in the folder to mp3\n\t  e.g.: if the voice folder has these files:\n\t\tvoice/a.amr\n\t\tvoice/other.txt\n\t\tvoice/sub/b.amr\n\t  result:\n\t\tvoice/a.mp3\n\t\tvoice/sub/b.mp3", name))
	fmt.Println()
}
