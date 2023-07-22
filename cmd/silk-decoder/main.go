package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/youthlin/silk"
	"github.com/youthlin/t"
)

//go:embed *.po
var poFiles embed.FS
var (
	lang       = flag.String("l", "", "")
	input      = flag.String("i", "", "")
	output     = flag.String("o", "", "")
	sampleRate = flag.Int("sampleRate", 24000, "")
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if *lang == "" {
		t.LoadFS(poFiles)
	} else {
		t.Load(*lang)
	}

	if *input == "" || *output == "" {
		printUsage()
		fmt.Println(t.T("[Error] both input file and output file are required.\n"))
		os.Exit(1)
	}

	in, err := os.Open(*input)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to open input file %q: %+v", *input, err))
		os.Exit(1)
	}
	buf, err := silk.Decode(in, silk.WithSampleRate(*sampleRate))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to decode input file %q: %+v", *input, err))
		os.Exit(1)
	}
	out, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to open output file %q: %+v", *output, err))
		os.Exit(1)
	}
	_, err = out.Write(buf)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to write output file %q: %+v", *output, err))
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println()
	fmt.Println(t.T("Silk decoder, Go version, based on v1.0.9 of C version"))
	fmt.Println(t.T("Decode silk v3 file to pcm, by youthlin"))
	fmt.Println(t.T("GitHub: https://github.comyouthlin/silk"))
	fmt.Println()
	fmt.Println(t.T("Usage: %s -i <input file> -o <output file> [-sampleRate <hz>] [-l <language>]", os.Args[0]))
	fmt.Println(t.T("  -i <input file>\tBitstream input to decoder"))
	fmt.Println(t.T("  -o <output file>\tSpeech output from decoder"))
	fmt.Println(t.T("  -sampleRate <hz>\tSample rate in Hz, default 24000"))
	fmt.Println(t.T("  -l <language>\t\tLanguage path(pointer to po file/dir)"))
	fmt.Println()
}
