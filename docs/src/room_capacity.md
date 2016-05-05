## 获取一个房间的容量

### 请求

`POST /api/v1/room/capacity`

### 参数

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Name | string | 房间名 |

### 响应

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Capacity | int | 最大可容人数 |
| Used | int | 当前人数 |
| Idle | int | 空闲数 Idle=Capacity-Used |

### 例子

```
$ curl -s -X POST -d '{"Name": "room_name"}' https://lawsroom.com/api/v1/room/capacity
```
