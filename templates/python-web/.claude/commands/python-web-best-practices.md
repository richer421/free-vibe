# Python Web Best Practices

## 新增接口流程（必须按顺序）

1. 检查 `knowledge/` 文档是否完备（业务定义 + 核心模型），不完备先补充。
2. 在 `src/__MODULE_NAME__/api/v1/` 下新建或编辑对应业务域的文件。
3. 定义请求/响应 Pydantic 模型（放在同文件顶部，或独立 `schemas/` 目录）。
4. 在 `api/router.py` 里 `include_router` 新路由。
5. 编写对应测试。

示例：
```python
# src/__MODULE_NAME__/api/v1/order.py
from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from ..deps import get_db

router = APIRouter()

class CreateOrderRequest(BaseModel):
    product_id: str
    quantity: int

class OrderResponse(BaseModel):
    id: int
    status: str

@router.post("/orders", response_model=OrderResponse)
async def create_order(
    body: CreateOrderRequest,
    db: AsyncSession = Depends(get_db),
) -> OrderResponse:
    ...
```

然后在 `router.py` 加：
```python
from .v1 import order
router.include_router(order.router, tags=["orders"])
```

## 存储层规范

1. ORM 模型放在 `storage/models/` 下，每个业务域一个文件，继承 `Base`。
2. Repository 放在 `storage/repositories/` 下，继承 `BaseRepository[T]`。
3. Repository 只封装数据库操作，不写业务规则。
4. 业务规则写在路由函数或单独的 service 层，不写在 Repository 里。

## 依赖注入规范

- 数据库 session 必须通过 `Depends(get_db)` 注入，禁止在路由函数里直接 import `_engine` 或 `_session_factory`。
- 如需 Redis 等其他依赖，在 `api/deps.py` 中统一定义注入函数。

## 反模式（禁止）

1. 在 `main.py` 里写业务逻辑。
2. 路由函数直接 `from __MODULE_NAME__.storage.db import _session_factory`。
3. Repository 里写价格计算、状态校验等业务规则。
4. 修改 `v1/` 已有接口的响应结构（应新建 `v2/`）。
