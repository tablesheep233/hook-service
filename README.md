# hook-service
git webhooks service

## About

使用go编写的git webhooks 服务，用于监听git发布的事件而执行运维脚本



## Getting Started

```cmd
set GOARCH=amd64 
set GOOS=linux
go build -o YOURPATH/server main.go
chmod +x server
./server
```

通过编写配置文件`config.yaml`进行配置

```yaml
port:
  17888         #使用端口
scripts:        #脚本路径 k:v对 目前key是 git推送的repository name, value是脚本的绝对路径
  pipe-detect: /home/deploy/install.sh
```

