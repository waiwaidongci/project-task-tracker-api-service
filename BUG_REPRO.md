# Bug Reproduce

## Bug 是什么

请求上下文已经取消后，任务列表仍会返回数据，取消信号没有向下传递。

## 如何触发

```sh
go test ./internal/service -run TestTaskListStopsWhenContextIsCanceled
```

## 错误信息

`expected canceled context to stop task list`
