---
name: python-web-best-practices
description: Use when adding API endpoints, storage models, or working in the FastAPI adapter layer of this python-web template.
---

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

新增模型 + 仓储示例：
```python
# storage/models/order.py
from sqlalchemy.orm import Mapped, mapped_column
from .base import Base

class Order(Base):
    __tablename__ = "orders"
    id: Mapped[int] = mapped_column(primary_key=True)
    product_id: Mapped[str]
    quantity: Mapped[int]
    status: Mapped[str] = mapped_column(default="pending")

# storage/repositories/order.py
from .base import BaseRepository
from ..models.order import Order

class OrderRepository(BaseRepository[Order]):
    pass
```

## 依赖注入规范

- 数据库 session 必须通过 `Depends(get_db)` 注入。
- 如需 Redis 等其他依赖，在 `api/deps.py` 中统一定义注入函数。

## 反模式（禁止）

1. 在 `main.py` 里写业务逻辑。
2. 路由函数直接 import 存储实例。
3. Repository 里写业务规则。
4. 修改 `v1/` 已有接口结构（应新建 `v2/`）。
