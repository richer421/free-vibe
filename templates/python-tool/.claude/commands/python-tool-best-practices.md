# Python Tool Best Practices

## 新增工具流程（必须按顺序）

1. 检查 `knowledge/` 文档是否完备（业务定义 + 核心模型），不完备先补充。
2. 在 `src/__MODULE_NAME__/tools/` 下新建一个文件（每个业务域一个文件）。
3. 定义输入/输出 Pydantic 模型（放在同文件顶部或 `schemas/` 中）。
4. 用 `@tool` 装饰器注册工具函数。
5. 运行测试，确认三种调用方式均可用。
6. 不需要修改 `http/`、`mcp/`、`cli/` 任何文件。

示例：
```python
# src/__MODULE_NAME__/tools/price.py
from __MODULE_NAME__.tools.registry import tool, ToolContext
from pydantic import BaseModel

class CalculatePriceInput(BaseModel):
    product_id: str
    quantity: int

class CalculatePriceOutput(BaseModel):
    total: float
    currency: str = "CNY"

@tool(description="计算商品价格")
async def calculate_price(params: CalculatePriceInput, ctx: ToolContext) -> CalculatePriceOutput:
    return CalculatePriceOutput(total=params.quantity * 10.0)
```

## 工具函数签名规范

1. 必须是 `async def`。
2. 第一个参数：单一 Pydantic `BaseModel`，命名为 `params`。
3. 第二个参数：`ctx: ToolContext`，由适配层注入。
4. 返回值：单一 Pydantic `BaseModel`。
5. 禁止工具函数直接 import 存储实例，必须通过 `ctx`。

## 存储层规范

1. ORM 模型放在 `storage/models/` 下，每个业务域一个文件。
2. Repository 放在 `storage/repositories/` 下，继承 `BaseRepository[T]`。
3. Repository 只封装数据库操作，不写业务规则。
4. 业务规则写在工具函数里，不写在 Repository 里。

## 适配层规范

`http/`、`mcp/`、`cli/` 三个适配层自动从 registry 读取所有工具，不需要手动注册。
禁止在适配层写任何业务逻辑。

## 反模式（禁止）

1. 工具函数直接 import 存储实例（必须通过 ctx）。
2. 在适配层里写业务逻辑。
3. 一个工具文件里混放多个不相关业务域的工具。
4. 新增工具时修改适配层代码。
5. Repository 里写业务规则。
