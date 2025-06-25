执行机制： 当 Gin 应用程序收到一个匹配 /users/create 路径的 POST 请求时，它并不会直接调用 userHandler.Create。相反，它会：

接收请求。
解析请求。
创建或获取一个上下文对象（例如 *gin.Context）。
然后，Gin 框架会“回调”或者说“执行”你传入的 userHandler.Create 函数，并将这个上下文对象作为参数传递给它。
