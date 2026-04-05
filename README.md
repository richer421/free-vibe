# FreeVibe

FreeVibe 是一个后端基础模板仓库，并提供 `freevibe` CLI：
- 以 `git submodule` 方式组织大项目模块
- 从远程仓库拉取模板并进行基础渲染

## 快速开始

### 1) 安装 CLI（最新版本）

```bash
curl -fsSL https://raw.githubusercontent.com/richer421/free-vibe/main/scripts/install.sh | bash
```

### 2) 更新 CLI（升级到最新版本）

```bash
curl -fsSL https://raw.githubusercontent.com/richer421/free-vibe/main/scripts/install.sh | bash -s -- --update
```

### 3) 初始化一个父项目（submodule 结构）

```bash
~/.local/bin/freevibe init my-monorepo --backend-name order-service
```

初始化后会生成：
- `freevibe.modules.yaml`：模块注册表
- `.gitmodules`：子模块定义
- 根 `Makefile`：`modules/status/pull`

## 常用命令

```bash
# 查看版本
freevibe version

# 新增模块
freevibe add --name payment-service --type backend

# 移除模块
freevibe remove payment-service
```

## 安装脚本参数

```bash
# 指定版本安装
curl -fsSL https://raw.githubusercontent.com/richer421/free-vibe/main/scripts/install.sh | \
  bash -s -- --version v0.1.0

# 指定安装目录
curl -fsSL https://raw.githubusercontent.com/richer421/free-vibe/main/scripts/install.sh | \
  bash -s -- --install-dir /usr/local/bin
```

## Release 版本管理

本仓库采用 `SemVer`（`vX.Y.Z`）+ GitHub Release：

1. 本地提交完成后打 tag：
```bash
./scripts/release.sh v0.1.0
```

2. 推送 tag 后，GitHub Actions 自动执行：
- 构建 `linux/darwin` + `amd64/arm64` 二进制
- 打包为 `freevibe_<version>_<os>_<arch>.tar.gz`
- 发布到 GitHub Release

工作流文件：
- `.github/workflows/release.yml`
- `.goreleaser.yaml`

## 本地开发

```bash
# 构建 CLI
make build-cli

# 运行帮助
./bin/freevibe --help
```
