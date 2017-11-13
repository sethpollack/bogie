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
	f := func(text string, w io.Writer) error {
		hasContent, buff, err := runTemplate(c, b, text)
		if err != nil {
			return err
		}

		if hasContent {
			if _, err := io.Copy(w, buff); err != nil {
				return err
			}
		}

		return nil
	}

	return template.FuncMap{
		"latestImage": ecr.LatestImage(b.SkipImageLookup),
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
