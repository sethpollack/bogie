package main

import "text/template"

func initFuncs(o *BogieOpts) template.FuncMap {
	env := &Env{}
	typeconv := &TypeConv{}
	file := &File{}
	ecr := &EcrInit{}

	ecr.ecrInit.Do(ecr.initEcr)

	f := template.FuncMap{
		"latestImage":  ecr.client.LatestImage,
		"readDir":      file.ReadDir(o),
		"readFile":     file.ReadFile(o),
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
