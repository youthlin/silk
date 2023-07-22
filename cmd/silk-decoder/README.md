# Silk Decoder

## Install
```
go install github.com/youthlin/silk/cmd/silk-decoder@latest
# Execute to see usage
silk-decoder
```

## Usage
```
Silk decoder, Go version, based on v1.0.9 of C version
Decode silk v3 file to pcm, by youthlin
GitHub: https://github.comyouthlin/silk

Usage: silk-decoder -i <input file> -o <output file> [-l <language>] [-sampleRate <hz>] [-quiet[=false]]
  -i <input file>       Bitstream input to decoder
  -o <output file>      Speech output from decoder
  -l <language>         Language path(pointer to po file/dir)
  -sampleRate <hz>      Sample rate in Hz, default 24000
  -quiet[=false]        Print out just some basic values

Silk 解码器，Go 语言版本，基于 v1.0.9 的 C 语言版本
将 silk v3 格式的文件解码为 pcm, 作者：youthlin
GitHub: https://github.comyouthlin/silk

用法：silk-decoder [-l <language>] [-sampleRate <hz>]
  -i <input file>       要解码的输入文件 silk 格式
  -o <output file>      解码后的 pcm 语音文件
  -l <language>         指定语言路径(po 文件或文件夹)
  -sampleRate <hz>      采样率，默认值为 24000
  -quiet[=false]        静默模式

```

### 翻译
提取翻译字符串
```
xgettext -C -c=TRANSLATORS: --from-code=UTF-8 -o messages.pot -kT:1  *.go
```

生成翻译文件
```
msginit -i messages.pot -l zh_CN
```

重新提取后, 更新翻译文件
```
msgmerge -U zh_CN.po messages.pot
```
