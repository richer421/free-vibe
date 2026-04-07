# Kratos Layout Best Practices

## Overview

基于 Kratos 官方 Project Layout，并强化"业务域同构"设计：
同一个业务域在 `biz/data/service` 中必须保持同名目录、同构边界和一致职责。

参考资料：
- https://go-kratos.dev/zh-cn/blog/kratos/go-project-layout/
- https://go-kratos.dev/zh-cn/docs/intro/layout/

## 设计思想

1. 目录设计&语义命名（项目架构）
- 目录先表达业务边界，再表达技术分层；优先按业务域同构组织 `biz/data/service`，避免按技术细节平铺。
- 命名先表达领域语义，再表达实现细节；目录名、包名、类型名应能直接回答"这个模块负责什么业务能力"。
- "看目录即可理解系统架构"是硬约束：新成员不读实现代码，也能从目录和命名看出模块职责、依赖方向与扩展点。

## 术语可执行定义

1. 业务域同构（domain-homomorphic modules）
- 可执行定义：同一业务域必须在 `internal/biz/<domain>/`、`internal/data/<domain>/`、`internal/service/<domain>/` 同时存在，并且分别承载 `usecase/repo/handler` 职责。
- 判定标准：
  - 三层同名域目录都存在。
  - `service/<domain>` 只调用 `biz/<domain>`，不直接访问数据库或 `sqlc` 查询对象。
  - `data/<domain>` 提供 `biz/<domain>` 仓储接口的实现。
- 失败示例：
  - 只有 `service/<domain>`，缺失 `biz` 或 `data` 同名目录。
  - `service` 直接 import `database/sql` 或直接调用 `internal/data/gen/sqlc`。

2. 目录设计&语义命名（项目架构）
- 可执行定义：目录名/包名/核心类型名必须优先表达业务语义，其次才是实现细节。
- 判定标准：
  - 目录名使用 `greeter/order/user` 等业务语义名。
  - 核心类型命名可映射职责：`Usecase`、`Repo`、`Service`。
- 失败示例：
  - 使用 `misc`、`temp`、`helper2`、`newlogic` 等非语义目录名承载核心业务。
  - 以 `manager/service/helper` 泛化名替代领域名。

3. SQL 归域（domain-scoped SQL）
- 可执行定义：SQL 输入文件必须位于 `internal/data/<domain>/sql/{schema,query}`，禁止全局混放。
- 判定标准：
  - `sqlc.yaml` 的 `schema/queries` 指向域目录。
  - 域 repo 仅依赖该域 SQL 输入。
- 失败示例：
  - 使用 `internal/data/sql` 全局目录混放多域 SQL。
  - 某域 repo 依赖其他域 query 且无显式依赖说明。

4. 执行步骤（AI 必须按顺序）
1. 枚举 `internal/biz`、`internal/data`、`internal/service` 下一级业务域目录。
2. 检查每个域是否三层同名同构齐全。
3. 抽样检查 import 关系，确认 `service -> biz -> data` 未被绕过。
4. 检查 `sqlc.yaml` 与 SQL 目录是否按域归属。
5. 输出"通过/不通过 + 文件路径 + 修复动作"。

## 分层职责基线

1. `cmd/<app>`：启动、配置加载、DI 装配。
2. `api/...`：协议契约与错误语义。
3. `internal/service`：协议适配与调用编排。
4. `internal/biz`：业务规则、实体、用例、仓储接口。
5. `internal/data`：仓储实现、数据库/缓存/外部系统访问。
6. `internal/server`：HTTP/gRPC 注册与中间件装配。
7. `pkg`：仅用于跨项目复用。

## 每层写什么（可复制提示词）

1. `api/<domain>/v1`（契约层）
- 提示词模板：
`请在 api/<domain>/v1/<domain>.proto 定义 <domain> 的 gRPC 接口与 HTTP 注解，包含请求/响应 message、错误码语义，并保持字段命名稳定可演进。`

2. `internal/biz/<domain>`（用例层）
- 提示词模板：
`请在 internal/biz/<domain>/ 中定义 Entity、Repo 接口和 Usecase，只表达业务规则与编排，不依赖数据库实现和传输层细节。`

3. `internal/data/<domain>`（仓储实现层）
- 提示词模板：
`请在 internal/data/<domain>/ 实现 biz 层的 Repo 接口，封装数据库/缓存/外部依赖调用，处理技术错误到业务错误的映射，不写协议层逻辑。`

4. `internal/data/<domain>/sql`（SQL 输入层）
- 提示词模板：
`请在 internal/data/<domain>/sql/schema 与 query 中维护该域 SQL；schema 负责表结构，query 负责命名查询，确保 sqlc 可生成类型安全代码。`

