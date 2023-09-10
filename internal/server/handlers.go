package server

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/biozz/paste/internal/storage_types"
	"github.com/labstack/echo/v4"
)

const (
	defaultNameLength = 10
	maxContentLength  = 20000
)

func (h *Web) Index(c echo.Context) error {
	ctx := map[string]string{"Slug": generateName(defaultNameLength)}
	return c.Render(http.StatusOK, "index", ctx)
}

type PasteIn struct {
	Content string `form:"content"`
	Slug    string `form:"slug"`
}

func (h *Web) New(c echo.Context) error {
	var pasteIn PasteIn
	err := c.Bind(&pasteIn)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if len(pasteIn.Content) > maxContentLength {
		return c.String(http.StatusBadRequest, "too long")
	}
	paste := storage_types.Paste{
		Content: pasteIn.Content,
		Slug:    pasteIn.Slug,
	}
	if paste.Slug == "" {
		paste.Slug = generateName(defaultNameLength)
	}
	paste, err = h.storage.CreatePaste(context.Background(), paste)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	pastePath := "/" + paste.Slug
	c.Response().Header().Set("HX-Redirect", "/"+paste.Slug)
	pasteURL := h.conf.BaseURL + pastePath
	return c.String(http.StatusOK, pasteURL)
}

func (h *Web) Redirect(c echo.Context) error {
	slug := c.Param("slug")
	paste, err := h.storage.GetPasteBySlug(context.Background(), slug)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Redirect(http.StatusTemporaryRedirect, paste.Content)
}

type ViewContext struct {
	Content template.HTML
	CSS     template.CSS
}

func (h *Web) View(c echo.Context) error {
	slug := c.Param("slug")
	slugParts := strings.Split(slug, ".")
	ext := ""
	if len(slugParts) == 2 {
		slug = slugParts[0]
		ext = slugParts[1]
	}
	paste, err := h.storage.GetPasteBySlug(context.Background(), slug)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if ext == "txt" {
		return c.String(http.StatusOK, paste.Content)
	}
	var lexer chroma.Lexer
	for _, name := range lexers.Names(true) {
		if name == ext {
			lexer = lexers.Get(name)
			break
		}
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	iterator, err := lexer.Tokenise(nil, paste.Content)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	content := bytes.NewBufferString("")
	css := bytes.NewBufferString("")

	noBackground := map[chroma.TokenType]string{
		chroma.Background: "background-color: transparent",
	}
	formatter := html.New(
		html.WithClasses(true),
		html.WithCustomCSS(noBackground),
		html.WithLineNumbers(true),
	)
	err = formatter.WriteCSS(css, style)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = formatter.Format(content, style, iterator)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctx := ViewContext{
		Content: template.HTML(content.String()),
		CSS:     template.CSS(css.String()),
	}
	return c.Render(http.StatusOK, "view", ctx)
}
