# Codex 规则

## 适用范围
- 本规则适用于 `shop-admin` 仓库全量目录（`server` + `web`）。

## 文档同步
- 每次修改代码后，必须同步检查并更新 `README.md`，确保文档与实际代码行为一致。
- 提交推送前，必须先完成 `README.md` 更新。
- 提交代码时，必须将 `README.md` 改动与本次代码改动一起提交。

## 生成与构建门禁
- 本项目依赖多段生成链路，提交前需按改动范围执行：
  1. 修改 DI/ProviderSet/构造函数后，执行 `cd server && make wire`。
  2. 修改 `server/api/protos/**` 后，执行：
     - `cd server && make api`
     - `cd server && make openapi`
     - `cd server && make ts`
  3. 修改前端页面或静态资源后，执行：
     - `cd web && pnpm build-only`
- 前端构建产物必须输出到：
  - `server/cmd/server/assets/web`
- OpenAPI 产物输出到：
  - `server/cmd/server/assets`

## 提交流程约定
- 用户要求“提交”时，默认执行完整动作：`git commit` + `git push`。
- 未明确指定分支时，推送当前分支到同名远程分支。
- `git commit -m` 信息默认使用中文，简洁描述本次变更。
- 若用户未指定提交信息，按变更内容自动生成中文提交信息。

## 项目专属约束
- `wireinject` 仅允许保留在 `server/cmd/server/wire.go`。
- `server/internal/configs|data|sdk|service` 下的 ProviderSet 文件禁止添加 `//go:build wireinject`。
- 修改 `server/internal/data` 时，必须同时核对：
  1. 是否仍有业务引用；
  2. 是否仍在 `ProviderSet` 中被使用。

## 注释与风格规范
- 后续新增或修改代码时，代码注释统一使用中文。
- 禁止提交 IDE/系统垃圾文件（如 `.idea/`、`.DS_Store`）。

