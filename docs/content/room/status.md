---
date: 2016-05-14T03:43:01Z
title: 获取一个房间的状态
---

## 请求

```
POST /v1/room/status
```

## 参数

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Name | string | 房间名 |

## 响应

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Capacity | int | 最大可容人数 |
| Used | int | 当前人数 |
| Idle | int | 空闲数 Idle=Capacity-Used |

## 示例

```
$ curl -s -X POST -d '{"Name": "room_name"}' https://api.lawsroom.com/v1/room/status
```
