package bogie

import (
	"io"
	"text/template"

	"github.com/sethpollack/bogie/crypto"
	"github.com/sethpollack/bogie/ecr"
	"github.com/sethpollack/bogie/file"
	"github.com/sethpollack/bogie/types"
)

func initFuncs(c *context, b *Bogie) template.FuncMap {
	f := func(text string, out io.Writer) {
		runTemplate(c, b, text, out)
	}

	return template.FuncMap{
		"latestImage": ecr.LatestImage,
		"readDir":     file.ReadDir(f),
		"readFile":    file.ReadFile(f),
		"decryptFile": file.DecryptFile(f),
		"decryptDir":  file.DecryptDir(f),
		"basicAuth":   crypto.BasicAuth,
		"json":        types.JSON,
		"jsonArray":   types.JSONArray,
		"toJSON":      types.ToJSON,
		"yaml":        types.YAML,
		"yamlArray":   types.YAMLArray,
		"toYAML":      types.ToYAML,
		"toml":        types.TOML,
		"toTOML":      types.ToTOML,
	}
}
