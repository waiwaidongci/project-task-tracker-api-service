# Bug Reproduce

## Bug 是什么

任务列表请求 `page_size=100` 时只返回 50 条，但响应里的 `page_size` 仍显示 100。

## 如何触发

```sh
go test ./internal/service -run TestTaskListHonorsRequestedPageSize
```

## 错误信息

`expected page_size=100 and 60 items, got page_size=100 len=50`
