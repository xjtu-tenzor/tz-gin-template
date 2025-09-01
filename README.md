# gin-example

gin应用示例

## tz-gin-cli

tz-gin-cli的热更新 `tz-gin run`

## 快速预览

```
tz-gin run
```

## 目录结构

```
├─common     	//各层级都会复用的结构体及函数，比如错误处理
├─config       	//配置文件
├─controller   	//所有与HTTP请求相关的业务逻辑都放在controller层中
├─middleware   	//中间件
├─model        	//模型
├─pkg        	//额外的功能的实现的包
├─router       	//路由
├─service      	//服务
│  └─validator 	//自定义数据校验
└─sql
```

## 示例代码

tz-gin提供了部分简单的示例代码，其放在 `*-example.go`下，作为快速上手的示例

## 数据库

数据库sql文件为 `gin-example.sql` ，请设置环境变量以指定数据库，该sql文件只存储基本表结构而不应存在数据，数据库的更新则应以增量更新的方式 `gin-example-unix_timestamp.up.sql`

## 环境变量

环境变量名及其默认值在 `config/config.go` 中定义,为了方便，**在开发过程中**，也可以通过 `.env`文件配置相关参数，项目开发时应拷贝 `.env.example`为 `.env`

- 项目**实际**上线时， `APP_PROD` 应设置为任意非空字符串，以开启生产模式
- 项目**实际**上线时， `APP_SECRET` 应设置为各应用互不相同的字符串并保密
- 项目**实际**上线时， `APP_ALLOW_HEADERS` `APP_ALLOW_ORIGINS` 应设置来防止存在的跨域 `CORS`风险，如果有多个则使用 `|`分开

## 日志

日志共有4种主要模式， `debug`、`info`、`warn`、`error`， 当然有隐藏模式 `trace` ，需要使用钩子开启。
可以在env文件里配置APP_LOG_LEVEL里面设置，默认`info`。注意日志的模式和设置项目是否为生产模式没有关系。

在生产模式下，日志会输出到 `log` 目录下。默认的日志记录信息包含了 `url,method,client_ip`等信息，如果发生错误则会输出包含栈信息的日志。
如果Gin处于Debug模式，日志将会同步输出到 `stdout`中。日志默认捕获GinLogger和GinRecovery, Gorm以及stderr信息。

**记录日志**：
项目开发过程中，建议使用 `logger.Infof`等方式输出，而不是使用 `fmt.Print `，具体用法可见 `logger/gin.go`,还有别用那个logger.database那个，那是给gorm兼容用的

**自定义日志输出**：

可以在 `config/logconf`里面的的 `Format`修改输出颜色等格式约定，也可以使用 `logger.WithFields`这个方法添加自定义字段。

**日志钩子**

需要实现 `logrus.Hook`，实现两个方法为 `Levels`和 `Fire`，例子可以见 `config/loghook`，通过使用方法 `logger.AddHook(&CustomHook{})`来为 `logger`添加钩子。

**可选择使用的日志钩子（默认关闭）**

1.TraceHook：将debug模式转换为trace模式，输出堆栈信息。

2.RemoteHook：将日志推送到远程服务器.

# 项目开发规范

#### tz-gin受到 koa 洋葱模型的思想，故设计如下的项目开发规范

## session

使用 `controller/session.go`下提供的函数进行session的处理，session的密钥应在**生产环境**中通过**环境变量**形式传入 `APP_SECRET`

## model

- `model` 中定义了与数据库相对应的模型，请在结构体的各字段中详细的写出相关的 `tag`
- 在 `model.go` 中提供了 `baseModel` ，在声明模型是应该包含该结构体
- 在 `scopes.go` 中提供了一些基础常见的服用逻辑，同时，在项目中，你也应该将一些复用通用的逻辑写在此处

## controller 的注册方式

将对相同资源处理的方法绑定在同一个结构体上，详情可见示例 `controller/hello-example.go`

同时将该结构体注册到 `controller/controller.go`的 `Controller`结构体中

## service 的注册方式

注册方式同controller相同，将对相同资源处理的方法绑定在同一个结构体上，示例见 `service/hello-example.go`

同时将该结构体注册到 `service/service.go`的 `Service`结构体中

## controller & service

- ~~在 `controller` 中对应的 `.go` 文件下，构造一个函数将 `Request` 绑定为 `model` 中的结构体~~
- 将该结构体传入 `service` ，`session` `page` `limit` 等作为其他参数传入
- `service` 的第一个返回值应为返回体中的 `data` 部分，`Response` 应在 `service` 中定义

## 错误异常处理

在出现非法操作和调用方法出现错误时，应调用 `common` 包下面的 `ErrNew` 方法并将错误返回到 `controller` 层统一处理（在边界情况明晰的情况下）,下面时 `ErrNew` 的函数原型及相关错误码

```go
func ErrNew(err error, errType gin.ErrorType) error

const (
	ParamErr gin.ErrorType = iota + 3   //参数错误
	SysErr                              //系统错误
	OpErr                               //操作错误
	AuthErr                             //鉴权错误
	LevelErr                            //权限错误
)
```

当你想自定义错误码时，请与前端进行沟通

## 自定义校验规则书写方式

- 校验规则应书写两个函数，一个函数为校验规则判断函数，另一个函数为翻译函数，俩函数原型如下

  ```go
  func(fl FieldLevel) bool		// 校验函数
  func(ut ut.Translator) error	// 翻译函数

  ```
  以下是我们提供的示例代码，也可在 `service/validator/validators.go` 和 `service/validator/translations.go` 看到

  ```go
  // 校验函数
  func timing(fl validator.FieldLevel) bool {			// 自定义应满足的时间
  	if date, ok := fl.Field().Interface().(time.Time); ok {
  		today := time.Now()
  		if today.After(date) {
  			return false
  		}
  	}
  	return true
  }


  // 翻译函数
  func timingTransZh(ut ut.Translator) error {
  	return ut.Add("timing", "{0}输入的时间不符合要求", true) // {0}表示会替代加了该校验的字段
  }
  ```
- 自定义校验的注册应放在 `service/validator/init.go` 的 `validatorHandleRouter` 中，`key` 值表示的是自定义校验的名称
- 校验规则应写在 `validators.go` 下面，翻译1则应写在 `translations.go` 下面

## 关于对函数式编程的支持
详细内容见 [pkg和函数式文档](pkg/README.md)

## 现有问题后续更新内容
- 后续更新添加更多pkg内容，也可能合并service和controller，毕竟现在写个东西改三五个地方太麻烦了
- session 问题已修复， 并修复了middleware/log的一个warning