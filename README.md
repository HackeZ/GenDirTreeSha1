## GenDirTreeSha1

一个生成目录树哈希的小工具

1. 用golang开发，代码放到github上，用github进行问题跟踪
2. 对整个目录下的所有文件进行遍历，获取所有文件的大小和计算文件的sha1哈希值，记录在一个文件里面
3. 结果文件格式：每一行一个文件，用逗号隔开，前面是文件名称，后面是哈希值，文件大小
4. 需要可以指定忽略哪些目录、文件，需要支持通配符
5. 代码实现简洁，运行性能高得分高
6. 要求通过测试代码自我证明代码能够可靠运行并正确实现上述功能
要求写安全稳定可靠的代码

### Usage 

```go
$ go build -o genDirTreeSha1 main.go
$ ./genDirTreeSha1 -r "xxx dir root" -d ".git,vie?s,*s" -f "*.go,*.t?t,*.*"
```