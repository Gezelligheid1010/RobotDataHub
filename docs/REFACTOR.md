# 业务模块化重构记录

本次从按技术层分包迁移至按业务域分包，保留全部 HTTP URL、数据库表名、请求/响应 JSON 字段与 Python Worker 协议。

- `server/internal/dataset`：数据集的 Handler、Service、Repository、模型、请求，以及对象存储文件业务。
- `server/internal/episode`：Episode 查询、Manifest 解析、带 Dataset 状态更新的事务导入。
- `server/internal/job`：任务增删查询、领取、完成和回收。
- `server/internal/platform`：共享的 PostgreSQL/MinIO 适配器和迁移初始化。
- `server/internal/transport/http`：共享响应格式和 ID 校验；API 路由及 Token 鉴权在 `server/cmd/server/router.go` 组装，避免循环依赖。
- `server/cmd/server`：依赖组装和进程启动，不承载业务逻辑。
- `web/src/modules/{dataset,episode,job}`：前端页面与对应 API 内聚。
- `web/src/shared/api/client.js`：Axios 请求基础设施。

跨域协作：Job Service 通过 Dataset Service 查询有效数据集；Episode 导入通过数据库事务同时写入 episodes 和 datasets。后者属于一个有意识的跨模块写事务，后续可进一步抽出 Import 用例协调器，不为当前 MVP 提前增加复杂度。

## 保留的 MVP 限制

- 任务领取依赖 Postgres `FOR UPDATE SKIP LOCKED`，尚不具备心跳续约/执行 token 校验。
- 文件上传仅做路径检查和大小检查，未提供完整恶意文件防护、对象存储与数据库的原子一致性保障。
- 通过 API Token 认证，不是 JWT/RBAC。
- LeRobot 的 Parquet/MP4 元数据和帧级可视化尚未完整实现。

## 构建

`cd server && go test ./... && go build ./cmd/server`

`cd web && npm install && npm run build`

`cd worker-python && python -m unittest discover -s tests`

`cd deploy && docker compose --env-file .env up --build -d`
