# gin-example

gin应用示例

## tz-gin-cli
tz-gin-cli的热更新`tz-gin run`

## 快速预览

```
tz-gin run
```

## 目录结构

```
├─config       //配置文件
├─controller   //所有与HTTP请求相关的业务逻辑都放在controller层中
├─middleware   //中间件
├─model        //模型
├─router       //路由
├─service      //服务
│  └─validator //自定义数据校验
└─sql
```

## 示例代码

tz-gin提供了部分简单的示例代码，其放在`*-example.go`下，作为快速上手的示例

## 数据库

数据库sql文件为 `gin-example.sql` ，请设置环境变量以指定数据库

## 环境变量

环境变量名及其默认值在 `util/config.go` 中定义,为了方便，**在开发过程中**，也可以通过`.env`文件配置相关参数，项目开发时应拷贝`.env.example`为`.env`

- 项目**实际**上线时， `APP_PROD` 应设置为任意非空字符串，以开启生产模式
- 项目**实际**上线时， `APP_SECRET` 应设置为各应用互不相同的字符串并保密
- 项目**实际**上线时， `APP_ALLOW_HEADERS` `APP_ALLOW_ORIGINS` 应设置来防止存在的跨域`CORS`风险，如果有多个则使用`|`分开

## 日志

在生产模式下，日志会输出到 `log` 目录下
