# Bogie
A templating tool for Kubernetes clusters.

## Running Bogie
Bogie can be run with a single template, but, is most useful with a manifest file, which describes the entire cluster state, and allows for greater template reuse. Generated manifests can be written to a directory, a single file, or stdout.

Run the following manifest with `bogie template -m path/to/manifest.yaml`.

```yaml
out_path: releases
out_file: release.yaml
out_format: dir
env_file: path/to/global/vars/values.yaml
ignore_file: ./.bogieignore
skip_image_lookup: false
app_regex: .*
applications:
- name: my-templates
  templates: path/to/templates
  values:
  - path/to/templates/values.yaml
  - path/to/templates/env.values.yaml
  override_vars:
  - app.secrets.key=value
- name: my-other-templates
  templates: path/to/templates
  env: production
  override_vars:
  - app.secrets.key=value
- name: my-other-other-templates
  templates: path/to/templates
```

> Note: Bogie supports file paths and remote Github repos for `templates`, `values`, `_helpers.tmpl`, and `.bogieignore` files. For private Github repos export a `GITHUB_TOKEN` variable. For more information on creating a personal API token see https://github.com/blog/1509-personal-api-tokens.

The manifest file supports the following keys.

 - `out_format` output format, (`dir`|`file`|`stdout`).
 - `out_path` output directory for `file` and `dir` outputs.
 - `out_file` output file name.
 - `env_file` global values file.
 - `ignore_file` global `.bogieignore` file.
 - `skip_image_lookup` skip the image lookup in `latestImage` template function.
 - `app_regex` only build apps with names that match regex
 - `applications` list of applications.
	 - `name` name of the output directory. (the directory maintains the original structure with the exception of this name)
	 - `templates` templates directory.
	 - `values` list of values files, when omitted, Bogie will attempt to load `values.yaml` from the templates directory root.
	 - `env` when set Bogie will attempt to load `env.values.yaml` from the templates directory root.
	 - `override_vars` list of individual values in the format of `app.secrets.key=value`.

Bogie can also be run without a manifest:

```sh
bogie template \
  -t path/to/templates \
  -v path/to/templates/env.values.yaml \
  -v path/to/templates/values.yaml \
  -e path/to/global/vars/values.yaml \
  -o dir
```

To see the full list of supported flags run `bogie template help`.

## Directory Structure
Bogie does not enforce a strict template directory structure. It walks the template directory and sends the files to the rendering engine. The only files required to be in template directory root are `_helpers.tmpl` and `.bogieignore`. When no values are specified bogie attempt to load `values.yaml` from the template directory root and when an `env` is provided it will also attempt to load `env.values.yaml`.

## Templates
Bogie uses the Golang `text/template` package for templating and defaults to triple-curlies `{{{ .Values.key }}}` so that manifests containing Golang templates can be generated.

## Values
Bogie values files are [sops](https://github.com/mozilla/sops) encrypted yaml files. Bogie will automatically decrypt the files and Unmasrshal them into a `Values` map when generating the templates and can be accessed in the templates like `{{{ .Values.key }}}`. The files get merged in the following order:

  1) The global values file.
  2) The values.yaml file (when no values files provided).
  3) The values files in order that they are provided.
  4) The env.values.yaml. (when env is provided).

Finally override vars in the format of `path.to.key=value` get merged into the struct.

## Partials
Partials can be added to the templates directory root as `_helpers.tmpl` files, and will be merged in with the templates.

## Ignore Files
Files excluded in `.bogieignore` will not be templated.

## Functions
Bogie adds all the template functions from [Masterminds/sprig](github.com/Masterminds/sprig) and adds a few custom template functions as well:

- `latestImage` - `bogie:{{{ latestImage "latest" }}}` using commit SHA tags in kubernetes deployments and `latest` in the manifest will trigger a new deployment whenever the manifest is applied, even if nothing changed. `latestImage` looks up the corresponding commit SHA tag based on `latest` or `staging` etc.
- `readDir` - `{{{ readDir "path/to/dir" }}}` returns a map of file names and file contents.
- `decryptDir` - `{{{ decryptDir "path/to/dir" }}}` returns a map of file names and dycrypted file contents.
- `readFile` - `{{{ readFile "path/to/file" }}}` returns file contents.
- `decryptFile` - `{{{ decryptFile "path/to/file" }}}` returns dycrypted file contents.
- `basicAuth` - `{{{ basicAuth .Values.user .Values.password }}}` returns a SHA-1 hash of your password in the format of `user:{SHA}password`.
- `json` - `{{{ json .Values.json_string }}}` converts a JSON string into an object.
- `jsonArray` - `{{{ json .Values.json_array_string }}}` converts a JSON string into a slice.
- `toJSON` - `{{{ toJSON .Values.object }}}` converts an object to a JSON document.
- `yaml` - `{{{ yaml .Values.yaml_string }}}` converts a YAML string into an object.
- `yamlArray` - `{{{ yaml .Values.yaml_array_string }}}` converts a YAML string into a slice.
- `toYAML` - `{{{ toYAML .Values.object }}}` converts an object to a YAML document.
- `toml` - `{{{ yaml .Values.toml_string }}}` converts a TOML string into an object.
- `toTOML` - `{{{ toTOML .Values.object }}}` converts an object to a TOML document.

## Credits
- This project began as a fork of [Gomplate](https://github.com/hairyhenderson/gomplate) but at this point is almost a complete rewrite.
- The ignore package was copied from [Helm](https://github.com/kubernetes/helm) and slightly modified.



