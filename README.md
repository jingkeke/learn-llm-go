# Go Auth + SQLite + Node.js + Next.js MVP Demo

这个项目展示了如何使用 `github.com/smhanov/auth` 和 SQLite 作为核心权限服务，并使用 Node.js 作为部分后端业务逻辑服务，最后通过 Next.js 前端进行整合。

## 架构说明

本项目由三个独立的服务组成：

1. **Go Auth Service (端口 8080):**
   - 核心认证服务。
   - 负责处理用户注册 (`/user/create`)、登录 (`/user/auth`)、注销 (`/user/signout`) 和获取用户信息 (`/user/get`)。
   - 使用 SQLite 作为数据库 (开启 WAL 模式以支持高并发)。
   - 签发 session cookie，设置 `HttpOnly` 以增强安全性。
   - 处理跨域请求 (CORS)，允许来自前端和 Node.js 后端的请求携带凭证 (`Credentials`)。

2. **Node.js Backend Service (端口 4000):**
   - 业务逻辑后端服务。
   - 包含受保护的 API `/api/protected-data`。
   - 拦截前端请求，读取 cookie 中的 `session` 字段，并携带该 cookie 请求 Go Auth Service 的 `/user/get` 接口验证权限。
   - 如果 Go 服务验证通过，Node.js 允许访问资源；否则返回 401 未授权。

3. **Next.js Frontend (端口 3000):**
   - 用户界面。
   - 提供注册和登录表单，直接请求 Go Auth Service。
   - 登录成功后，浏览器会自动保存 Go Auth Service 签发的 `session` cookie。
   - 之后请求 Node.js 保护资源时，浏览器自动将 cookie 发送给 Node.js 后端。
   - 页面状态会根据 `goAuthApi.get('/get')` 的结果自动更新。

## 如何运行项目

### 1. 启动 Go Auth Service
```bash
cd go-auth-service
go run main.go
```
服务将在 `http://localhost:8080` 启动。

### 2. 启动 Node.js Backend Service
```bash
cd node-backend
npm install
node server.js
```
服务将在 `http://localhost:4000` 启动。

### 3. 启动 Next.js Frontend
```bash
cd next-frontend
npm install
npm run dev
```
服务将在 `http://localhost:3000` 启动。

在浏览器中打开 `http://localhost:3000`，你可以测试注册、登录以及获取 Node.js 受保护数据的完整流程。

---

## 架构对比分析：Go Auth (SQLite) vs Supabase

根据作者 `smhanov` 的文章，使用 SQLite + WAL 模式不仅性能卓越且支持高并发，结合他编写的 `smhanov/auth` 库，可以快速实现完整的认证逻辑，避免复杂的外部依赖（如远程 Postgres 等）。

下面我们与目前流行的 Backend-as-a-Service (BaaS) 平台 **Supabase** 进行全方位的对比：

### 1. 简单性与上手难度
* **Go Auth (SQLite) + 自建 Node.js**
  * **简单点：** 没有外部依赖，SQLite 只是一个本地文件 (`.db`)，零配置。Go 的 `auth` 库提供了非常"Boring"但是开箱即用的路由和表结构自动创建。所有的权限逻辑在你的掌握之中。
  * **麻烦点：** 这种架构是"拼装式"的（如本 Demo）。你仍然需要自己编写微服务间的认证校验逻辑（比如 Node.js 去请求 Go 服务验签），还需要自己处理服务器部署、CORS、SSL 证书等。
* **Supabase**
  * **简单点：** 提供开箱即用的身份验证 (Auth)、数据库 (Postgres)、实时订阅和存储。前端直接集成 `@supabase/supabase-js`，无需自己写后端中间件。注册、登录、OAuth 以及鉴权都在控制台一键配置。
  * **麻烦点：** 它是基于 Postgres Row Level Security (RLS) 的，学习 RLS 语法和权限控制策略有一定门槛。当业务逻辑变得复杂且无法仅用 SQL/RLS 表示时，需要写 Edge Functions 或者自己接后端。

### 2. 性能与并发
* **Go Auth (SQLite)**
  * **优势：** SQLite + WAL + C/GO 接口访问速度极快。没有网络跳数 (TCP Network Hop) 延迟。在单机 NVMe 驱动器上，轻松应对数千并发。对于 99% 的新项目和中小项目，单机性能完全过剩。
* **Supabase**
  * **优势：** 基于强大的 PostgreSQL，架构可扩展，适合真正需要分布式、海量数据的企业级应用。
  * **劣势：** 默认的网络连接需要 TCP 握手，相对单机 SQLite 有网络延迟。免费层或低配实例在并发极高时需要升级配置。

### 3. 部署与运维成本
* **Go Auth (SQLite)**
  * **极低成本：** 你只需要一台最便宜的 VPS（例如 5刀/月的机器）将 Go 二进制文件、Node.js 脚本和 SQLite `.db` 文件放上去即可运行。
  * **运维：** 备份极其简单，复制一个文件 (`users.db`) 就是全量备份。但如果要做到高可用（横向扩展多台机器），SQLite 共享会比较困难（需要如 LiteFS 之类的技术方案）。
* **Supabase**
  * **成本：** 如果使用托管服务，前期有免费层，但当流量、数据库大小超标后，费用会显著增加。如果选择自托管 (Self-hosting)，部署一整套 Supabase (涉及十几个 Docker 容器：Kong, GoTrue, PostgREST, Realtime 等) 极其繁琐，运维成本极高。

### 4. 灵活性与控制力
* **Go Auth (SQLite)**
  * **高控制力：** 源代码在手，认证逻辑想改就改（例如定制响应结构、挂载回调 `OnAuthEvent`）。你的数据、你的数据库、你的网络栈完全由你控制，真正的"数据主权"。
* **Supabase**
  * **强生态约束：** 强依赖于 Supabase 的生态和 Postgres 的 RLS。某些特定的企业级需求或者深度的认证流定制可能受限于 Supabase 提供的现有功能。

### 总结
* 如果你是一名喜欢 **"Keep It Simple, Stupid (KISS)"** 原则的开发者，喜欢拥有代码和数据的绝对控制权，目标是以极低的成本快速启动一个单体/小型微服务项目，并且你能熟练使用 Go/Node 处理前后端粘合工作，那么 **Go Auth + SQLite 方案非常棒，它极其简单、极快且廉价**。
* 如果你是一名前端开发者，不想写任何后端代码，希望快速拥有"认证 + 数据库 + 接口"的全套解决方案，或者预见项目未来一定会迅速成长为需要复杂关系和实时能力的大型平台，那么 **Supabase 是更省时省力的选择**。
