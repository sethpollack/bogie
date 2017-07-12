package main

import "text/template"

func initFuncs(c *Context, b *Bogie) template.FuncMap {
	typeconv := &TypeConv{}
	file := &File{}
	ecr := &EcrInit{}
	env := &Env{}

	ecr.ecrInit.Do(ecr.initEcr)

	f := template.FuncMap{
		"latestImage":  ecr.client.LatestImage,
		"readDir":      file.ReadDir(c, b),
		"readFile":     file.ReadFile(c, b),
		"getenv":       env.Getenv,
		"json":         typeconv.JSON,
		"jsonArray":    typeconv.JSONArray,
		"toJSON":       typeconv.ToJSON,
		"toJSONPretty": typeconv.toJSONPretty,
		"yaml":         typeconv.YAML,
		"yamlArray":    typeconv.YAMLArray,
		"toYAML":       typeconv.ToYAML,
		"toml":         typeconv.TOML,
		"toTOML":       typeconv.ToTOML,
	}

	return f
}
