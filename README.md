# FreeVibe Backend Template

This repository is the backend base template and now also contains the `freevibe` CLI.

## FreeVibe CLI

Build CLI:
```bash
make build-cli
./bin/freevibe --help
```

Initialize a parent project with submodule layout (pulls module from remote and renders it):
```bash
./bin/freevibe init my-monorepo \
  --backend-name order-service \
  --template-repo-url https://github.com/richer421/free-vibe.git
```

Add/remove modules in an existing parent project:
```bash
./bin/freevibe add --name payment-service --type backend
./bin/freevibe remove payment-service
```

Generated parent project includes:
- `freevibe.modules.yaml`: module registry
- `Makefile`: `modules/status/pull`
- Git submodules for each module

## Kratos Base Usage

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## SQLC + MySQL（已接入默认 Greeter）
```bash
# 1) 准备 MySQL 表（示例 schema）
# internal/data/greeter/sql/schema/001_greeters.sql

# 2) 生成 sqlc 代码
make sqlc

# 3) 生成 sqlc + wire 并构建
make generate
make build

# 4) 运行服务（确保 configs/config.yaml 中数据库可连通）
./bin/free-vibe-coding -conf ./configs/config.yaml
```

关键目录：
- `sqlc.yaml`
- `internal/data/greeter/sql/schema`
- `internal/data/greeter/sql/query`
- `internal/data/gen/sqlc` (generated)
- `internal/data/greeter/repo.go`（仓储适配，调用 sqlc）

## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```