5. `internal/service/<domain>`（协议适配层）
- 提示词模板：
`请在 internal/service/<domain>/ 实现 gRPC/HTTP handler，完成参数校验与 DTO/Entity 转换，仅调用 biz usecase，不直接访问数据库。`

6. `internal/server`（传输装配层）
- 提示词模板：
`请在 internal/server 注册 <domain> 的 HTTP/gRPC 服务与中间件，保持 server 只做 transport wiring，不包含业务判断。`

7. `cmd/<app>`（启动装配层）
- 提示词模板：
`请在 cmd/<app> 完成配置加载、logger 初始化、wire 依赖装配和 app 启停逻辑，禁止写业务规则。`

8. 聚合 Provider（域同构接线）
- 提示词模板：
`请更新 biz.go、service.go、wire.go，把 <domain> 的 provider 接入依赖图，确保 service -> biz -> data 单向依赖且可编译。`

## 业务域同构原则（核心）

1. 每个业务域在三层都建同名目录，例如 `greeter`：
- `internal/biz/greeter/`
- `internal/data/greeter/`
- `internal/service/greeter/`

2. 禁止"平铺式单文件域实现"长期存在，例如：
- `internal/biz/greeter.go`
- `internal/data/greeter.go`
- `internal/service/greeter.go`

3. 目录层面的同构优先于文件命名偏好：先保证域边界，再谈代码风格。

## 业务模块同构模板

```text
api/
└── <domain>/v1/
    ├── <domain>.proto
    └── *.pb.go (generated)

internal/
├── biz/
│   ├── biz.go                       # provider 聚合
│   └── <domain>/
│       └── <domain>.go              # Entity / Repo / Usecase
├── data/
│   ├── data.go                      # Data 资源容器
│   ├── gen/sqlc/                    # generated
│   └── <domain>/
│       ├── repo.go                  # Repo 实现
│       ├── wire.go                  # 域 provider
│       └── sql/
│           ├── schema/
│           └── query/
└── service/
    ├── service.go                   # provider 聚合
    └── <domain>/
        └── service.go               # gRPC/HTTP handler
```

## SQL 归属规则（必须）

1. SQL 必须归属业务域目录，不放在 `internal/data/sql` 这种全局混合目录。
2. 每个域维护自己的：
- `internal/data/<domain>/sql/schema`
- `internal/data/<domain>/sql/query`
3. `sqlc` 输出可统一到 `internal/data/gen/sqlc`，但输入必须按域拆分。
4. 跨域表查询必须在调用方域明确声明依赖，不允许在"公共 SQL 目录"隐式耦合。

## sqlc 配置建议

单域：
- `schema: internal/data/<domain>/sql/schema`
- `queries: internal/data/<domain>/sql/query`

多域：
- 在 `sqlc.yaml` 里为每个域声明一个 `sql` 条目；
- 或保持一个条目但以域目录分组输入路径，禁止回退到平铺全局 SQL 目录。

## 依赖方向与边界规则

1. 固定依赖方向：`service -> biz -> data`（通过 `biz` 接口反转实现依赖）。
2. `service` 不直连 DB/缓存/第三方 SDK。
3. `biz` 不依赖 transport 细节（HTTP/gRPC 上下文、路由对象）。
4. `server` 不承载业务规则。

## 新业务模块落地流程

1. 在 `api/<domain>/v1` 定义 proto 契约。
2. 在 `internal/biz/<domain>` 定义 `Entity/Repo/Usecase`。
3. 在 `internal/data/<domain>` 实现 repo，并放置该域 SQL。
4. 在 `internal/service/<domain>` 实现 handler。
5. 在 `biz.go / service.go / wire.go` 聚合 provider。
6. 运行生成与构建命令，确保可编译可运行。

## PR 检查清单（必须）

1. 是否做到同业务域三层同名目录同构。
2. SQL 是否归属到对应业务域目录。
3. 是否存在 `service -> data` 直接调用。
4. 是否存在 `biz` 依赖具体数据库客户端。
5. `cmd/server` 是否只保留启动装配职责。
6. 变更后是否通过生成、构建、运行验证。

## 常见反模式与修正

1. 反模式：按层平铺一个大文件承载多个域。
修正：按业务域拆目录，保持域边界独立。

2. 反模式：全局 `internal/data/sql` 混放所有业务 SQL。
修正：迁移到 `internal/data/<domain>/sql`。

3. 反模式：`service` 里做业务规则、`data` 里做业务判断。
修正：业务规则集中在 `biz/usecase`。

## 输出约定

使用本 Skill 处理"搭建/重构/评审"任务时，输出必须包含：
1. 分层边界结论（逐层通过/不通过）。
2. 域同构结论（目录是否同构，SQL 是否归域）。
3. 违规点清单（文件路径 + 修复动作）。
4. 最小改造顺序（按依赖安全顺序执行）。
