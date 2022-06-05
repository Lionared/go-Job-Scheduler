# Job Scheduler

## 编译  
* linux  
```shell
make linux
```  
* macos  
```shell
make macos
```  
* docker  
```shell
make docker
```  

## 运行  
* linux  
```shell
./dist/goscheduler-linux -p 20001 --store-type=redis --store-host=127.0.0.1 --store-port=6379 --store-password=123456
```
* macos  
```shell
./dist/goscheduler-darwin -p 20001 --store-type=redis --store-host=127.0.0.1 --store-port=6379 --store-password=123456
```  
* docker  
```shell
docker run --rm -p 20001:20001 jobscheduler /dist/goscheduler-docker -h 0.0.0.0 -p 20001 --store-type=redis --store-host=192.168.5.108 --store-port=6379 --store-password=123456
```