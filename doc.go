package silk

import (
	"io"

	"github.com/youthlin/silk/internal"
)

// -------------------- Decode --------------------

// Decode decodes silk encode src to pcm.
// 解码 silk 格式为 pcm 格式.
func Decode(src io.Reader, opts ...internal.DecodeOpt) ([]byte, error) {
	return internal.Decode(src, opts...)
}

// WithSampleRate set decode option, sample rate, default 24000
// 设置 sample rate 参数，默认值 24000
func WithSampleRate(sampleRate int) internal.DecodeOpt {
	return func(dc *internal.DecodeCfg) { dc.SampleRate = sampleRate }
}

// -------------------- Encode --------------------

// Encode encode pcm file to silk v3 type.
// 将 pcm 格式编码为 silk v3 格式.
func Encode(src io.Reader, opts ...internal.EncodeOpt) ([]byte, error) {
	return internal.Encode(src, opts...)
}

// SampleRate set sample rate, default 24000.
// 设置 sample rate 参数，默认值 24000
func SampleRate(sampleRate int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.SampleRate = sampleRate }
}

// MaxInternal maximum internal sampling rate in Hz, default: 24000
// 最大采样率, 默认值 24000Hz
func MaxInternal(sampleRate int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.MaxInternalSampleRate = sampleRate }
}

// PacketSizeMs set packet interval in ms, default: 20
// 设置每个数据包的长度，默认值 20 毫秒
func PacketSizeMs(packageLength int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.PacketSizeMs = packageLength }
}

// PackageLossPct stt the uplink loss estimate, in percent (0-100); default: 0
// 设置数据丢失率，默认值 0
func PackageLossPct(loss int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.PacketLossPct = loss }
}

// InbandFEC Enable inband FEC usage; default: false
// 是否开启音频带内 FEC(前向纠错)
func InbandFEC(enable bool) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.UseInBandFEC = enable }
}

// UseDTX Enable DTX; default: false
func UseDTX(enable bool) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.UseDTX = enable }
}

// Complexity set complexity, 0: low, 1: medium, 2: high; default: 2
// 设置复杂模式，0=低，1=中，2=高
func Complexity(mode int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.ComplexityMode = mode }
}

// BitRate set target bitrate; default: 25000
// 设置比特率，默认值 25000
func BitRate(bitRate int) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.BitRate = bitRate }
}

// Stx set stx flag to compatible with QQ/Wechat: add stx flag before file header and remove footer block.
// 设置 stx 标记，在文件头之前添加 stx 标记(0x02)，并移除 footer 数据块
func Stx(enable bool) internal.EncodeOpt {
	return func(ec *internal.EncodeCfg) { ec.Stx = enable }
}
