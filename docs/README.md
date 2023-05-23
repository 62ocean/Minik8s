### 概述
此文件说明了miniK8s项目的代码规范、接口规范和命令行格式，供小组三人协作时进行交流记录。

### 代码规范
#### 代码日志
技术参考：https://zhuanlan.zhihu.com/p/159482200
全局的Log配置已在cmd/main.go init()函数中配置，prefix要在各个模块init函数中重新配置

- 使用Go语言内置Log包，flag属性配置为时间（精确到毫秒）+ 短文件名
- 日志输出前缀设置为英文方括号 + 组件名称（Service、Kubectl、Pod等）
- 日志输出在log文件夹中，文件名为{启动时间}.log
```
log.SetFlags(log.Lshortfile | log.Lmicroseconds)
log.SetOutput(logFile)

log.SetPrefix("[Pod]")
```
#### 错误处理
技术参考：https://juejin.cn/post/7197361396220100663

约定每个函数最后一个返回值为err
使用defer panic recover进行异常检查与恢复，并将错误信息直接通过` fmt.Println`直接输出至终端

#### 注释与git推送
为方便阅读，统一使用中文

### 接口规范

### 命令行格式