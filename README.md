# 夸克网盘 MCP Server

基于 Go 实现的夸克网盘 [Model Context Protocol](https://modelcontextprotocol.io) (MCP) 服务器，通过 stdio 提供文件管理、搜索、上传下载和分享转存能力。

仓库：[`https://github.com/chiehw/quark-mcp`](https://github.com/chiehw/quark-mcp)

可执行文件名为 `quark-mcp`（Windows 为 `quark-mcp.exe`）。

## 功能特性

- **列出文件/文件夹** - 列出目录内容
- **创建文件夹** - 创建新文件夹
- **重命名** - 重命名文件和文件夹
- **删除** - 删除文件和文件夹
- **移动** - 移动文件和文件夹
- **复制** - 复制文件和文件夹
- **获取信息** - 获取文件或文件夹的详细信息
- **下载文件** - 异步下载，支持并发连接（3线程，10MB分片）
- **上传文件** - 上传本地文件到夸克网盘
- **正则重命名** - 使用正则表达式批量重命名文件
- **搜索** - 根据关键词递归搜索文件
- **获取分享文件列表** - 获取分享链接根目录的文件列表
- **获取分享详情** - 获取分享链接中指定文件夹的文件列表
- **保存分享文件** - 将分享链接中的文件保存到自己的网盘

## 快速开始

1. 从 [GitHub Releases](https://github.com/chiehw/quark-mcp/releases) 下载对应系统和 CPU 架构的安装包。
2. 准备夸克网盘 Cookie（见下方「如何获取 Cookie」）。
3. 通过 `QUARK_COOKIE` 启动服务器，或在支持 MCPB 的桌面客户端中安装 `.mcpb` 并在界面中填写 Cookie。

```bash
export QUARK_COOKIE="你的夸克网盘cookie"
./quark-mcp
```

不要把真实 Cookie 写进仓库、文档、命令行历史以外的公开位置，或 MCPB 包内容。

## 安装

GitHub Release 同时提供普通压缩包和 MCPB 包。选择时同时匹配**操作系统**和 **CPU 架构**。

| 系统 | 架构 | 普通压缩包 | MCPB |
|------|------|------------|------|
| macOS | Apple Silicon (arm64) | `quark-mcp_*_darwin-arm64.tar.gz` | `quark-mcp_*_darwin-arm64.mcpb` |
| macOS | Intel (amd64) | `quark-mcp_*_darwin-amd64.tar.gz` | `quark-mcp_*_darwin-amd64.mcpb` |
| Linux | x64 (amd64) | `quark-mcp_*_linux-amd64.tar.gz` | `quark-mcp_*_linux-amd64.mcpb` |
| Linux | ARM64 | `quark-mcp_*_linux-arm64.tar.gz` | `quark-mcp_*_linux-arm64.mcpb` |
| Windows | x64 (amd64) | `quark-mcp_*_windows-amd64.zip` | `quark-mcp_*_windows-amd64.mcpb` |
| Windows | ARM64 | `quark-mcp_*_windows-arm64.zip` | `quark-mcp_*_windows-arm64.mcpb` |

最新文件见 [Releases](https://github.com/chiehw/quark-mcp/releases)。MCPB 清单里的平台字段只区分操作系统，架构以 Release 文件名为准。

### MCPB 桌面扩展

`.mcpb` 可在支持 MCPB 的桌面客户端（例如 Claude Desktop）中直接安装：

1. 下载与本机系统和架构匹配的 `.mcpb`。
2. 打开该文件，或将其拖到客户端的 MCP 扩展安装入口。
3. 在安装配置界面填写夸克网盘 Cookie。该字段为必填敏感字符串，客户端会通过 `QUARK_COOKIE` 传给服务器。

### 普通二进制

解压后将 `quark-mcp`（Windows 为 `quark-mcp.exe`）放到可执行路径，然后用环境变量或 JSON 配置启动。

### 从源码编译

```bash
git clone https://github.com/chiehw/quark-mcp.git
cd quark-mcp
go build -o quark-mcp .
```

## Cookie 配置

优先级：

1. 环境变量 `QUARK_COOKIE`（非空时直接使用，不要求存在 JSON 文件）
2. `~/.quark-mcp/config.json`，或 `-config` 指定的 JSON 文件。若新路径不存在，仍会读取旧路径 `~/.quark-nd-disk/config.json`

两种来源都没有 Cookie 时，服务器会明确提示可设置 `QUARK_COOKIE` 或写入 JSON 配置。JSON 凭据文件会尽量收紧为仅当前用户可读写（`0600`）。

### 环境变量

```bash
export QUARK_COOKIE="你的夸克网盘cookie"
./quark-mcp
```

Windows PowerShell：

```powershell
$env:QUARK_COOKIE = "你的夸克网盘cookie"
.\quark-mcp.exe
```

### JSON 兼容配置

未设置 `QUARK_COOKIE` 时，继续读取 JSON 配置。默认写入 `~/.quark-mcp/config.json`：

```json
{
  "cookie": "你的夸克网盘cookie"
}
```

```bash
quark-mcp config init
quark-mcp config set cookie "你的夸克网盘cookie"
quark-mcp config show
./quark-mcp -config /path/to/config.json
```

`config get` / `config show` 会掩码显示 Cookie，不会完整输出。旧版 `~/.quark-nd-disk/config.json` 若仍存在且新路径尚未创建，会继续被读取。

### 如何获取 Cookie

1. 在浏览器中登录 [夸克网盘](https://pan.quark.cn)
2. 打开开发者工具（F12）
3. 切换到 Network（网络）标签
4. 刷新页面
5. 找到任意请求到 `drive.quark.cn` 的请求
6. 从请求头中复制 `Cookie` 的值

Cookie 等同于账号凭证。不要提交到 Git、不要写进 MCPB 包、不要出现在公开的 Registry 元数据或命令参数里。

## 使用方法

### 运行 MCP 服务器

```bash
./quark-mcp
```

或使用自定义配置路径：

```bash
./quark-mcp -config /path/to/config.json
```

### 配合 Claude Desktop 使用

在 Claude Desktop 配置文件中添加（macOS 路径：`~/Library/Application Support/Claude/claude_desktop_config.json`）：

```json
{
  "mcpServers": {
    "quark-mcp": {
      "command": "/path/to/quark-mcp",
      "env": {
        "QUARK_COOKIE": "你的夸克网盘cookie"
      }
    }
  }
}
```

如果已经用 MCPB 安装，客户端会代为写入启动命令和 `QUARK_COOKIE`，无需再手改这份 JSON。

## GitHub Release 与 MCPHub Registry

两者相关但不是同一件事：

- **GitHub Release**：推送 `v*` 标签后，GitHub Actions 会构建 6 个平台的普通压缩包和 `.mcpb`，校验 MCPB 后再附加到 Release。用户从这里下载安装包。
- **MCPHub Registry**（[mcphub.app](https://mcphub.app)）：公共目录收录需要维护者登录后单独提交，GitHub Actions **不会**自动上架。条目应链接本仓库和 Release 页面，安装说明需包含平台/架构选择、可执行文件安装和 `QUARK_COOKIE` 配置。公开元数据中不要填写真实 Cookie。

MCPB 文件不会自动完成 Registry 上架。

## MCPHub 自托管网关

MCPHub 网关启动 stdio Server 时，必须在**网关所在主机或容器**中找到对应平台的 `quark-mcp` 二进制。用户自己电脑上的二进制不能直接给远端网关使用。

在网关环境中：

1. 安装与网关操作系统/架构匹配的二进制。
2. 用环境变量注入 Cookie，例如 `QUARK_COOKIE`，不要把 Cookie 写进网关的公开配置或日志。
3. 以 stdio 方式启动 `quark-mcp`。

本次范围是 Registry 收录和本地 stdio 接入。如果需要 MCPHub Cloud 远程托管运行，需要另行设计 HTTP 传输、认证和部署。

## 可用工具

### list_files（列出文件）
列出指定目录下的所有文件和文件夹。

参数：
- `path`（字符串）：目录路径。使用 `/` 或空字符串表示根目录。例如：`/我的文档`

### create_folder（创建文件夹）
在指定路径创建新文件夹。

参数：
- `path`（字符串）：新文件夹的完整路径。例如：`/父文件夹/新文件夹` 或 `新文件夹`（在根目录创建）

### rename（重命名）
重命名文件或文件夹。

参数：
- `old_path`（字符串）：文件或文件夹的当前路径。例如：`/文件夹/旧名称`
- `new_name`（字符串）：新名称（不是完整路径）。例如：`新名称`

### delete（删除）
删除指定路径的文件或文件夹。

参数：
- `paths`（字符串数组）：要删除的路径列表。例如：`["/文件夹/文件.txt", "/文件夹/子文件夹"]`

### move（移动）
移动文件或文件夹到目标目录。

参数：
- `source_paths`（字符串数组）：要移动的源路径列表。例如：`["/文件夹/文件.txt"]`
- `dest_path`（字符串）：目标文件夹路径。使用 `/` 表示根目录。例如：`/目标目录`

### copy（复制）
复制文件或文件夹到目标目录。

参数：
- `source_paths`（字符串数组）：要复制的源路径列表。例如：`["/文件夹/文件.txt"]`
- `dest_path`（字符串）：目标文件夹路径。使用 `/` 表示根目录。例如：`/目标目录`

### get_info（获取信息）
获取文件或文件夹的详细信息。

参数：
- `path`（字符串）：文件或文件夹的路径。例如：`/文件夹/文件.txt`

### download_file（下载文件）
启动异步下载任务，将夸克网盘文件下载到本地。返回任务ID用于跟踪进度。

参数：
- `source_path`（字符串）：夸克网盘中的文件路径。例如：`/文件夹/文件.txt`
- `local_path`（字符串）：本地保存路径。例如：`/Users/name/Downloads/文件.txt`

返回：
- `task_id`（字符串）：用于 `get_download_status` 查询进度
- `file_size`（int64）：文件总大小（字节）

### get_download_status（获取下载状态）
查询下载任务的状态和进度。

参数：
- `task_id`（字符串）：`download_file` 返回的下载任务ID。例如：`dl_1234567890`

返回：
- `status`（字符串）：`pending`（等待）、`running`（运行中）、`completed`（已完成）、`failed`（失败）或 `canceled`（已取消）
- `progress`（字符串）：进度百分比。例如：`45.23%`
- `downloaded`（int64）：已下载字节数
- `total_size`（int64）：文件总大小（字节）
- `error`（字符串）：失败时的错误信息

### cancel_download（取消下载）
取消正在运行的下载任务。

参数：
- `task_id`（字符串）：要取消的下载任务ID。例如：`dl_1234567890`

### list_downloads（列出下载任务）
列出所有下载任务及其状态信息。

返回下载任务对象数组，包含状态、进度和文件信息。

### upload_file（上传文件）
将本地文件上传到夸克网盘。

参数：
- `dest_path`（字符串）：目标文件夹路径。使用 `/` 表示根目录。例如：`/我的文档`
- `local_path`（字符串）：本地文件路径。例如：`/Users/name/Downloads/文件.txt`

### regex_rename（正则重命名）
使用正则表达式批量重命名文件夹内的文件。

参数：
- `path`（字符串）：包含要重命名文件的目录路径。例如：`/photos`
- `pattern`（字符串）：匹配文件名的正则表达式模式。例如：`IMG_(\d+)`
- `replacement`（字符串）：替换字符串。使用 `$1`、`$2` 表示捕获组。例如：`Photo_$1`

### search（搜索）
在指定目录下递归搜索文件或文件夹。

参数：
- `path`（字符串）：搜索的目录路径。使用 `/` 表示根目录。例如：`/我的文档`
- `keyword`（字符串）：搜索关键词，匹配文件/文件夹名称

### get_share_file_list（获取分享文件列表）
获取夸克分享链接根目录的文件列表。

参数：
- `share_url`（字符串）：夸克分享链接。例如：`https://pan.quark.cn/s/abc123` 或带密码 `https://pan.quark.cn/s/abc123?pwd=xyz`

返回：
- 文件列表，包含 id、name、size、is_folder、category、时间戳等信息

### get_share_detail（获取分享详情）
获取分享链接中指定文件夹内的文件列表。

参数：
- `share_url`（字符串）：夸克分享链接。例如：`https://pan.quark.cn/s/abc123`
- `folder_id`（字符串）：文件夹ID（fid）。使用 `0` 表示根目录，或传入文件夹的 fid 浏览子目录

返回：
- 文件列表，包含 id、name、size、is_folder、category、时间戳等信息

### save_from_share（保存分享文件）
将分享链接中的文件保存到自己的夸克网盘。

参数：
- `share_url`（字符串）：夸克分享链接。例如：`https://pan.quark.cn/s/abc123`
- `folder_id`（字符串，可选）：文件所在的文件夹ID。使用 `0` 表示根目录（默认）。如果要保存子目录中的文件，请先通过 `get_share_detail` 获取文件夹的 fid。例如：`abc123def456`
- `dest_path`（字符串）：目标保存路径。使用 `/` 表示根目录。例如：`/来自分享`
- `file_ids`（字符串数组，可选）：要保存的文件ID列表。留空则保存该文件夹下全部文件。例如：`["file_id_1", "file_id_2"]`

返回：
- `task_id`（字符串）：保存任务ID
- `saved_count`（int）：已保存文件数量
- `saved_files`（字符串数组）：已保存的文件名列表
- `folder_id`（字符串）：源文件夹ID

## 下载流程示例

```
1. download_file(source_path="/视频/电影.mp4", local_path="/tmp/电影.mp4")
   → 返回：{ "task_id": "dl_1709876543210", "file_size": 1500000000 }

2. get_download_status(task_id="dl_1709876543210")
   → 返回：{ "status": "running", "progress": "35.50%", "downloaded": 532500000 }

3. get_download_status(task_id="dl_1709876543210")
   → 返回：{ "status": "completed", "progress": "100.00%" }
```

## 分享文件操作示例

### 保存根目录全部文件
```
save_from_share(share_url="https://pan.quark.cn/s/abc123", dest_path="/来自分享")
→ 返回：{ "saved_count": 5, "saved_files": ["文件1.mp4", "文件夹A", ...] }
```

### 保存根目录特定文件
```
1. get_share_file_list(share_url="https://pan.quark.cn/s/abc123")
   → 返回：[{ "id": "xxx", "name": "文件.mp4", "is_folder": false }, ...]

2. save_from_share(share_url="https://pan.quark.cn/s/abc123", dest_path="/来自分享", file_ids=["xxx"])
   → 返回：{ "saved_count": 1, "saved_files": ["文件.mp4"] }
```

### 保存子目录中的特定文件
```
1. get_share_file_list(share_url="https://pan.quark.cn/s/abc123")
   → 返回：[{ "id": "folder_abc", "name": "文件夹A", "is_folder": true }, ...]

2. get_share_detail(share_url="https://pan.quark.cn/s/abc123", folder_id="folder_abc")
   → 返回：[{ "id": "file_xyz", "name": "电影.mp4", "is_folder": false }, ...]

3. save_from_share(share_url="https://pan.quark.cn/s/abc123", folder_id="folder_abc", file_ids=["file_xyz"], dest_path="/来自分享")
   → 返回：{ "saved_count": 1, "saved_files": ["电影.mp4"], "folder_id": "folder_abc" }
```

## 开发构建

```bash
go test ./...
go build -o quark-mcp .
```

本地打包全部 MCPB 产物（需要 Node.js，以及全局 `mcpb` 或 `npx`）：

```bash
./scripts/pack-mcpb.sh 1.3.0
```

产物写入 `mcpb-dist/`，文件名包含系统和架构。推送 `v*` 标签后，GitHub Actions 会先跑测试并校验 MCPB，再和 GoReleaser 普通压缩包一起发布；MCPB 校验失败时不会发布不完整的 MCPB 资产。

## 常见问题

**如何选择安装包？**  
先看操作系统，再看 CPU：Apple Silicon 用 `darwin-arm64`，Intel Mac 用 `darwin-amd64`，常见 Linux/Windows 电脑用 `amd64`，ARM 设备用 `arm64`。

**MCPB 和 tar.gz/zip 有什么区别？**  
MCPB 面向支持该格式的桌面客户端，安装时收集 Cookie。普通压缩包适合手动放置二进制，或在 MCPHub 自托管网关里启动 stdio Server。

**为什么 Registry 上看不到这个项目？**  
Registry 上架需要维护者在 MCPHub 单独提交，不会随 GitHub Release 自动完成。

**网关报找不到可执行文件？**  
把对应平台的二进制放到网关运行环境中，不要指望用户本地安装的文件能被远端网关调用。

**Cookie 无效或过期？**  
重新从浏览器复制 Cookie，更新 `QUARK_COOKIE` 或 JSON 配置后重启 Server。不要在日志里打印 Cookie。

## 许可证

[MIT](LICENSE)
