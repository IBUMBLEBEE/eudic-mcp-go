# eudic-mcp-go

欧路词典（Eudic）生词本 OpenAPI 的 [MCP](https://modelcontextprotocol.io/) 服务端，使用官方 [Go SDK](https://github.com/modelcontextprotocol/go-sdk) 实现。

本项目受 [safeblood/eudic-mcp](https://github.com/safeblood/eudic-mcp) 启发：在功能与工具命名上与其对齐，并用 Go 重写为单个可执行文件，面向 Cursor / Claude Desktop 等客户端，**无需本机 Python**。

## 功能

工具名与 [safeblood/eudic-mcp](https://github.com/safeblood/eudic-mcp) 对齐，便于 english-coach 等 skill 直接切换：

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

```powershell
cd C:\data\code\go\src\eudic-mcp-go
go build -o eudic-mcp-go.exe .
```

需要 Go 1.22+（开发机）。运行时用户只需 `eudic-mcp-go.exe`。

## 配置 Token

1. 打开 [OpenAPI 授权页](https://my.eudic.net/OpenAPI/Authorization) 获取 Token  
2. **只填 `NIS` 后面的字符串**，不要带 `NIS ` 前缀（程序会自动拼接）

## Cursor MCP

编辑 `%USERPROFILE%\.cursor\mcp.json`：

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

然后**完全重启 Cursor**。

## english-coach

配合 skill 时，生词本名称使用 **`english-coach`**：先 `eudic_list_categories`，没有则 `eudic_create_category`，再 `eudic_add_word` / `eudic_add_note`。

## 致谢

感谢 [safeblood/eudic-mcp](https://github.com/safeblood/eudic-mcp) 项目提供的思路与工具设计；本仓库在此基础上用 Go 实现了兼容的 MCP 服务端。

## License

MIT
