package internal

/*
#cgo CFLAGS: -Wno-shift-negative-value -Wno-constant-conversion
#include "SKP_Silk_SDK_API.h"
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

const (
	MAX_BYTES_PER_FRAME = 250 // Equals peak bitrate of 100 kbps
	MAX_INPUT_FRAMES    = 5
	FRAME_LENGTH_MS     = 20
	MAX_API_FS_KHZ      = 48
)

type EncodeCfg struct {
	SampleRate            int
	MaxInternalSampleRate int
	PacketSizeMs          int
	PacketLossPct         int
	UseInBandFEC          bool
	UseDTX                bool
	ComplexityMode        int
	BitRate               int
	Stx                   bool
}

type EncodeOpt func(*EncodeCfg)

func Encode(src io.Reader, opts ...EncodeOpt) ([]byte, error) {
	var cfg = buildCfg(opts...)
	if cfg.SampleRate > MAX_API_FS_KHZ*1000 || cfg.SampleRate < 0 {
		return nil, fmt.Errorf("error: API sampling rate = %d out of range, valid range 8000 - 48000", cfg.SampleRate)
	}
	/* Print options */
	log("encode options: %#v", cfg)

	var out = &bytes.Buffer{}

	/* Add Silk header to stream */
	if cfg.Stx {
		out.Write([]byte{STX})
	}
	out.Write([]byte(Header))

	/* Create Encoder */
	var encSizeBytes = getEncoderSize()
	var psEnc, free = malloc(encSizeBytes)
	defer free()

	/* Reset Encoder */
	initEncode(psEnc)

	if err := doEncode(src, out, cfg, psEnc); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func buildCfg(opts ...EncodeOpt) *EncodeCfg {
	var cfg = &EncodeCfg{
		SampleRate:            defaultSampleRate,
		MaxInternalSampleRate: 0,
		BitRate:               25000,
		PacketSizeMs:          20,
		ComplexityMode:        2,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.MaxInternalSampleRate == 0 {
		cfg.MaxInternalSampleRate = defaultSampleRate
		if cfg.SampleRate < cfg.MaxInternalSampleRate {
			cfg.MaxInternalSampleRate = cfg.SampleRate
		}
	}
	return cfg
}

func getEncoderSize() (encSizeBytes int32) {
	C.SKP_Silk_SDK_Get_Encoder_Size((*C.SKP_int32)(unsafe.Pointer(&encSizeBytes)))
	return
}

func initEncode(psEnc unsafe.Pointer) {
	var encStatus C.SKP_SILK_SDK_EncControlStruct
	/*************************/
	/* Init or reset encoder */
	/*************************/
	// SKP_int SKP_Silk_SDK_InitEncoder(
	//   void                          *encState,      /* I/O: State                                           */
	//   SKP_SILK_SDK_EncControlStruct *encStatus      /* O:   Encoder Status                                  */
	// );
	C.SKP_Silk_SDK_InitEncoder(psEnc, (*C.SKP_SILK_SDK_EncControlStruct)(unsafe.Pointer(&encStatus)))
}

func buildEncControl(cfg *EncodeCfg) *C.SKP_SILK_SDK_EncControlStruct {
	var encControl C.SKP_SILK_SDK_EncControlStruct
	encControl.API_sampleRate = C.SKP_int32(cfg.SampleRate)
	encControl.maxInternalSampleRate = C.SKP_int32(cfg.MaxInternalSampleRate)
	encControl.packetSize = C.SKP_int(cfg.PacketSizeMs * cfg.SampleRate / 1000)
	encControl.packetLossPercentage = C.SKP_int(cfg.PacketLossPct)
	if cfg.UseInBandFEC {
		encControl.useInBandFEC = C.SKP_int(1)
	} else {
		encControl.useInBandFEC = C.SKP_int(0)
	}
	if cfg.UseDTX {
		encControl.useDTX = C.SKP_int(1)
	} else {
		encControl.useDTX = C.SKP_int(0)
	}
	encControl.complexity = C.SKP_int(cfg.ComplexityMode)
	if cfg.BitRate > 0 {
		encControl.bitRate = C.SKP_int32(cfg.BitRate)
	} else {
		encControl.bitRate = C.SKP_int32(0)
	}
	return &encControl
}

func doEncode(reader io.Reader, out io.Writer, cfg *EncodeCfg, psEnc unsafe.Pointer) error {
	const frameSizeReadFromFile_ms = 20
	var (
		/* Set Encoder parameters */
		encControl = buildEncControl(cfg)
		frameSize  = frameSizeReadFromFile_ms * cfg.SampleRate / 1000
		// C 源码中是按 sizeof( SKP_int16 ) 读取的
		// 每次读取 frameSize 个 SKP_int16 大小
		// 这里我们的 in 是 []byte 类型，所以需要 *2
		in         = make([]byte, frameSize*2)
		nBytes     = int16(MAX_BYTES_PER_FRAME * MAX_INPUT_FRAMES)
		payload    = make([]byte, nBytes)
		blockIndex int
	)
	log("encode frameSize=%d", frameSize)
	for {
		blockIndex++

		// 读取一段数据
		n, err := io.ReadFull(reader, in)
		log("block=%d, read n=%d, err=%+v", blockIndex, n, err)
		if errors.Is(err, io.EOF) {
			err = nil
			log("block=%d, EOF when read data", blockIndex)
			break
		}
		if n < frameSize {
			break
		}

		// 编码
		nBytes = MAX_BYTES_PER_FRAME * MAX_INPUT_FRAMES
		/**************************/
		/* Encode frame with Silk */
		/**************************/
		// SKP_int SKP_Silk_SDK_Encode(
		//     void                                *encState,      /* I/O: State                                           */
		//     const SKP_SILK_SDK_EncControlStruct *encControl,    /* I:   Control status                                  */
		//     const SKP_int16                     *samplesIn,     /* I:   Speech sample input vector                      */
		//     SKP_int                             nSamplesIn,     /* I:   Number of samples in input vector               */
		//     SKP_uint8                           *outData,       /* O:   Encoded output vector                           */
		//     SKP_int16                           *nBytesOut      /* I/O: Number of bytes in outData (input: Max bytes)   */
		// );
		ret := C.SKP_Silk_SDK_Encode(
			psEnc,
			encControl,
			(*C.SKP_int16)(unsafe.Pointer(&in[0])),
			C.SKP_int(n/2), // in 是 []byte 类型，看做 SKP_int16 数组的话，长度需要 / 2
			(*C.SKP_uint8)(unsafe.Pointer(&payload[0])), // 接收 encode 后的数据
			(*C.SKP_int16)(unsafe.Pointer(&nBytes)),     // 接收 encode 后的长度
		)
		log("encode ret code=%d, encode payload size=%d, data=%x", ret, nBytes, payload[:nBytes])
		if ret != 0 {
			warn("encode failed, ret=%d", ret)
			continue // or break?
		}

		// 写入编码后的长度、内容
		err = binary.Write(out, binary.LittleEndian, nBytes)
		if err != nil {
			warn("failed to write block size, err=%+v", err)
			return fmt.Errorf("failed to write block size: %w", err)
		}
		_, err = out.Write(payload[:nBytes])
		if err != nil {
			warn("failed to write block data, err=%+v", err)
			return fmt.Errorf("failed to write block data: %w", err)
		}
	}

	if !cfg.Stx {
		binary.Write(out, binary.LittleEndian, int16(-1)) // footer block
	}

	return nil
}
