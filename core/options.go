package core

import (
	"CuiRi/gologger"
	"flag"
	"os"
)

type Options struct {
	FileName	string	//	shellcode文件名
	Manual		bool	//	shellcode生成方法
	Stdin		bool
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.FileName, "f", "","通过shellcode生成免杀马")
	flag.BoolVar(&options.Manual, "manual", false, "查看shellcode生成方法")
	flag.Parse()

	options.Stdin = hasStdin()
	ShowBanner()

	if len(options.FileName) == 0 && !options.Manual {
		flag.Usage()
		os.Exit(0)
	}
	if options.FileName != "" && !FileExists(options.FileName) {
		gologger.Fatalf("文件 %s 不存在!\n", options.FileName)
		os.Exit(0)
	}

	return options
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}