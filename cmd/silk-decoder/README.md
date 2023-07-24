# Silk Decoder

## Install
> You may need to run `sudo apt-get install libmp3lame-dev` to get lame lib on Linux

```
go install github.com/youthlin/silk/cmd/silk-decoder@latest
# Execute to see usage
silk-decoder
```

## Usage
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

```

支持中文，如果你的环境默认语言不是中文，可以设置 `LANG=zh` 再启动软件

`LANG=zh silk-decoder`

```
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
$s -i a.amr
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

### 翻译
提取翻译字符串
```
xgettext -C -c=TRANSLATORS: --from-code=UTF-8 -o messages.pot -kT:1 -kX:2,1c  *.go
```

生成翻译文件
```
msginit -i messages.pot -l zh_CN
```

重新提取后, 更新翻译文件
```
msgmerge -U zh_CN.po messages.pot
```
