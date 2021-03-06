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
$ go get -u github.com/HackeZ/GenDirTreeSha1
$ cd $GOPATH/bin
$ ./GenDirTreeSha1 -r `the-dir-root` -d ".git,vie?s,*s" -f "*.go,*.t?t" -g 2048
```

### Dependce

- [alecthomas/kingpin](https://github.com/alecthomas/kingpin)

### Dev-Log

基本功能已经实现，这是一个真正完全并发同步的版本，计算每一个文件的 SHA1 值都是一个单独的 Groutine，再也不用怕大文件啦！

异步并发版本最大的问题是每次输出结果的顺序都是不一样的..

-------

目前有一个很大的问题就是没有限制 Goroutine 的上限，如果同时处理大量容量很大的文件，栈中的空间很容易就将内存塞满。

处理这种情况有两种办法：
    - 添加多核支持，使密集的 SHA1 值计算加快，减少阻塞的 Goroutine
    - 可以设置默认的 Goroutine 上限，另外用一个 Channel 进行同步，不让内存中存在无上限个 G。

--------

当前最新的版本将写入文件这个步骤都已经是并发，当计算出一个结果之后马上就会将其写入文件中，而不是像之前的将所有结果读取到一个变量之后再统一写入文件。

### Benchmark

```
400 Files => Used Times : 119ms    
10286 Files => Used Times : 1.02s
```
