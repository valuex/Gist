# Gist

[![License: GPL v2](https://img.shields.io/badge/License-GPL_v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html) [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/9bingyin/Gist) [![zread](https://img.shields.io/badge/Ask_Zread-_.svg?style=flat&color=00b0aa&labelColor=000000&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQuOTYxNTYgMS42MDAxSDIuMjQxNTZDMS44ODgxIDEuNjAwMSAxLjYwMTU2IDEuODg2NjQgMS42MDE1NiAyLjI0MDFWNC45NjAxQzEuNjAxNTYgNS4zMTM1NiAxLjg4ODEgNS42MDAxIDIuMjQxNTYgNS42MDAxSDQuOTYxNTZDNS4zMTUwMiA1LjYwMDEgNS42MDE1NiA1LjMxMzU2IDUuNjAxNTYgNC45NjAxVjIuMjQwMUM1LjYwMTU2IDEuODg2NjQgNS4zMTUwMiAxLjYwMDEgNC45NjE1NiAxLjYwMDFaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00Ljk2MTU2IDEwLjM5OTlIMi4yNDE1NkMxLjg4ODEgMTAuMzk5OSAxLjYwMTU2IDEwLjY4NjQgMS42MDE1NiAxMS4wMzk5VjEzLjc1OTlDMS42MDE1NiAxNC4xMTM0IDEuODg4MSAxNC4zOTk5IDIuMjQxNTYgMTQuMzk5OUg0Ljk2MTU2QzUuMzE1MDIgMTQuMzk5OSA1LjYwMTU2IDE0LjExMzQgNS42MDE1NiAxMy43NTk5VjExLjAzOTlDNS42MDE1NiAxMC42ODY0IDUuMzE1MDIgMTAuMzk5OSA0Ljk2MTU2IDEwLjM5OTlaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik0xMy43NTg0IDEuNjAwMUgxMS4wMzg0QzEwLjY4NSAxLjYwMDEgMTAuMzk4NCAxLjg4NjY0IDEwLjM5ODQgMi4yNDAxVjQuOTYwMUMxMC4zOTg0IDUuMzEzNTYgMTAuNjg1IDUuNjAwMSAxMS4wMzg0IDUuNjAwMUgxMy43NTg0QzE0LjExMTkgNS42MDAxIDE0LjM5ODQgNS4zMTM1NiAxNC4zOTg0IDQuOTYwMVYyLjI0MDFDMTQuMzk4NCAxLjg4NjY0IDE0LjExMTkgMS42MDAxIDEzLjc1ODQgMS42MDAxWiIgZmlsbD0iI2ZmZiIvPgo8cGF0aCBkPSJNNCAxMkwxMiA0TDQgMTJaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00IDEyTDEyIDQiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8L3N2Zz4K&logoColor=ffffff)](https://zread.ai/9bingyin/Gist)

[![GitHub Release](https://img.shields.io/github/v/release/9bingyin/Gist)](https://github.com/9bingyin/Gist/releases/latest) [![Build Docker Image](https://github.com/9bingyin/Gist/actions/workflows/docker-build.yml/badge.svg)](https://github.com/9bingyin/Gist/actions/workflows/docker-build.yml)

轻量级自托管 RSS 阅读器，内置 AI 能力。

![desktop](docs/images/desktop.png)

移动端截图  
<img width="1938" height="1422" alt="image" src="https://github.com/user-attachments/assets/431221a6-0cc5-40d8-8428-b3e946f7ccc5" />

## 功能特性

- 全格式订阅，支持 RSS 2.0 / Atom / JSON Feed
- Readability 沉浸式阅读模式
- AI 摘要与翻译，支持 OpenAI / Anthropic / 兼容接口 (BYOK)
- 文件夹分层管理与内容分类
- 浅色 / 深色 / 跟随系统主题
- PWA，可安装到桌面和移动设备，滚动时可触发终端浏览器地址栏和工具栏隐藏，实现阅读界面最大化
- 多语言 (简体中文 / English)

## 部署

### Docker Compose (推荐)

```bash
curl -O https://raw.githubusercontent.com/9bingyin/Gist/main/docker-compose.yml
docker compose up -d
```

或手动创建 `docker-compose.yml`:

```yaml
services:
  gist:
    image: ghcr.io/9bingyin/gist:latest
    container_name: gist
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - GIST_LOG_LEVEL=info
    restart: always
```

访问 `http://localhost:8080`，数据持久化在 `./data` 目录。

###  通过镜像包安装
参照重新部署说明：[link(https://github.com/valuex/Gist/blob/main/reinstall_using_tar.md)

### Docker Run

```bash
docker run -d \
  --name gist \
  -p 8080:8080 \
  -v ./data:/app/data \
  ghcr.io/9bingyin/gist:latest
```

### 镜像标签

| 标签 | 说明 |
|------|------|
| `latest` | 最新稳定版 |
| `1.2.0` | 指定版本 |
| `1.2` | 该 minor 版本的最新 patch |
| `1` | 该 major 版本的最新 minor |
| `develop` | 每次推送 `main` 分支自动构建 |

所有镜像均为多架构 (`linux/amd64`, `linux/arm64`)。

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `GIST_ADDR` | `:8080` | 监听地址 |
| `GIST_DATA_DIR` | `/app/data` | 数据目录 |
| `GIST_STATIC_DIR` | `/app/static` | 静态文件目录 |
| `GIST_LOG_LEVEL` | `info` | 日志级别 (`debug` / `info` / `warn` / `error`) |

## 本地开发

### 前置依赖

- Go 1.25+
- [Bun](https://bun.sh/)

### 后端

```bash
cd backend
go mod download
go run ./cmd/server/main.go
```

### 前端

```bash
cd frontend
bun install
bun run dev
```

### 测试

```bash
# 后端
cd backend
make test    # 运行测试 (含 race 检测)
make lint    # 运行 golangci-lint

# 前端
cd frontend
bun run test
bun run lint
```

## 许可证

[GPL-2.0](./LICENSE)

## 订阅 API（NAS / 自动化场景）

支持通过 HTTP API 从容器外部添加 RSS 订阅，适用于 NAS 上脚本、快捷指令等自动化场景。

所有接口需携带 `Authorization: ****** 请求头（token 通过登录接口获取）。

### 1. 新增单条订阅

```bash
# 获取 token
TOKEN=$(curl -s -X POST "http://NAS_IP:PORT/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"identifier":"your_username","password":"your_password"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 新增订阅（201 = 新建，200 = 已存在）
curl -X POST "http://NAS_IP:PORT/api/subscriptions" \
  -H "Authorization: ******" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/rss.xml","title":"可选标题","category":"可选分类"}'
```

**请求体**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `url` | string | ✅ | RSS/Atom feed 地址（http/https） |
| `title` | string | | 自定义标题（不填则使用 feed 原标题） |
| `category` | string | | 分类（文件夹）名称，不存在则自动创建 |

**响应示例（201 新建）**

```json
{
  "id": "1234567890",
  "title": "Example Feed",
  "url": "https://example.com/rss.xml",
  "folderId": "9876543210",
  "isNew": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

`isNew: false` 表示订阅已存在（返回 HTTP 200，幂等）。

### 2. 批量新增订阅

```bash
curl -X POST "http://NAS_IP:PORT/api/subscriptions/batch" \
  -H "Authorization: ******" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"url":"https://a.example.com/rss","category":"Tech"},
      {"url":"https://b.example.com/feed.xml","title":"Blog B"}
    ]
  }'
```

**响应示例**

```json
{
  "created": 1,
  "skipped": 1,
  "errors": 0,
  "results": [
    {"url":"https://a.example.com/rss","status":"created","feed":{...}},
    {"url":"https://b.example.com/feed.xml","status":"exists","feed":{...}}
  ]
}
```

`status` 取值：`"created"`（新建）、`"exists"`（已存在）、`"error"`（失败）。

> **注意**：批量接口不拉取 feed 内容（速度更快），单条接口会拉取 feed 以获取初始条目。

## 便捷操作特性
- 特定feed支持三种视图自定义：常规模式（显示rss feed提供的内容）；阅读模式（显示全文）；浏览器（将文章在新tab打开，适用于需要登录查看全文的网站）
- 文章列表右滑显示feed列表
- 自动触发隐藏终端浏览器的地址栏和工具栏，实现阅读界面最大化
- 滚动出顶部工具栏时自动将文章标记为已读
- 浮动按钮便于将整个目录或feed下文章标记为已读
- 记住feed文章列表和文章阅读位置，便于继续从上次位置开始阅读

