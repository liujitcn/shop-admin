# shop-admin

商城管理端项目，包含 Go 后端与 Vue3 前端：

- `server`：Kratos + Wire + Buf + GORM
- `web`：Vite + Vue3 + Element Plus

## 目录结构

```text
shop-admin/
├── server/
│   ├── api/                      # proto 与 buf 生成配置
│   ├── cmd/server/               # 后端启动入口（含 wire、静态资源嵌入）
│   │   └── assets/
│   │       ├── openapi.yaml
│   │       └── web/              # 前端打包产物（由 web 构建输出）
│   ├── configs/                  # 后端配置
│   ├── internal/                 # 业务代码
│   └── Makefile
├── web/                          # 前端工程
└── README.md
```

## 环境要求

- Go `>= 1.26`（见 `server/go.mod`）
- Node.js `>= 18`
- pnpm
- MySQL、Redis 等依赖（请按本地环境自行启动）

## 快速开始

### 1. 启动后端

```bash
cd server
go mod tidy
make wire
go run ./cmd/server -conf ./configs
```

默认端口（见 `server/configs/server.yaml`）：

- HTTP: `8091`
- gRPC: `6001`

访问地址：

- 管理页面：`http://localhost:8091/`（或 `http://localhost:8091/web/`）
- Swagger 文档：`http://localhost:8091/docs`

### 2. 启动前端开发模式

```bash
cd web
pnpm install
pnpm dev
```

## 前端打包并嵌入后端

当前前端构建输出目录已固定为：

- `server/cmd/server/assets/web`

执行：

```bash
cd web
pnpm build-only
```

后端通过 `embed` 提供静态资源（`server/cmd/server/assets/assets.go`）：

- `//go:embed all:web/*`

说明：使用 `all:web/*` 是为了包含以下划线 `_` 开头的资源文件，避免运行时 404。

## 常用命令

### server

```bash
cd server
make init        # 安装开发工具
make api         # 生成 Go proto/grpc/http/error
make openapi     # 生成 OpenAPI 文档
make ts          # 生成 TS proto
make wire        # 生成 wire 注入代码
make run         # 生成 API+OpenAPI 后启动服务
go test ./...    # 运行测试
```

### web

```bash
cd web
pnpm dev
pnpm build
pnpm build-only
```

## 开发联动建议

1. 修改 `server/api/protos/**` 后执行：`make api && make openapi && make ts`
2. 修改依赖注入构造函数或 ProviderSet 后执行：`make wire`
3. 前端页面需要由后端托管时执行：`pnpm build-only`
4. 提交前建议至少执行：`cd server && go test ./...`
