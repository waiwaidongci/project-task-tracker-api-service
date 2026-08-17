# 项目：待办事项API

从0做一个Go待办事项API，用Gin开发，数据存SQLite。支持维护项目和任务两级数据：项目包含名称和描述，任务关联项目，包含标题、优先级、状态、截止日期和标签。支持新增修改删除、按状态和优先级筛选、查看今日到期任务、统计各项目未完成任务数。代码按Go企业分层结构组织：cmd/server/main.go、internal/config、internal/model、internal/repository、internal/service、internal/handler、internal/router、internal/middleware、migrations。查询接口统一分页，状态变更通过service层校验，错误响应统一格式。
