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
$ ./genDirTreeSha1 -r `the-dir-root` -d ".git,vie?s,*s" -f "*.go,*.t?t"
```

### Dev-Log

基本功能已经实现，这是一个真正完全并发同步的版本，计算每一个文件的 SHA1 值都是一个单独的 Groutine，再也不用怕大文件啦！

异步并发版本最大的问题是每次输出结果的顺序都是不一样的..

-------

目前有一个很大的问题就是没有限制 Goroutine 的上限，如果同时处理大量容量很大的文件，栈中的空间很容易就将内存塞满。
处理这种情况有两种办法：
    - 添加多核支持，使密集的 SHA1 值计算加快，减少阻塞的 Goroutine
    - 可以设置默认的 Goroutine 上限，另外用一个 Channel 进行同步，不让内存中存在无上限个 G。

结果：

![MultiCPU vs SingleCPU](http://7xsxev.com1.z0.glb.clouddn.com/GenDirTreeSHA1%20MultiCPU%20vs%20SingleCPU.png)

跑出 `21s` 是上限设置为 `25600` 个 Goroutine 并且没有开启多核支持的结果。
而 开启了多核支持，而且 Goroutine 上限是默认的 `1024` 个，却能跑出 14s 的成绩，同时也大幅度减少了栈空间所占内存。

由此可见，优化之后带来的结果是可喜的。
