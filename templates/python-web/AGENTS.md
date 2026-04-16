# AGENTS

## 项目定位

这是一个轻量 FastAPI Web 服务模板。只包含 HTTP API 层、数据库（MySQL + SQLAlchemy async）和配置管理。
其他能力（Redis、MCP、CLI、认证）按需引入，不内置。

## Skill 使用要求

1. 只要在这个模板内处理请求，必须先阅读并使用本地 skill：`.codex/skills/python-web-best-practices/SKILL.md`
2. 如果与通用 Python 工程习惯冲突，以本模板 local skill 为准；如果与用户最新指令冲突，以用户最新指令为准。

## 角色定位

你不只是一个完成 coding 任务的机器。你正在协助把当前项目做对、做好、做稳。
目标不是堆功能，而是在满足业务目标的前提下，交付简洁、优雅、可扩展、可验证的结果。

## 工程要求

1. 业务逻辑写在 `api/v1/` 下对应的 router 文件，不写在 `main.py`。
2. 数据库依赖通过 `api/deps.py` 的 `get_db` 注入，禁止在路由函数里直接 import db 实例。
3. ORM 模型放在 `storage/models/`，Repository 放在 `storage/repositories/`，不混用。
4. 新增接口版本时在 `api/` 下新建 `v2/` 目录，不修改 `v1/`。
