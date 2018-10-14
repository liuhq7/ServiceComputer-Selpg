# ServiceComputer-Selpg
## 1.关于Selpg
使用Golang开发Linux命令行中的selpg，selpg是与 cat、ls、pr 和 mv 等标准命令类似的 Linux 命令行实用程序，这个名称代表 SELect PaGes。selpg 允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件或另一个进程。通用 Linux 实用程序的编写者应该在代码中遵守某些准则。这些准则经过了长期发展，它们有助于确保用户以更灵活的方式使用实用程序，特别是在与其它命令（内置的或用户编写的）以及 shell 的协作方面 ― 这种协作是利用 Linux 作为开发环境的能力的手段之一。</br>
## 2.实验准备
### 2.1.安装和导入pflag
使用go get命令来安装pflag包。</br>
![]()
我们通过go test来确认是否导入文件包。</br>
![]()
这样我们就可以在代码中导入我们想要的关于pflag的包了。flag和pflag都是源自于Google，工作原理甚至代码实现基本上都是一样的。flag虽然是Golang官方的命令行参数解析库，但是pflag却得到更加广泛的应用，因为支持更精细的变量类型和更丰富的功能。</br>
### 2.2.引入所需要的包
![]()
io，实现了一系列非平台相关的 IO 相关接口和实现，比如提供了对 os 中系统相关的 IO 功能的封装。我们在进行流式读写（比如读写文件）时，通常会用到该包。</br>
os/exec，执行外部命令，它包装了 os.StartProcess 函数以便更容易映射到 stdin 和 stdout，并且利用 pipe 连接 I/O。</br>
bufio，在 io 的基础上提供了缓存功能。在具备了缓存功能后， bufio 可以比较方便地提供 ReadLine 之类的操作。</br>
os，提供了对操作系统功能的非平台相关访问接口。接口为Unix风格。提供的功能包括文件操作、进程管理、信号和用户账号等。</br>
fmt，实现格式化的输入输出操作，其中的 fmt.Printf() 和 fmt.Println() 是开发者使用最为频繁的函数。</br>
pflag，提供命令行参数的规则定义和传入参数解析的功能。绝大部分的 CLI 程序都需要用到这个包。</br>
### 2.2.定义结构体
![]()
## 3.函数部分
### 3.1.main函数
![]()
### 3.2.用pflag包进行对Selpg程序中对应的变量的参数绑定
![]()
### 3.3.检查参数的正确性
![]()
### 3.4.关于文件的读和写
![]()
![]()