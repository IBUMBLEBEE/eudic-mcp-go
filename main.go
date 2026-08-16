package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"eudic-mcp-go/internal/eudic"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "eudic-mcp-go",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_list_categories",
		Description: "获取所有生词本分组（category）列表。",
	}, listCategories)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_create_category",
		Description: "创建新的生词本分组。",
	}, createCategory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_rename_category",
		Description: "重命名生词本分组。",
	}, renameCategory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_delete_category",
		Description: "删除生词本分组。",
	}, deleteCategory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_list_words",
		Description: "获取指定生词本中的单词列表。",
	}, listWords)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_add_word",
		Description: "添加单个单词到生词本（可带 context_line 例句）。",
	}, addWord)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_add_words_bulk",
		Description: "批量导入单词到指定生词本（已存在不会重复添加）。",
	}, addWordsBulk)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_delete_words",
		Description: "从指定生词本中删除单词。",
	}, deleteWords)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_get_word",
		Description: "查询某个单词是否已在生词本中。",
	}, getWord)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_list_mastered_words",
		Description: "查询已掌握单词列表。",
	}, listMasteredWords)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_list_notes",
		Description: "获取所有单词笔记列表。",
	}, listNotes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_get_note",
		Description: "获取某个单词的笔记。",
	}, getNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_add_note",
		Description: "为某个单词添加或更新笔记。",
	}, addNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_delete_note",
		Description: "删除某个单词的笔记。",
	}, deleteNote)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "eudic_list_sentences",
		Description: "获取用户例句列表。",
	}, listSentences)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func clientOrErr() (*eudic.Client, *mcp.CallToolResult, error) {
	c, err := eudic.NewClientFromEnv()
	if err != nil {
		return nil, &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			IsError: true,
		}, nil
	}
	return c, nil, nil
}

func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

func apiFail(err error) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		IsError: true,
	}, nil, nil
}

type languageIn struct {
	Language string `json:"language,omitempty" jsonschema:"语言代码 en/fr/de/es，默认 en"`
}

