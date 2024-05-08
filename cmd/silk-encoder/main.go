package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/youthlin/silk"
	"github.com/youthlin/silk/internal"
	"github.com/youthlin/t"
)

//go:embed *.po
var poFiles embed.FS
var args = &appArgs{}

func main() {
	parseArgs()
	i18n()
	if args.input == "" || args.output == "" {
		printUsage()
		fmt.Println(t.T("[Error] both input file and output file are required.\n"))
		os.Exit(1)
	}

	input, err := os.Open(args.input)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to open input file %q: %+v", args.input, err))
		os.Exit(1)
	}

	buf, err := silk.Encode(
		input,
		silk.SampleRate(args.FsAPI),
		silk.MaxInternal(args.FsMaxInternal),
		silk.PacketSizeMs(args.PacketLength),
		silk.PackageLossPct(args.LossPct),
		silk.InbandFEC(args.InbandFEC),
		silk.UseDTX(args.DTX),
		silk.Complexity(args.Complexity),
		silk.BitRate(args.Rate),
		silk.Stx(args.STX),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to encode input file %q: %+v", args.input, err))
		os.Exit(1)
	}

	output, err := os.OpenFile(args.output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to open output file %q: %+v", args.output, err))
		os.Exit(1)
	}
	_, err = output.Write(buf)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, t.T("failed to write output file %q: %+v", args.output, err))
		os.Exit(1)
	}
}

type appArgs struct {
	lang          string
	input         string
	output        string
	FsAPI         int
	FsMaxInternal int
	PacketLength  int
	Rate          int
	LossPct       int
	Complexity    int
	InbandFEC     bool
	DTX           bool
	STX           bool
	Verbose       bool
}

func i18n() {
	if args.lang == "" {
		t.LoadFS(poFiles)
	} else {
		t.Load(args.lang)
	}
}

func parseArgs() {
	flag.StringVar(&args.lang, "l", "", "")
	flag.StringVar(&args.input, "i", "", "")
	flag.StringVar(&args.output, "o", "", "")
	flag.IntVar(&args.FsAPI, "Fs_API", 24000, "")
	flag.IntVar(&args.FsMaxInternal, "Fs_maxInternal", 24000, "")
	flag.IntVar(&args.PacketLength, "packetlength", 20, "")
	flag.IntVar(&args.Rate, "rate", 25000, "")
	flag.IntVar(&args.LossPct, "loss", 0, "")
	flag.BoolVar(&args.InbandFEC, "inbandFEC", false, "")
	flag.IntVar(&args.Complexity, "complexity", 2, "")
	flag.BoolVar(&args.DTX, "DTX", false, "")
	flag.BoolVar(&args.STX, "STX", true, "")
	flag.BoolVar(&args.Verbose, "verbose", false, "")
	flag.Usage = printUsage
	flag.Parse()
	internal.Verbose = args.Verbose
}

func printUsage() {
	fmt.Println()
	fmt.Println(t.T("Silk encoder, Go version, based on v1.0.9 of C version"))
	fmt.Println(t.T("Encode pcm file to silk v3 type, by youthlin"))
	fmt.Println(t.T("GitHub: https://github.comyouthlin/silk"))
	fmt.Println()
	fmt.Println(t.T("Usage: %s [settings]", os.Args[0]))
	fmt.Println(t.T("  [settings]"))
	fmt.Println(t.T("    -l <path to po file>\tlanguage path(pointer to po file/dir)"))
	fmt.Println(t.T("    -i <input file>\t\tSpeech input to encoder"))
	fmt.Println(t.T("    -o <output file>\t\tBitstream output from encoder"))
	fmt.Println(t.T("    -Fs_API <Hz>\t\tAPI sampling rate in Hz, default: 24000"))
	fmt.Println(t.T("    -Fs_maxInternal <Hz>\tMaximum internal sampling rate in Hz, default: 24000"))
	fmt.Println(t.T("    -packetlength <ms>\t\tPacket interval in ms, default: 20"))
	fmt.Println(t.T("    -rate <bps>\t\t\tTarget bitrate; default: 25000"))
	fmt.Println(t.T("    -loss <perc>\t\tUplink loss estimate, in percent (0-100); default: 0"))
	fmt.Println(t.T("    -inbandFEC[=false]\t\tEnable inband FEC usage, default: false"))
	fmt.Println(t.T("    -complexity <comp>\t\tSet complexity, 0: low, 1: medium, 2: high; default: 2"))
	fmt.Println(t.T("    -DTX\t\t\tEnable DTX; default: false"))
	fmt.Println(t.T("    -stx[=false]\t\tAdd STX flag before file header and remove footer block, default true"))
	fmt.Println(t.T("    -verbose\t\t\tprint verbose log, default false"))
	fmt.Println()
}
