# Bug Reproduce

## Bug 是什么

创建带标签的任务时，标签去重逻辑写入了一个未初始化的 map，导致 panic。

## 如何触发

```sh
go test ./internal/service -run TestCreateTaskWithTagsDoesNotPanic
```

## 错误信息

`assignment to entry in nil map`
