package main

import (
	"net/http"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/lang"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/render/mhtml"
)

func render(w http.ResponseWriter, r *http.Request, buf []byte, pathname string) {
	ext := path.Ext(pathname)
	switch ext {
	case ".md", ".markdown", ".txt", ".text":
		renderer, doc := newRendererMmark(buf, pathname)
		x := markdown.Render(doc, renderer)
		w.Write(x)
	default:
		w.Write(buf)
	}

}

func newRendererMmark(buf []byte, pathname string) (markdown.Renderer, ast.Node) {
	init := mparser.NewInitial(pathname)
	p := parser.NewWithExtensions(mparser.Extensions)
	parserFlags := parser.FlagsNone
	p.Opts = parser.Options{
		ReadIncludeFn: init.ReadInclude,
		Flags:         parserFlags,
	}

	doc := markdown.Parse(buf, p)

	mparser.AddBibliography(doc)
	mparser.AddIndex(doc)

	mhtmlOpts := mhtml.RendererOptions{
		Language: lang.New("en"), // TODO(miek): should come from xdoc.yaml.
	}
	opts := html.RendererOptions{
		Comments:       [][]byte{[]byte("//"), []byte("#")},
		RenderNodeHook: mhtmlOpts.RenderHook,
		Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks,
		Generator:      `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.miek.nl`,
	}
	opts.Flags |= html.CompletePage

	return html.NewRenderer(opts), doc
}
