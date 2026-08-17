# Bug Reproduce

## Bug 是什么

按状态、优先级或项目等条件筛选任务列表时，查询参数会被组装错，数据库侧报缺少参数。

## 如何触发

创建一个项目和一条待处理任务，再按 status 筛选并分页列出任务：

```sh
go test ./internal/repository -run TestTaskRepositoryFiltersAndStats
```

## 错误信息

`missing argument with index 3`
