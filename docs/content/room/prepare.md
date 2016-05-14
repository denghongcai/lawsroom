---
date: 2016-05-14T03:07:13+08:00
title: 获取一个新的空房间
---

## 请求

```
GET /v1/room/prepare
```

## 参数

无

## 响应

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| Name | string | 房间名 |

> 提示: 获取成功后浏览器直接访问 `https://lawsroom.com/room/:Name` 即可进入房间

## 示例

```
$ curl -s -X GET https://api.lawsroom.com/v1/room/prepare
```
