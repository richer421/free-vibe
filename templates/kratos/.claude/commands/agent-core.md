# Free Vibe Coding Agent Core (Proxy)

## Overview

这是仓库级 **Skill Proxy**，不是业务规则集合。
它的唯一职责是：先判定任务类型，再路由到正确的子 Skill，确保代理不会"随缘发挥"。

## 路由协议（必须执行）

1. 接到任务后，先做任务分类，不直接写代码。
2. 按路由表选择子 Skill。
3. 若命中"新模块开发"类型，先执行 Knowledge Gate（知识门禁）。
4. 明确声明路由结果（命中哪个 Skill，为什么）。
5. 子 Skill 执行完后，再由本 Proxy 汇总输出结果。

## 路由表

1. 新模块开发（新增业务域/新增业务模块）
- 触发条件：包含"新增模块/新增业务域/new module/new domain/从 0 到 1 做一个模块"等意图。
- 必选子 Skill（按顺序）：
  - `/business-domain`
  - `/kratos-layout`
- 路由动作：
  - 先执行 Knowledge Gate。
  - Gate 未通过时：只输出缺失项与填写模板，不进入代码实现。
  - Gate 通过后：再进入分层设计与实现。

2. Kratos 架构/分层/目录/模块/重构/评审/实现任务
- 触发条件：任务涉及 `api`、`biz`、`data`、`service`、`server`、`wire`、`sqlc` 任一层。
- 必选子 Skill（按顺序）：
  - `/business-domain`
  - `/kratos-layout`
- 路由动作：先完成业务语义与模型映射，再按业务域同构规则落地实现。

3. 纯仓库流程或通用协作问题（不涉及代码改动）
- 触发条件：仅讨论流程、策略、协作方式、产出结构。
- 必选子 Skill：无
- 路由动作：由 Proxy 直接处理，但仍遵守本仓库输出规范。

## 路由输出规范（必须）

在开始实施前，先输出一行路由声明：
- `ROUTE: <task-type> -> <skill-name | direct>`
若为新模块开发，再输出一行门禁状态：
- `KNOWLEDGE_GATE: pass`
- `KNOWLEDGE_GATE: fail (<missing-items>)`

示例：
- `ROUTE: new-module-development -> /business-domain + /kratos-layout`
- `KNOWLEDGE_GATE: fail (missing business domains, missing core models)`
- `ROUTE: kratos-module-refactor -> /business-domain + /kratos-layout`
- `ROUTE: process-discussion -> direct`

## 失败与降级策略

1. 若目标子 Skill 不存在或不可读：明确报错并给出最接近的替代执行规则，不得静默跳过。
2. 若任务跨多个类型：先执行"架构约束最强"的子 Skill，再补充其他要求。
3. 若用户要求与架构规则冲突：先指出冲突，再给两个可执行选项。

## 执行边界

1. 本 Proxy 不负责定义具体分层细则，细则由子 Skill 提供。
2. 本 Proxy 负责"先路由、再执行、后汇总"的顺序保障。
3. 本 Proxy 负责保持任务焦点，不被历史上下文带偏。
