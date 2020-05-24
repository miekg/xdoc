package main

type Asset struct {
	contenttype string
	content     string
}

// Assets hold various assets like css contents.
var Assets = map[string]Asset{
	"style.css": styleAsset,
}

var styleAsset = Asset{
	contenttype: "text/css; charset=UTF-8",
	content: `
html, body { height: 100%; }
`,
}
