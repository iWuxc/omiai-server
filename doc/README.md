# 文档中心

* [micros脚手架的安装与使用](http://git.microdreams.com/golang-common/micros)
* [go-kit使用文档](http://github.com/iWuxc/go-wit)
* [项目运行问题排查指南](FAQ.md)

## 项目开发约定

> 为保证项目开发中的一致性以及可维护性, 请遵循以下规范来编写代码

### 代码边界

主要边界限定:

1. internal/api 层提供全局通用的能力, 其他地方都可以依赖
2. internal/biz 层定义数据资源接口, 提供数据的持久化操作
3. internal/data 层是对 internal/biz 层接口的具体实现, 直接管理数据内容

### 项目规范

1. 减少依赖的层层传递, 例如: gin.Context 仅应该在 **internal/controller** 中依赖
2. 所有依赖的第三方应该在 **internal/server** 中统一初始化, 并且在需要的地方声明依赖, ([go-kit](http://github.com/iWuxc/go-wit) 提供了开箱即用的方式,
   因此不需要在 server 中初始化)
3. 所有接口的相应都应该由 **pkg/response** 统一返回, 并且在响应结果中携带 code (0: 成功, 非0: 错误代码)
4. 减低数据流转周期: 简单逻辑可直接在 internal/controller 中返回, 如果有复杂业务逻辑, 可以写在 internal/service 中, 并由 internal/controller 引入
5. 项目中所有的配置项, 都应该在 **configs/config.yaml.bak** 中明确写出
