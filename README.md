# ServiceComputer-Selpg
## 1.关于Selpg
使用Golang开发Linux命令行中的selpg，selpg是与 cat、ls、pr 和 mv 等标准命令类似的 Linux 命令行实用程序，这个名称代表 SELect PaGes。selpg 允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件或另一个进程。通用 Linux 实用程序的编写者应该在代码中遵守某些准则。这些准则经过了长期发展，它们有助于确保用户以更灵活的方式使用实用程序，特别是在与其它命令（内置的或用户编写的）以及 shell 的协作方面 ― 这种协作是利用 Linux 作为开发环境的能力的手段之一。</br>
## 2.实验准备
### 2.1.安装和导入pflag
使用go get命令来安装pflag包。</br>
![](https://github.com/liuhq7/ServiceComputer-Selpg/blob/master/0~HMG205QVEP%7B~NAS8ZT%25Z2.png)</br>
我们通过go test来确认是否导入文件包。</br>
![](https://github.com/liuhq7/ServiceComputer-Selpg/blob/master/MOSXM_9%24J%7BV9%256)1RE3T%5D%5DY.png)</br>
这样我们就可以在代码中导入我们想要的关于pflag的包了。flag和pflag都是源自于Google，工作原理甚至代码实现基本上都是一样的。flag虽然是Golang官方的命令行参数解析库，但是pflag却得到更加广泛的应用，因为支持更精细的变量类型和更丰富的功能。</br>
### 2.2.引入所需要的包
```go
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
```
io，实现了一系列非平台相关的 IO 相关接口和实现，比如提供了对 os 中系统相关的 IO 功能的封装。我们在进行流式读写（比如读写文件）时，通常会用到该包。</br>
os/exec，执行外部命令，它包装了 os.StartProcess 函数以便更容易映射到 stdin 和 stdout，并且利用 pipe 连接 I/O。</br>
bufio，在 io 的基础上提供了缓存功能。在具备了缓存功能后， bufio 可以比较方便地提供 ReadLine 之类的操作。</br>
os，提供了对操作系统功能的非平台相关访问接口。接口为Unix风格。提供的功能包括文件操作、进程管理、信号和用户账号等。</br>
fmt，实现格式化的输入输出操作，其中的 fmt.Printf() 和 fmt.Println() 是开发者使用最为频繁的函数。</br>
pflag，提供命令行参数的规则定义和传入参数解析的功能。绝大部分的 CLI 程序都需要用到这个包。</br>
### 2.3.定义结构体
```go
type selpgArgs struct {
	// selpg -s startPage -e endPage [-l linePerPage | -f ][-d dest] filename
	startPage  int    //开始页码
	endPage    int    //结束页码
	filename   string //文件名
	pageLength int    // 每页长度
	pageType   bool   // 页类型
	printDest  string // 输出地址
}
```
## 3.函数部分
### 3.1.main函数
```go
func main() {
	selpg := selpgArgs{0, 0, "", 60, false, ""}
	progressName = os.Args[0]
	InitArgs(&selpg)
	pflag.Parse()
	ProcessArgs(len(os.Args), &selpg)
	ProcessInput(&selpg)
}
```
### 3.2.用pflag包进行对Selpg程序中对应的变量的参数绑定
```go
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
```
### 3.3.检查参数的正确性
```go
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
```
### 3.4.关于文件的读和写
```go
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
```
## 4.程序测试
![](https://github.com/liuhq7/ServiceComputer-Selpg/blob/master/E%5DHN85D%7D%7D%5D%7BHG_BPOB6%40E%60I.png)</br>
![](https://github.com/liuhq7/ServiceComputer-Selpg/blob/master/BF%24FVX)%241MFI%5B%7B4L%40%7DT%606X5.png)</br>
