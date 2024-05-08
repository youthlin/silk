package internal

/*
#cgo CFLAGS: -Wno-shift-negative-value -Wno-constant-conversion
#include "SKP_Silk_SDK_API.h"
#include <stdlib.h>
*/
import "C"

// 导入 C 包说明： 必须单独写一行，标记需要 cgo
// 在这单独的一行上面的注释内容会当做 C 源码编译
// 注释和导入 C 包必须紧挨这 空间不能空行

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"unsafe"
)

const (
	STX       byte = 2           // 文件开头如果是 0x02 解码时需要丢弃
	Header         = "#!SILK_V3" // 文件头
	HeaderLen      = len(Header) // 文件头长度 = 9
	// 默认值
	defaultSampleRate = 24000
)

type DecodeCfg struct {
	SampleRate int
}

type DecodeOpt func(*DecodeCfg)

func Decode(src io.Reader, opts ...DecodeOpt) ([]byte, error) {
	/* set option */
	var cfg = &DecodeCfg{SampleRate: defaultSampleRate}
	for _, opt := range opts {
		opt(cfg)
	}
	log("decode option: %#v", cfg)

	var reader = bufio.NewReader(src)

	/* Check Silk header */
	if err := checkHeader(reader); err != nil {
		return nil, err
	}

	/* Create decoder */
	var decSize = getDecoderSize()
	var psDec, free = malloc(decSize)
	defer free()

	/* Reset decoder */
	initDecoder(psDec)

	out := &bytes.Buffer{}
	/* decode */
	if err := doDecode(reader, psDec, cfg.SampleRate, out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func checkHeader(reader *bufio.Reader) error {
	first, err := reader.Peek(1)
	if err != nil {
		warn("io error / failed to peek first byte: %+v", err)
		return fmt.Errorf("failed to peek first byte: %w", err)
	}
	// 如果第一位是 0x02 需要丢弃
	// 安卓移植版说明:
	// https://wufengxue.github.io/2019/04/17/wechat-voice-codec-amr.html
	// kn007 兼容版本:
	// https://github.com/kn007/silk-v3-decoder/blob/master/silk/test/Decoder.c#L187
	// 原始开源版本:(不识别 0x02 开头的文件)
	// https://github.com/gaozehua/SILKCodec/blob/master/SILK_SDK_SRC_ARM/test/Decoder.c#L182
	if first[0] == STX {
		log("first byte is STX(%x), read it", STX)
		stx, err := reader.ReadByte()
		if err != nil {
			warn("read first byte error: %+v", err)
			return fmt.Errorf("failed to read first byte: %w", err)
		}
		if stx != STX {
			warn("read first byte not STX: %x", stx)
			return fmt.Errorf("invalid first byte: %d, expected=%d", stx, STX)
		}
	}
	// 文件头
	var header = make([]byte, HeaderLen)
	n, err := io.ReadFull(reader, header)
	if err != nil {
		warn("failed to read file header: %+v", err)
		return fmt.Errorf("failed to read file header: %w", err)
	}
	if n != HeaderLen {
		warn("invalid file header, read %d bytes, expected %d", n, HeaderLen)
		return fmt.Errorf("invalid file header, length=%d, expected=%d", n, HeaderLen)
	}
	if string(header) != Header {
		warn("invalid file header %q expected %q", header, HeaderLen)
		return fmt.Errorf("invalid file header, got=%q, expected=%q", header, Header)
	}
	return nil
}

// getDecoderSize wrap the C function SKP_Silk_SDK_Get_Decoder_Size,
// return the size of SKP_Silk_decoder_state, which is implemented in SKP_Silk_dec_API.c
// 包装 C 函数 SKP_Silk_SDK_Get_Decoder_Size, 返回 sizeof(SKP_Silk_decoder_state) 大小，
// C 语言实现文件 SKP_Silk_dec_API.c
func getDecoderSize() (decSize int32) {
	// 函数声明
	/***********************************************/
	/* Get size in bytes of the Silk decoder state */
	/***********************************************/
	// SKP_int SKP_Silk_SDK_Get_Decoder_Size(
	//     SKP_int32 *decSizeBytes   /* O: Number of bytes in SILK decoder state */
	// );
	//
	// 入参类型是 SKP_int32 *,
	// SKP_Silk_SDK_API.h include 了 SKP_Silk_typedef.h:
	//     #define SKP_int32       int
	// 可以看到 SKP_int32 就是 int 类型，对应 Go 语言中 int32， 所以 decSize 是 int32 类型
	// 函数声明中注释了这个是一个 Output 参数，所以是一个 SKP_int32 指针，才能用来接受结果
	// 因为是指针，所以要用 unsafe.Pointer 取地址, 再通过类型强转指明是 SKP_int32 指针类型
	C.SKP_Silk_SDK_Get_Decoder_Size((*C.SKP_int32)(unsafe.Pointer(&decSize)))
	return
}

// malloc wrap the C funcion malloc, return the pointer and free function.
// 包装了 C 中的 malloc 函数，返回指向申请的内存的指针，和一个释放内存的 free 函数
func malloc(size int32) (p unsafe.Pointer, free func()) {
	p = C.malloc(C.ulong(size))
	return p, func() {
		C.free(p)
	}
}

// initDecoder wrap the C function SKP_Silk_SDK_InitDecoder
// 初始化 Decoder
func initDecoder(psDec unsafe.Pointer) {
	/*************************/
	/* Init or Reset decoder */
	/*************************/
	// SKP_int SKP_Silk_SDK_InitDecoder(
	//     void *decState    /* I/O: State */
	// );
	C.SKP_Silk_SDK_InitDecoder(psDec)
}

func doDecode(reader io.Reader, psDec unsafe.Pointer, sampleRate int, out io.Writer) (err error) {
	var (
		blockIndex int // for debug log
		decControl C.SKP_SILK_SDK_DecControlStruct
		// in 对应 C 源码中 payload(SKP_uint8 数组), buf 对应 out(SKP_int16 数组)
		in = make([]byte, 1024) // Decoder.c 中 MAX_BYTES_PER_FRAME 和 Encoder.c 不一样哦
		// 20ms FRAME_LENGTH_MS=20 MAX_API_FS_KHZ=48
		frameSize = (FRAME_LENGTH_MS * MAX_API_FS_KHZ) << 1
		// frameSize 个 SKP_int16，这里是 []byte 所以 *2
		buf = make([]byte, frameSize*2) // 相当于 [frameSize]int16 大小
	)
	decControl.API_sampleRate = C.SKP_int32(sampleRate)
	decControl.framesPerPacket = C.SKP_int(1)

	// https://github.com/kn007/silk-v3-decoder/blob/master/silk/test/Decoder.c
	// https://github.com/gaozehua/SILKCodec/blob/master/SILK_SDK_SRC_ARM/test/Decoder.c
	// C 版本的 decoder 模拟了数据丢失 然后一顿操作靠其他帧修复 这里省略, 直接按 frame 解码
	// 这个 go 库也是这么做的
	// https://github.com/wdvxdr1123/go-silk/blob/main/silk.go
	for {
		blockIndex++
		// 参考格式说明
		// https://wufengxue.github.io/2019/04/17/wechat-voice-codec-amr.html
		// 文件头之后，就是每个 block, 先是 16 字节的 block 大小 n，然后是 n 个字节内容
		// 最后是 footer 部分，内容是 0xffff, 也可以看做是一个 block(大小是 -1，没有内容)

		var nByte int16 // 先读取 block 大小, 占两个字节，用 int16 接收
		err = binary.Read(reader, binary.LittleEndian, &nByte)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				log("packet=%d, EOF when read block size", blockIndex)
				break
			}
			log("packet=%d, read block size err=%+v", blockIndex, err)
			return fmt.Errorf("failed to read block size: %w", err)
		}
		log("packet=%d, block size=%d", blockIndex, nByte)
		if nByte < 0 {
			break // 是 footer 部分, 没有 block 内容
		}
		if int(nByte) > len(in) { // 兜底 or 报错?
			in = make([]byte, nByte)
		}

		// 再读取 block 内容，长度就是 nByte
		n, err := io.ReadFull(reader, in[:nByte])
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				log("packet=%d, EOF when read block data", blockIndex)
				break
			}
			warn("packet=%d, read block data err=%+v", blockIndex, err)
			return fmt.Errorf("failed to read block: %w", err)
		}
		if n != int(nByte) {
			log("packet=%d, read block data invalid, read %d bytes, expected %d", blockIndex, n, nByte)
			return fmt.Errorf("invalid block")
		}

		// 解码
		C.SKP_Silk_SDK_Decode(
			psDec,                                   // State
			&decControl,                             // Control Structure
			0,                                       // 0: no loss, 1 loss
			(*C.SKP_uint8)(unsafe.Pointer(&in[0])),  // Encoded input vector
			C.SKP_int(n),                            // Number of input bytes
			(*C.SKP_int16)(unsafe.Pointer(&buf[0])), // Decoded output speech vector
			(*C.SKP_int16)(unsafe.Pointer(&nByte)),  // Number of samples (vector/decoded)
		)
		// buf 是 []byte 类型，但是实际上 SKP_Silk_SDK_Decode 输出的是 SKP_int16 数组
		// 也就是说写入了 nByte 个 SKP_int16, 所以按 []byte 计算需要 *2
		_, err = out.Write(buf[:nByte*2])
		if err != nil {
			return fmt.Errorf("failed to write decode data: %w", err)
		}
	}
	return nil
}

var Verbose = false

func log(msg string, args ...any) {
	if Verbose {
		fmt.Printf("[INFO][silk] "+msg+"\n", args...)
	}
}

func warn(msg string, args ...any) {
	if Verbose {
		fmt.Fprintf(os.Stderr, "[Warn][silk] "+msg+"\n", args...)
	}
}
