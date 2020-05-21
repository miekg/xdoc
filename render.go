package main

import (
	"io"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mparser"
)

func (d *Doc) render(w http.ResponseWriter, r *http.Request, proj GitLab, path string) {
	fullpath := d.FullPath(proj, path)

	renderer := newRendererMmark(fullpath)
	x := markdown.Render(doc, renderer)
	io.WriteString(w, x)
}

func newRendererMmark(path string) *markdown.Renderer {
	init := mparser.NewInitial(path)
	p := parser.NewWithExtensions(mparser.Extensions)
	parserFlags := parser.FlagsNone
	parserFlags |= parser.SkipFootnoteList
	p.Opts = parser.Options{
		ReadIncludeFn: init.ReadInclude,
		Flags:         parserFlags,
	}

	doc := markdown.Parse(buf, p)
	mparser.AddBibliography(doc)
	mparser.AddIndex(doc)
	mhtmlOpts := mhtml.RendererOptions{
		Language: lang.New(documentLanguage),
	}
	opts := html.RendererOptions{
		Comments:  [][]byte{[]byte("//"), []byte("#")},
		Flags:     html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks,
		Generator: `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.miek.nl`,
	}
	opts.Flags |= html.CompletePage

	return html.NewRenderer(opts)
}
