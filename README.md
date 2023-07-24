# silk
Silk v3 Encoder / Decoder implement in Golang.

Go 语言版本的 Silk v3 编码/解码器。

可用于解码国内通信软件语音文件。

## Install 获取
```
go get -u github.com/youthlin/silk
```

## API 接口
```
func Decode(src io.Reader, opts ...internal.DecodeOpt) ([]byte, error)
func Encode(src io.Reader, opts ...internal.EncodeOpt) ([]byte, error)
```
see [API doc](https://pkg.go.dev/github.com/youthlin/silk)

## Comandline tool 命令行
### [silk-decoder](./cmd/silk-decoder/) 解码器

> You may need to run `sudo apt-get install libmp3lame-dev` to get lame lib on Linux

```
go install github.com/youthlin/silk/cmd/silk-decoder@latest
# Execute to see usage
silk-decoder
```

```
Silk decoder, Go version, based on v1.0.9 of C version
Decode silk v3 file to pcm or mp3, by youthlin
GitHub: https://github.comyouthlin/silk

Usage: silk-decoder -i <input file> [settings]
  -i <input file>       Input file or input folder(should with -d settings)
  [settings]
    -d <pattern>        Input is a dir, and use the regexp <pattern> to test input file
    -sampleRate <hz>    Sample rate in Hz, default 24000
    -mp3[=false]        Output as mp3 file, default true, set false to output as pcm file
    -o <output file>    Output file name, or output file extension name when input is folder.
                        If not provide, output name is <input>.mp3 or <input>.pcm(when -mp3=false)
    -l <language>       Language path(pointer to po file/dir)

Example:
silk-decoder -i a.amr
        decode a.amr to a.mp3
silk-decoder -i amr.1
        decode amr.1 to amr.mp3
silk-decoder -i file
        decode file to file.mp3
silk-decoder -i a.amr -o b.mp3
        decode a.amr to b.mp3
silk-decoder -i a.amr -mp3=false
        decode a.amr to a.pcm
silk-decoder -i a.amr -mp3=false -o b.pcm
        decode a.amr to b.pcm
silk-decoder -i voice -d ".*\.amr"
        decode files in the folder to mp3
          e.g.: if the voice folder has these files:
                voice/a.amr
                voice/other.txt
                voice/sub/b.amr
          result:
                voice/a.mp3
                voice/sub/b.mp3

Silk 解码器，Go 语言版本，基于 v1.0.9 的 C 语言版本
将 silk v3 格式的文件解码为 pcm 或 mp3, 作者：youthlin
GitHub: https://github.comyouthlin/silk

用法：silk-decoder -i <输入文件> [选项]
  -i <输入文件>         输入文件或输入文件夹(需要和 -d 连用)
  [选项]
    -d <正则表达式>             指明 -i 的参数是文件夹，对输入文件夹(及子文件夹中)中，文件名符合正规表达式的文件进行解码
    -sampleRate <采样率>        单位为赫兹，默认值为 24000
    -mp3[=false]        输出为 mp3 格式，默认 true, 设置为 flase 以输出 pcm 格式
    -o <输出文件>       指定输出文件名，或指定输出文件后缀名（当使用-d 时）。
                        如果为空输出文件会根据自动推断为 mp3 或 pcm
    -l <语言>           指定语言路径(po 文件或文件夹)

示例：
silk-decoder -i a.amr
        将 a.amr 解码为 a.mp3
silk-decoder -i amr.1
        将 amr.1 解码为 amr.mp3
silk-decoder -i file
        将 file 解码为 file.mp3
silk-decoder -i a.amr -o b.mp3
        将 a.amr 解码为 b.mp3
silk-decoder -i a.amr -mp3=false
        将 a.amr 解码为 a.pcm
silk-decoder -i a.amr -mp3=false -o b.pcm
        将 a.amr 解码为 b.pcm
silk-decoder -i voice -d ".*\.amr"
          例如：voice 文件夹下有如下文件：
                voice/a.amr
                voice/other.txt
                voice/sub/b.amr
          转换结果：
                voice/a.mp3
                voice/sub/b.mp3
```

### [silk-encoder](./cmd/silk-encoder/) 编码器
```
go install github.com/youthlin/silk/cmd/silk-encoder@latest
# Execute to see usage
silk-encoder
```

```
Silk encoder, Go version, based on v1.0.9 of C version
Encode pcm file to silk v3 type, by youthlin
GitHub: https://github.comyouthlin/silk

Usage: silk-encoder [settings]
  [settings]
    -l <path to po file>        language path(pointer to po file/dir)
    -i <input file>             Speech input to encoder
    -o <output file>            Bitstream output from encoder
    -Fs_API <Hz>                API sampling rate in Hz, default: 24000
    -Fs_maxInternal <Hz>        Maximum internal sampling rate in Hz, default: 24000
    -packetlength <ms>          Packet interval in ms, default: 20
    -rate <bps>                 Target bitrate; default: 25000
    -loss <perc>                Uplink loss estimate, in percent (0-100); default: 0
    -inbandFEC[=false]          Enable inband FEC usage, default: false
    -complexity <comp>          Set complexity, 0: low, 1: medium, 2: high; default: 2
    -DTX[=false]                Enable DTX; default: false
    -stx[=false]                Add STX flag before file header and remove footer block, default true

Silk 编码器，Go 语言版本，基于 v1.0.9 的 C 语言版本
将 pcm 文件编码为 silk v3 类型，作者： youthlin
GitHub: https://github.comyouthlin/silk

用法: silk-encoder [选项]
  [选项]
    -l <语言路径>               指向 po/mo 文件或所在文件夹
    -i <输入文件>               待编码的输入语音文件
    -o <输出文件>               编码后的文件
    -Fs_API <采样率>            单位赫兹(Hz), 默认值为 24000
    -Fs_maxInternal <赫兹>      最大采样率，单位赫兹(Hz), 默认值为 24000
    -packetlength <毫秒>        数据包长度，单位毫秒(ms), 默认值为 20
    -rate <比特率>              比特率，默认值为 25000
    -loss <损耗比>              上行链路预计损耗比例，取值(0-100), 默认值为 0
    -inbandFEC[=false]          开启音频带内 FEC(前向纠错), 默认值为 false
    -complexity <模式>          设置复杂模式, 0=低，1=中，2=高，默认值为 2
    -DTX[=false]                开启 DTX, 默认值为 false
    -stx[=false]                在文件头之前添加 STX 标记，并移除 footer 块(兼容国内通信软件语音格式), 默认值为 true
```

## See also 致谢
- https://github.com/gaozehua/SILKCodec    源码
- https://github.com/kn007/silk-v3-decoder 兼容国内软件的版本
- https://github.com/wdvxdr1123/go-silk    ccgo 转写为 go 的版本
- https://github.com/zxfishhack/go-silk    可直接转 wav 的版本
- [Go语言高级编程 - 第 2 章 CGO 编程](https://chai2010.cn/advanced-go-programming-book/ch2-cgo/index.html)
- https://github.com/sunicy/go-lame 将 pcm 转换为 mp3，适配 go mod 的 fork 版本：https://github.com/youthlin/go-lame

## LICENSE
MIT.

C 源码开源协议见每个文件头部注释。
