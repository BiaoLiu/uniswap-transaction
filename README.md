# uniswap-transaction

## 起步

### 配置

所有配置信息使用Protobuf进行定义，具体使用请参考以下流程

##### 修改配置
1. 修改proto文件
```
cd app/uniswap/job/internal/conf
vim conf.proto
```

2. 编译proto文件，生成go代码
```
make config
```

3. 修改本地配置文件
```
vim app/uniswap/job/configs/config.yaml
```

### 依赖注入

#### 1.提供Provider(构造器)

Provider为负责创建对象的方法，如：NewData，NewRepo，NewService。在service.go、biz.go、data.go的ProviderSet集合中添加Provider。

#### 2.编写Injector(注入器)

Injector为负责根据对象的依赖，依次构造依赖对象，最终构造目的对象的方法，如：wire.go中的initApp方法。

#### 3.执行wire命令

执行`make wire`命令，生成wire_gen.go，通过自动生成代码的方式在**编译期**完成依赖注入。

### 编译运行
#### 编译
```
cd app/uniswap/job
make build
```

#### 运行
```
cd app/uniswap/job/cmd/server
go run main.go -conf ./configs
-conf: 指定配置文件路径
```

### 静态代码检查
```
make lint
```