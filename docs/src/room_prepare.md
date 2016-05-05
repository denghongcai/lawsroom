## 获取一个新的空房间

### 请求

`GET /api/v1/room/prepare`

### 参数

无

### 响应

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Name | string | 房间名 |

> 提示: 获取成功后浏览器直接访问 `https://lawsroom.com/room/:Name` 即可进入房间

### 例子

```
$ curl -s -X GET https://lawsroom.com/api/v1/room/prepare
```
