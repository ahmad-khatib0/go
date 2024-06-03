package validate

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
)

func ValidateHostURL(origUrl string) (string, error) {
	hostUrl, err := url.Parse(origUrl)
	if err != nil {
		return "", err
	}
	if !hostUrl.IsAbs() {
		return "", errors.New("host_url must be absolute")
	}
	if hostUrl.Hostname() == "" {
		return "", errors.New("invalid host_url")
	}
	if hostUrl.Fragment != "" {
		return "", errors.New("fragment is not allowed in host_url")
	}
	if hostUrl.Path == "" {
		hostUrl.Path = "/"
	}
	return hostUrl.String(), nil
}

func ExecuteTemplate(template *template.Template, parts []string, params map[string]interface{}) (map[string]string, error) {
	content := map[string]string{}
	buffer := new(bytes.Buffer)

	if parts == nil {
		if err := template.Execute(buffer, params); err != nil {
			return nil, err
		}
		content[""] = buffer.String()
	} else {
		for _, part := range parts {
			buffer.Reset()
			if templBody := template.Lookup(part); templBody != nil {
				if err := templBody.Execute(buffer, params); err != nil {
					return nil, err
				}
			}
			content[part] = buffer.String()
		}
	}

	return content, nil
}

func ResolveTemplatePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	curwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(curwd, path)), nil
}

func ReadTemplateFile(pathTempl *template.Template, lang string) (*template.Template, string, error) {
	buffer := bytes.Buffer{}
	err := pathTempl.Execute(&buffer, map[string]interface{}{"Language": lang})
	path := buffer.String()
	if err != nil {
		return nil, path, fmt.Errorf("reading %s: %w", path, err)
	}

	templ, err := template.ParseFiles(path)
	return templ, path, err
}