type createCategoryIn struct {
	Name     string `json:"name" jsonschema:"新分组名称，例如 english-coach"`
	Language string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type renameCategoryIn struct {
	CategoryID any    `json:"category_id" jsonschema:"分组 id"`
	Name       string `json:"name" jsonschema:"新名称"`
	Language   string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type deleteCategoryIn struct {
	CategoryID any    `json:"category_id" jsonschema:"分组 id"`
	Name       string `json:"name" jsonschema:"分组名称"`
	Language   string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type listWordsIn struct {
	CategoryID any    `json:"category_id,omitempty" jsonschema:"分组 id，默认 0"`
	Language   string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
	Page       int    `json:"page,omitempty" jsonschema:"页码，默认 0"`
	PageSize   int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 100"`
}

type addWordIn struct {
	Word         string `json:"word" jsonschema:"要添加的单词"`
	CategoryIDs  []any  `json:"category_ids,omitempty" jsonschema:"分组 id 列表"`
	Star         int    `json:"star,omitempty" jsonschema:"星级 1-5，默认 2"`
	ContextLine  string `json:"context_line,omitempty" jsonschema:"语境例句"`
	Language     string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type addWordsBulkIn struct {
	Words      []string `json:"words" jsonschema:"单词数组"`
	CategoryID any      `json:"category_id,omitempty" jsonschema:"分组 id，默认 0"`
	Language   string   `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type deleteWordsIn struct {
	Words      []string `json:"words" jsonschema:"要删除的单词"`
	CategoryID any      `json:"category_id,omitempty" jsonschema:"分组 id，默认 0"`
	Language   string   `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type wordIn struct {
	Word     string `json:"word" jsonschema:"单词"`
	Language string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

type listMasteredIn struct {
	Language   string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
	CategoryID any    `json:"category_id,omitempty" jsonschema:"可选分组 id"`
	Page       int    `json:"page,omitempty" jsonschema:"页码，默认 0"`
	PageSize   int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 100"`
}

type pageIn struct {
	Language string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
	Page     int    `json:"page,omitempty" jsonschema:"页码，默认 0"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 100"`
}

type addNoteIn struct {
	Word     string `json:"word" jsonschema:"单词"`
	Note     string `json:"note" jsonschema:"笔记内容"`
	Language string `json:"language,omitempty" jsonschema:"语言代码，默认 en"`
}

func listCategories(ctx context.Context, _ *mcp.CallToolRequest, in languageIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.ListCategories(ctx, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func createCategory(ctx context.Context, _ *mcp.CallToolRequest, in createCategoryIn) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Name) == "" {
		return apiFail(fmt.Errorf("name is required"))
	}
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.CreateCategory(ctx, in.Name, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func renameCategory(ctx context.Context, _ *mcp.CallToolRequest, in renameCategoryIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.RenameCategory(ctx, eudic.IDString(in.CategoryID), in.Name, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func deleteCategory(ctx context.Context, _ *mcp.CallToolRequest, in deleteCategoryIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.DeleteCategory(ctx, eudic.IDString(in.CategoryID), in.Name, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func listWords(ctx context.Context, _ *mcp.CallToolRequest, in listWordsIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.ListWords(ctx, eudic.IDString(in.CategoryID), in.Language, in.Page, in.PageSize)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func addWord(ctx context.Context, _ *mcp.CallToolRequest, in addWordIn) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Word) == "" {
		return apiFail(fmt.Errorf("word is required"))
	}
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	ids := make([]string, 0, len(in.CategoryIDs))
	for _, id := range in.CategoryIDs {
		ids = append(ids, eudic.IDString(id))
	}
	out, err := c.AddWord(ctx, in.Word, in.Language, ids, in.Star, in.ContextLine)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func addWordsBulk(ctx context.Context, _ *mcp.CallToolRequest, in addWordsBulkIn) (*mcp.CallToolResult, any, error) {
	if len(in.Words) == 0 {
		return apiFail(fmt.Errorf("words is required"))
	}
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.AddWordsBulk(ctx, in.Words, eudic.IDString(in.CategoryID), in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func deleteWords(ctx context.Context, _ *mcp.CallToolRequest, in deleteWordsIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.DeleteWords(ctx, in.Words, eudic.IDString(in.CategoryID), in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func getWord(ctx context.Context, _ *mcp.CallToolRequest, in wordIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.GetWord(ctx, in.Word, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func listMasteredWords(ctx context.Context, _ *mcp.CallToolRequest, in listMasteredIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	cat := ""
	if in.CategoryID != nil {
		cat = eudic.IDString(in.CategoryID)
	}
	out, err := c.ListMasteredWords(ctx, in.Language, cat, in.Page, in.PageSize)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func listNotes(ctx context.Context, _ *mcp.CallToolRequest, in pageIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.ListNotes(ctx, in.Language, in.Page, in.PageSize)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func getNote(ctx context.Context, _ *mcp.CallToolRequest, in wordIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.GetNote(ctx, in.Word, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func addNote(ctx context.Context, _ *mcp.CallToolRequest, in addNoteIn) (*mcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Word) == "" || strings.TrimSpace(in.Note) == "" {
		return apiFail(fmt.Errorf("word and note are required"))
	}
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.AddNote(ctx, in.Word, in.Note, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func deleteNote(ctx context.Context, _ *mcp.CallToolRequest, in wordIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.DeleteNote(ctx, in.Word, in.Language)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}

func listSentences(ctx context.Context, _ *mcp.CallToolRequest, in pageIn) (*mcp.CallToolResult, any, error) {
	c, fail, err := clientOrErr()
	if fail != nil || err != nil {
		return fail, nil, err
	}
	out, err := c.ListSentences(ctx, in.Language, in.Page, in.PageSize)
	if err != nil {
		return apiFail(err)
	}
	return jsonResult(out)
}
