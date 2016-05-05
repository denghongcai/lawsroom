## 约定

### 请求

* 接口所有的请求URL前缀均为: `https://lawsroom.com`.
* 请求如果带有body, 那么格式必须为json.

### 响应

* 接口响应的HTTP状态码遵循标准HTTP协议. 通常应该为200.
* 正常情况下, 响应如果带有body, 那么格式为json. 结构如下:

    | 字段 | 类型 | 说明 |
    | --- | --- | --- |
    | Error | null/object | 错误消息 |
    | Result | null/mixed  | 期望的数据 |

    * 如果Error不为null, 那么其结构如下:

        | 字段 | 类型 | 说明 |
        | --- | --- | --- |
        | Code | int | 错误码 |
        | Message | string  | 错误描述 |

    * 如果请求的接口不期望其返回任何数据, 则Result为null, 否则Result为期望的数据.

    > 提示: 下面的接口在描述响应时将只描述Result.
