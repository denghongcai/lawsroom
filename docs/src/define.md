## 约定

* 接口所有的请求URL前缀均为: `https://lawsroom.com`.
* 接口响应的HTTP状态码遵循标准HTTP协议. 通常应该为200.
* 响应如果带有body. 那么格式一定为 `json`. 结构如下:

    | 字段 | 类型 | 说明 |
    | --- | --- | --- |
    | error | null/object | 错误消息 |
    | result | null/mixed  | 期望的数据 |

    * 如果error不为null, 那么其结构如下:

        | 字段 | 类型 | 说明 |
        | --- | --- | --- |
        | code | int | 错误码 |
        | message | string  | 错误描述 |

    * 如果请求的接口不期望其返回任何数据, 则result为null, 否则result为期望的数据.
    * 下面的接口在描述响应时将只描述result.
