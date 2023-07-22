# Silk Encoder

## Install
```
go install github.com/youthlin/silk/cmd/silk-encoder@latest
# Execute to see usage
silk-encoder
```

## Usage
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
    -inbandFEC                  Enable inband FEC usage, default: false
    -complexity <comp>          Set complexity, 0: low, 1: medium, 2: high; default: 2
    -DTX                        Enable DTX; default: false
    -quiet                      Print only some basic values
    -stx                        Add STX flag before file header and remove footer block, default true


Silk 编码器，Go 语言版本，基于 v1.0.9 的 C 语言版本
将 pcm 文件编码为 silk v3 类型，作者： youthlin
GitHub: https://github.comyouthlin/silk

用法: silk-encoder [选项]
  [选项]
    -l <path to po file>        指定语言路径(po 文件或文件夹)
    -i <input file>             待编码的输入语音文件
    -o <output file>t           编码后的文件
    -Fs_API <Hz>                采样率，单位赫兹(Hz), 默认值为 24000
    -Fs_maxInternal <Hz>        内部最大采样率，单位赫兹(Hz), 默认值为 24000
    -packetlength <ms>          数据包长度，单位毫秒(ms), 默认值为 20
    -rate <bps>                 比特率，默认值为 25000
    -loss <perc>                上行链路预计损耗比例，取值(0-100), 默认值为 0
    -inbandFEC                  开启音频带内 FEC(前向纠错), 默认值为 false
    -complexity <comp>          设置复杂模式, 0=低，1=中，2=高，默认值为 2
    -DTX                        开启 DTX, 默认值为 false
    -quiet                      只打印基本数据
    -stx                        在文件头之前添加 STX 标记，并移除 footer 块(兼容国内通信软件语音格式), 默认值为 true

```


### 翻译
提取翻译字符串
```
xgettext -C -c=TRANSLATORS: --from-code=UTF-8 -o messages.pot -kT:1  *.go
```

生成翻译文件
```
# GETTEXTCLDRDIR=./cldr
msginit -i messages.pot -l zh_CN
```

重新提取后, 更新翻译文件
```
msgmerge -U zh_CN.po messages.pot
```
