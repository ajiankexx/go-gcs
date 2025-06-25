# Module 
| 模块             | 说明                                     |
| -------------- | -------------------------------------- |
| `cmd/`         | 应用启动主函数目录，便于多命令支持（例如 server / cli）     |
| `api/handler/` | 类似 Java 的 `controller/`，负责路由处理和参数绑定    |
| `api/router/`  | 使用 Gin / Echo 注册路由、绑定中间件               |
| `service/`     | 核心业务逻辑（如认证、权限、仓库管理）                    |
| `dao/`         | 对数据库的 CURD 操作（结合 gorm/原生 SQL）          |
| `model/`       | 所有数据结构体，包括请求 DTO、返回 VO、数据库 PO          |
| `middleware/`  | JWT 校验、错误恢复等中间件                        |
| `config/`      | 统一配置管理（使用 `viper` 加载 `.yaml` 或 `.env`） |
| `util/`        | 常用工具封装（如加密、UUID、时间转换）                  |
| `pkg/`         | 公共模块，如标准响应格式封装、统一错误码定义                 |


