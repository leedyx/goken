# goken

start /b main.exe

netstat -aon|findstr 38080

taskkill -f -pid 28532

```bash

docker run -d --name=goken -p 38080:38080 -v D:\Data\docker_data_path\goken\token:/opt/app/token  -v D:\Data\docker_data_path\goken\log:/opt/app/log goken:v1

```



