# file-sync
同步文件到服务器

## 使用方式

### 编译
```shell
$ git clone https://github.com/Joeyscat/file-sync.git
$ cd file-sync
$ bash scripts/build.sh
$ ls output 
fs    # client
fsd   # server
``` 

### 在服务器执行fsd，等待上传
```shell
$ ./fsd -up path_for_upload -p 8002
```

### 在客户端执行fs，监听文件并进行上传
```shell
$ ./fs -dir path_to_sync -url http://server:8002/upload
```

之后在 path_to_sync 目录下的文件会被同步到服务端的 path_for_upload 目录下

注意，path_to_sync 的子目录文件也会上传，但是都会直接上传到 path_for_upload，而不是它的子目录


