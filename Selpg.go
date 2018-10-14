package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/pflag"
)

const MAX_INT = int(^uint(0) >> 1)

type selpgArgs struct {
	// selpg -s startPage -e endPage [-l linePerPage | -f ][-d dest] filename
	startPage  int    //开始页码
	endPage    int    //结束页码
	filename   string //文件名
	pageLength int    // 每页长度
	pageType   bool   // 页类型
	printDest  string // 输出地址
}

var progressName string

func main() {
	selpg := selpgArgs{0, 0, "", 60, false, ""}
	progressName = os.Args[0]
	InitArgs(&selpg)
	pflag.Parse()
	ProcessArgs(len(os.Args), &selpg)
	ProcessInput(&selpg)
}

func InitArgs(p *selpgArgs) {
	pflag.Usage = Usage
	pflag.IntVarP(&p.startPage, "start", "s", 1, "Start Page Number")
	pflag.IntVarP(&p.endPage, "end", "e", 1, "End Page Number")
	pflag.IntVarP(&p.pageLength, "LinePerPage", "l", 60, "Line in Each Page")
	pflag.BoolVarP(&p.pageType, "printType", "f", false, "flag Splits Page")
	pflag.StringVarP(&p.printDest, "printDest", "d", "", "Destinaton of Print")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage:\tselpg -s startPage -e endPage [-l linePerPage | -f ][-d dest] filename\n")
	fmt.Fprintf(os.Stderr, "\t-s=Number\tStart Page Number\n")
	fmt.Fprintf(os.Stderr, "\t-e=Number\tend Page Number\n")
	fmt.Fprintf(os.Stderr, "\t-l=Number\tLine in Each Page\n")
	fmt.Fprintf(os.Stderr, "\t-f=BoolValue(true or false)\tflag Splits Page\n")
	fmt.Fprintf(os.Stderr, "\t-d=string\tDestinaton of Print\n")
	fmt.Fprintf(os.Stderr, "\t[filename]\t\n")
}

func ProcessArgs(ac int, selpg *selpgArgs) {
	if ac < 3 {
		fmt.Fprintf(os.Stderr, "Not enough arguments\n")
		pflag.Usage()
		os.Exit(1)
	} else if os.Args[1][0] != '-' && os.Args[1][1] != 's' {
		fmt.Fprintf(os.Stderr, "First arg should be -s=Number\n")
		pflag.Usage()
		os.Exit(2)
	} else if selpg.startPage < 1 || selpg.startPage > (MAX_INT-1) {
		fmt.Fprint(os.Stderr, "StartPage Number invalid\n")
		pflag.Usage()
		os.Exit(3)
	} else if os.Args[3][0] != '-' && os.Args[3][1] != 'e' {
		fmt.Fprintf(os.Stderr, "Second arg should be -e=Number\n")
		pflag.Usage()
		os.Exit(4)
	} else if selpg.endPage < 1 || selpg.endPage > (MAX_INT-1) || selpg.startPage < selpg.endPage {
		fmt.Fprintf(os.Stderr, "EndPage Number invalid\n")
		pflag.Usage()
		os.Exit(5)
	} else if selpg.pageLength < 1 {
		fmt.Fprintf(os.Stderr, "PageLength invalid\n")
		pflag.Usage()
		os.Exit(6)
	}

	if pflag.NArg() > 0 {
		selpg.filename = pflag.Arg(0)
		file, err := os.Open(selpg.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "The filename doesn't exits\n")
			os.Exit(7)
		}
		file.Close()
	}
}

func ProcessInput(selpg *selpgArgs) {
	fout := os.Stdout
	result := ""
	pageCount := 0
	lineCount := 0
	reader := bufio.NewReader(os.Stdin)
	if selpg.filename != "" {
		file, err := os.Open(selpg.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not open the file\n")
			os.Exit(8)
		}
		defer file.Close()
		reader = bufio.NewReader(file)		
	}

	if selpg.printDest != "" {
		file, err := os.OpenFile(selpg.printDest,os.O_RDWR,0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not open the destination file!\n")
			os.Exit(9)
		}
		defer file.Close()
		fout = file
	}
	if selpg.pageType == true {
		pageCount = 1
		for {
			str, err := reader.ReadString('\f')
			if err == io.EOF || pageCount > selpg.endPage {
				break
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "Read File Error1!\n")
				os.Exit(10)
			}
			pageCount++
			fout.Write([]byte(str))
			result = strings.Join([]string{result, str}, "")
		}
	} else {
		pageCount = 1
		lineCount = 0
		for {
			str, err := reader.ReadString('\n')
			if err == io.EOF || pageCount > selpg.endPage {
				break
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "Read File Error2!\n")
			}
			lineCount++
			if lineCount > selpg.pageLength {
				lineCount = 1
				pageCount++
			}
			fout.Write([]byte(str))
			result = strings.Join([]string{result, str}, "")
		}
	}
	if selpg.printDest != "" {
		cmd := exec.Command("lp", "-d"+selpg.printDest)
		cmd.Stdin = strings.NewReader(result)
		var stderrOut bytes.Buffer
		cmd.Stderr = &stderrOut
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, stderrOut.String())
		}
	}
}
