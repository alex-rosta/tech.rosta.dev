package markdown

import (
	"bytes"
	"os"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkHtml "github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.DefinitionList,
			emoji.Emoji,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("paraiso-dark"),
				highlighting.WithCSSWriter(&bytes.Buffer{}),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkHtml.WithHardWraps(),
			goldmarkHtml.WithXHTML(),
		),
	)

	// Generate CSS file for syntax highlighting
	GenerateChromaCSS()
}

func ToHTML(content string) string {
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return content // fallback to plain text
	}
	return buf.String()
}

// GenerateChromaCSS generates the CSS for syntax highlighting
func GenerateChromaCSS() {
	// Create web/static directory if it doesn't exist
	os.MkdirAll("web/static", 0755)

	// Create CSS file
	file, err := os.Create("web/static/syntax-highlight.css")
	if err != nil {
		return
	}
	defer file.Close()

	style := styles.Get("paraiso-dark")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithClasses(true),
		html.PreventSurroundingPre(true),
	)

	err = formatter.WriteCSS(file, style)
	if err != nil {
		return
	}
}
