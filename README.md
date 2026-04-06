# FreeVibe

FreeVibe 是一个多模板脚手架仓库，并提供 `freevibe` CLI：
- 以 `git submodule` 方式组织大项目模块
- 以模板仓库作为脚手架来源，把代码初始化到目标模块仓库
- 仓库根目录负责 CLI 与发布维护，模板目录位于 `templates/`

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

### 3) 查看可用模板

```bash
freevibe template ls
```

当前内置模板：

- `kratos`：Kratos 后端服务模板
- `console-react`：中后台 React 前端模板，技术栈为 React 18.3、TypeScript 5.6、Ant Design 5.24、Tailwind CSS 3.4、Vite 5.4

### 4) 初始化一个父项目（submodule 结构）

```bash
freevibe init my-monorepo \
  --template kratos \
  --name order-service \
  --repo https://github.com/<owner>/order-service.git
```

说明：
- `--template` 必填，表示要使用的项目模板
- `--repo` 必填，表示目标模块仓库
- `--name` 可选；不传时默认取仓库名
- 目标仓库必须已存在

初始化后会生成：
- `freevibe.modules.yaml`：模块注册表
- `.gitmodules`：子模块定义
- 根 `Makefile`：`modules/status/pull`

### 5) 验证初始化

```bash
cd my-monorepo
git submodule status
```

## 常用命令

```bash
# 查看版本
freevibe version

# 查看模板
freevibe template ls

# 新增 React 前端模块
freevibe add --template console-react --repo https://github.com/<owner>/console-web.git

# 新增模块
freevibe add --template kratos --repo https://github.com/<owner>/payment-service.git

# 指定模块名
freevibe add --template kratos --name payment-service --repo https://github.com/<owner>/payment-service.git

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
- `templates/console-react/`：中后台 React 前端模板，生成模块代码时直接使用这个目录
