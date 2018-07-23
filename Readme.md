# Debug Service

## 流程介绍

1. 启动时，从项目根目录读取`go-online.yml`文件内容，并创建`Debug`文件目录
3. 系统在项目根目录查找`Makefile`文件，若未找到，则结束调试；若找到，则生成可执行文件
4. 启动调试进程，从`Debug`目录载入可执行文件`main`，若未找到，则结束调试
5. 从`stdin`中以行缓冲方式获得输入，并直接传入gdb进程
6. 从gdb进程获取输出，并格式化成json后，输出到stdout
7. 程序收到退出指令后，先发送停止指令给gdb进程，等待5s之后终止自身，从而强制gdb和容器退出

## 注意事项

1. 当前只支持c/c++语言的调试
2. 默认当前目录即为项目根目录
3. 项目运行会自动创建`Debug`目录以存放调试过程中产生的文件
4. 调试程序的运行目录为项目根目录
5. 默认会将所有的cpp/c文件作为编译的对象
6. gdb service默认接受的数据都是带有`\n`的