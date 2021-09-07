package core

import (
	"CuiRi/gologger"
	"github.com/gookit/color"
)

const Version = "1.0"
const Banner = `
   ____   __    __    _____      ______      _____  
  / ___)  ) )  ( (   (_   _)    (   __ \    (_   _) 
 / /     ( (    ) )    | |       ) (__) )     | |   
( (       ) )  ( (     | |      (    __/      | |   
( (      ( (    ) )    | |       ) \ \  _     | |   
 \ \___   ) \__/ (    _| |__    ( ( \ \_))   _| |__ 
  \____)  \______/   /_____(     )_) \__/   /_____( 

摧日：一款红队专用免杀木马生成器，基于shellcode生成绕过所有杀软的木马
`
const LinkAndAuthor = "https://github.com/NyDubh3/CuiRi   Author:Dubh3\n"
const Warning = "警告：\n1.本工具仅用于企业内部测试，请勿用于任何非法犯罪活动，否则后果自负\n2.本工具需要Go语言环境，且使用时需要关闭杀软"

func ShowBanner() {
	color.RGBStyleFromString("210,105,30").Println(Banner)
	color.RGBStyleFromString("255,0,0").Println(Warning)
	color.RGBStyleFromString("30,144,255").Println(LinkAndAuthor)
	//gologger.Printf(Banner)
	gologger.Infof("Current Version: %s\n", Version)
}