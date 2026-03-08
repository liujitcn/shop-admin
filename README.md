# shop-admin

商城管理端项目，包含：

- `server`：Go 后端（Kratos + Wire + Buf）
- `web`：Vue3 前端管理台（Vite + Element Plus）

## 技术栈

- 后端：Go `1.26`、Kratos、Wire、Buf、GORM Gen、go-sdk（auth/cache/queue/gorm 等）
- 前端：Vue3、TypeScript、Vite、Pinia、Element Plus、UnoCSS、pnpm
- API：Proto3 + Buf 生成 Go/HTTP/OpenAPI/TS 代码

## 项目结构

```text
shop-admin/
├── server/                        # 后端
│   ├── cmd/server/                # 启动入口、wire、静态资源嵌入
│   │   └── assets/
│   │       ├── assets.go          # embed openapi + web
│   │       └── web/               # 前端打包产物（嵌入后端）
│   ├── internal/
│   │   ├── configs/               # 配置解析与 ProviderSet
│   │   ├── const/                 # 常量定义
│   │   ├── core/                  # 核心运行时（ShopCore）
│   │   ├── middleware/            # 中间件实现（含 logging）
│   │   ├── sdk/                   # SDK ProviderSet（auth/cache/gorm/queue）
│   │   ├── data/                  # Repo 层
│   │   ├── service/               # Biz + gRPC/HTTP Service
│   │   └── server/                # HTTP Server 注册
│   ├── api/                       # proto、buf 生成规则
│   ├── configs/                   # 运行配置（data/logger/server/auth...）
│   └── Makefile
├── web/                           # 前端工程
└── README.md
```

## 环境要求

- Go `>= 1.26`
- Node.js `>= 18`
- pnpm（前端包管理）
- 常用生成工具（通过 `make init` 安装）
- `protoc-gen-ts_proto`（执行 `make ts` 需要，需额外安装并加入 PATH）

## 快速开始

### 1) 前端开发

```bash
cd web
pnpm install
pnpm dev
```

### 2) 后端开发

```bash
cd server
go mod tidy
make wire
go run ./cmd/server -conf ./configs
```

启动成功后，管理端默认访问地址为：`http://localhost:8091/`  

### 3) 前端打包并嵌入后端

当前已配置为：前端构建产物直接输出到后端资源目录  
`/server/cmd/server/assets/web`

```bash
cd web
pnpm build-only
```

构建成功后，后端通过 `embed` 自动加载该目录资源（见 `server/cmd/server/assets/assets.go`）。
为避免以下划线开头的静态资源（如 `_initCloneObject.*.js`）在运行时 404，嵌入规则使用 `all:web/*`。

## 常用命令

### server/Makefile

- `make init`：安装 protoc / kratos / buf / lint 等开发工具
- `make api`：生成 Go 的 proto/grpc/http/error 代码（`server/api/gen/go`）
- `make openapi`：生成 OpenAPI 文档（输出到 `server/cmd/server/assets`）
- `make ts`：生成 TS proto（输出到 `web/src/rpc`）
- `make wire`：生成 `cmd/server/wire_gen.go`
- `make test`：运行单测
- `make compose-up`：启动依赖服务（docker compose）
- `make compose-down`：停止依赖服务
- `make run`：生成 API + OpenAPI 后启动服务

### web/package.json

- `pnpm dev`：本地开发
- `pnpm build`：类型检查 + 打包
- `pnpm build-only`：仅打包（产物输出到 `server/cmd/server/assets/web`）

## 代码生成与改动联动

以下改动需要联动执行对应命令：

1. 修改 `server/api/protos/**`：
   - 运行 `make api && make openapi && make ts`
   - 确认 `server/api/gen/go` 与 `web/src/rpc` 同步
2. 修改依赖注入构造函数（`NewXxx`）或 ProviderSet：
   - 运行 `make wire`
3. 前端路由/静态页面需要由后端统一承载：
   - 运行 `pnpm build-only`，确认产物在 `server/cmd/server/assets/web`
4. 修改 `server/internal/service/**` 下的服务构造函数：
   - 统一移除 `ctx *bootstrap.Context` 入参，改为注入 `sc *core.ShopCore`
   - 服务结构体统一嵌入 `*core.ShopCore`

## Codex Rules（继承 go-sdk + 本项目追加）

以下规则用于本仓库协作开发，默认必须遵守。

### A. 继承自 go-sdk 的基础规则

1. 代码有改动时，必须同步检查并更新 `README.md`（文档与代码同提交）。
2. Git 提交信息使用中文。
3. 新增/修改注释优先使用中文，避免无意义注释。
4. 禁止提交 IDE/系统垃圾文件（如 `.idea/`、`.DS_Store`）。

### B. shop-admin 追加规则

1. `wireinject` 仅用于 `server/cmd/server/wire.go`。  
   `internal/configs|data|sdk|service` 下的 ProviderSet 文件不要再加 `//go:build wireinject`。
2. 修改 `internal/data` 仓储代码时，必须同步校验：
   - 是否仍被业务引用
   - 是否仍在 `ProviderSet` 中使用
3. 前端打包输出目录固定为：
   - `server/cmd/server/assets/web`
   - 不要改回默认 `web/dist`
4. 变更 Proto 或 OpenAPI 时，不手改生成文件逻辑代码；优先改 `proto` 与生成配置。
5. 提交前最小校验建议：
   - 后端：`make wire`
   - 前端：`pnpm build-only`
   - 如涉及接口：`make api && make openapi && make ts`

### C. 提交规范建议

推荐单次提交遵循：

1. 代码变更
2. 生成物更新（wire/api/openapi/ts）
3. README 同步
4. 中文提交信息

## 备注

- `server/cmd/server/assets/web` 目录内容为构建产物，通常体积较大；清理或重建时注意不要影响后端启动。
- 若 `make wire` 报 provider 缺失，优先检查 `internal/sdk`、`internal/configs` 的构造函数签名与 `ProviderSet` 是否闭环。
