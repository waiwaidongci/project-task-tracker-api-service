# Bug Reproduce

## Bug 是什么

给任务传入不存在的项目 ID 时，错误类型在传递中丢失，接口返回内部错误而不是 404。

## 如何触发

```sh
go test ./internal/service -run TestTaskServiceCreateRequiresProject
```

## 错误信息

`got project not found`
