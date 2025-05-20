# Api层
1. 每一个接口都应编写请求schema，以便进行参数校验
2. 响应错误码必须进行软编码

# Service层
1. Service都应包含context上下文

# Conf层
1. toml文件管理全局应用配置、各服务配置

# Common层
1. 请求响应格式化、请求参数校验
2. 统一错误码
3. Redis服务
4. Log服务
5. Util工具集
6. I18n对错误码进行国际化

# Middleware层
1. 实现`AuthMiddleWare`鉴权中间件，采用Redis+Token方式
2. 实现`Cors`跨域中间件
3. 实现`I18n`国际化中间件

# Model层
1. 每一个`Dao`都应继承`BaseDao`
2. 每一个`Model`根据是否需要软删，选择继承`BaseModel`或`SoftDeleteBaseModel`
3. 每一个子类`Dao`都应实现以下方法：
    - `GetByID(id int64) (*Model, error)`: 根据主键 ID 获取记录
    - `Select(where Where, appends Appends) ([]Model, error)`: 查询数据
    - `Insert(model *Model) (int64, error)`: 新增记录
    - `Update(model Model) (bool, error)`: 根据主键 ID 更新记录
    - `Delete(id int64) (bool, error)`: 根据主键 ID 删除记录（硬删/软删）
    - <font color="red">注：
        - `Insert`应向上层暴露最后一条插入记录的主键 ID
        - `Update`和`Delete`应向上层暴露是否更新/删除成功</font>