# eudic-mcp-go

欧路词典（Eudic）生词本 OpenAPI 的 [MCP](https://modelcontextprotocol.io/) 服务端，使用官方 [Go SDK](https://github.com/modelcontextprotocol/go-sdk) 实现。

本项目受 [safeblood/eudic-mcp](https://github.com/safeblood/eudic-mcp) 启发：工具命名对齐，Go 单二进制，**无需 Python**。适合 Cursor 本机与 **Remote SSH（Linux）**。

## 功能

| 工具 | 说明 |
|------|------|
| `eudic_list_categories` | 列出生词本分组 |
| `eudic_create_category` | 创建分组（如 `english-coach`） |
| `eudic_rename_category` / `eudic_delete_category` | 重命名 / 删除分组 |
| `eudic_list_words` / `eudic_add_word` / `eudic_add_words_bulk` / `eudic_delete_words` | 单词读写 |
| `eudic_get_word` / `eudic_list_mastered_words` | 查询单词 / 已掌握 |
| `eudic_list_notes` / `eudic_get_note` / `eudic_add_note` / `eudic_delete_note` | 笔记 |
| `eudic_list_sentences` | 用户例句列表 |

## 构建

需要 Go 1.22+。

```bash
# 当前平台
go build -ldflags="-s -w" -o eudic-mcp-go .

# 多平台（含 Cursor Remote 用的 Linux）
bash scripts/build.sh          # Linux/macOS
# 或
powershell -File scripts/build.ps1   # Windows
```

产物：

| 路径 | 用途 |
|------|------|
| `dist/linux-amd64/eudic-mcp-go` | 远程 Linux x86_64 |
| `dist/linux-arm64/eudic-mcp-go` | 远程 Linux arm64 |
| `dist/windows-amd64/eudic-mcp-go.exe` | 本机 Windows |

## Token

1. [OpenAPI 授权](https://my.eudic.net/OpenAPI/Authorization)  
2. 只填 `NIS` **后面**的字符串（程序会自动加 `NIS `）

## Cursor 配置

### Remote SSH / Linux（MCP 跑在远程）

把 Linux 二进制放到远程，例如 `~/bin/eudic-mcp-go`，编辑**远程** `~/.cursor/mcp.json`：

```json
{
  "mcpServers": {
    "eudic": {
      "command": "/home/你的用户名/bin/eudic-mcp-go",
      "env": {
        "EUDIC_API_TOKEN": "你的token"
      }
    }
  }
}
```

然后 Reload Window。详细步骤见 english-coach 仓库 `docs/eudic-mcp-setup.md`。

### Windows 本机

```json
{
  "mcpServers": {
    "eudic": {
      "command": "C:\\data\\code\\go\\src\\eudic-mcp-go\\eudic-mcp-go.exe",
      "env": {
        "EUDIC_API_TOKEN": "你的token"
      }
    }
  }
}
```

## english-coach

生词本名 **`english-coach`**：`list_categories` → 必要时 `create_category` → `add_word` / `add_note`。

## 致谢

感谢 [safeblood/eudic-mcp](https://github.com/safeblood/eudic-mcp)。

## License

MIT
