# FreeVibe

FreeVibe 是一个后端基础模板仓库，并提供 `freevibe` CLI：
- 以 `git submodule` 方式组织大项目模块
- 以模板仓库作为脚手架来源，把代码初始化到目标模块仓库
- 仓库根目录负责 CLI 与发布维护，Kratos 后端代码模板位于 `templates/kratos/`

## 快速开始

### 1) 安装 CLI（最新版本）

```bash
curl -fsSL https://github.com/richer421/free-vibe/releases/latest/download/install.sh | bash
```

说明：再次执行同一条安装命令，就是升级（无需单独 update 命令）。

### 2) 验证安装

```bash
freevibe version
```

### 3) 初始化一个父项目（submodule 结构）

```bash
freevibe init my-monorepo \
  --backend-name order-service \
  --repo https://github.com/<owner>/order-service.git
```

说明：
- `--repo` 必填，表示目标模块仓库
- `--backend-name` 可选；不传时默认取仓库名
- 目标仓库必须已存在

初始化后会生成：
- `freevibe.modules.yaml`：模块注册表
- `.gitmodules`：子模块定义
- 根 `Makefile`：`modules/status/pull`

### 4) 验证初始化

```bash
cd my-monorepo
git submodule status
```

## 常用命令

```bash
# 查看版本
freevibe version

# 新增模块
freevibe add --repo https://github.com/<owner>/payment-service.git

# 指定模块名
freevibe add --name payment-service --repo https://github.com/<owner>/payment-service.git

# 移除模块
freevibe remove payment-service
```

## 指定版本安装

```bash
# 直接安装 v0.1.3
curl -fsSL https://github.com/richer421/free-vibe/releases/download/v0.1.3/install.sh | bash
```

可选参数（高级）：
```bash
# 指定安装目录
curl -fsSL https://github.com/richer421/free-vibe/releases/latest/download/install.sh | \
  bash -s -- --install-dir /usr/local/bin
```

## 本地开发

```bash
# 构建 CLI
make build-cli

# 运行帮助
./bin/freevibe --help
```

仓库结构：

- 根目录：`freevibe` CLI、发布脚本、仓库维护文件
- `templates/kratos/`：Kratos 后端服务模板，生成模块代码时直接使用这个目录
