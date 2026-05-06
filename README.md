# Campaign Platform

活动页快速搭建平台。Go 后端 + Next.js 前端 + PostgreSQL。

## 快速开始

```bash
# 1. 建库
createdb campaign_platform
psql -d campaign_platform -f migrations/001_init.sql

# 2. 启动 API
cd backend && go run ./cmd/server

# 3. 启动前端
cd web && npm install && npm run dev

# 4. 创建并上线活动
curl -X POST http://localhost:8080/api/v1/campaigns \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test","slug":"test","config":{"lang":"en","sections":[{"component":"hero_banner","props":{"title":{"en":"Hello"}}}]}}'

curl -X PATCH http://localhost:8080/api/v1/campaigns/1/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"active"}'
```

## 项目结构

```
backend/          — Go API（9 个接口）
  cmd/server/     入口
  internal/
    handler/      HTTP 层
    service/      业务逻辑
    repository/   数据库
    model/        数据模型
web/              — Next.js 渲染
  app/            页面路由
  components/     活动组件
  lib/            工具库
migrations/       — SQL 迁移
deploy/           — systemd + nginx + 部署脚本
.github/workflows — CI/CD
```

## 部署

```bash
git push main        # GitHub Actions 自动部署
# 或手动：
./deploy/deploy.sh user@server
```
