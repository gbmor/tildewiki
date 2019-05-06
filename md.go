package main

import (
	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

// Sets parameters for the markdown->html renderer
func setupMarkdown(css, title string) *bf.HTMLRenderer {
	// if using local CSS file, use the virtually-served css
	// path rather than the actual file name
	if cssLocal([]byte(confVars.cssPath)) {
		css = "/css"
	}

	// return the parameters used for the rendering
	// of markdown to html.
	var params = bf.HTMLRendererParameters{
		CSS:   css,
		Title: title,
		Icon:  "/icon",
		Meta: map[string]string{
			"name=\"application-name\"": "TildeWiki " + twvers + " :: https://github.com/gbmor/tildewiki",
			"name=\"viewport\"":         "width=device-width, initial-scale=1.0",
		},
		Flags: bf.CompletePage | bf.Safelink,
	}
	return bf.NewHTMLRenderer(params)
}

// Wrapper function to generate the parameters above and
// pass them to the blackfriday library's parsing function
func render(data []byte, title string) []byte {
	return bf.Run(data, bf.WithRenderer(setupMarkdown(confVars.cssPath, title)))
}
