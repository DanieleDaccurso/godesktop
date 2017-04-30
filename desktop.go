package godesktop

import (
	"io/ioutil"
	"strings"
)

type DesktopFile map[string]string

func (d DesktopFile) Get(key string) string {
	return d[key]
}

func (d DesktopFile) GetLocalizedOrFallback(key, locale string) string {
	val, exists := d[key+"["+locale+"]"]
	if exists {
		return val
	}
	return d[key]
}

func (d DesktopFile) IsTerminal() bool {
	return d["Terminal"] == "true"
}

func GetFiles(dir string) ([]DesktopFile, error) {
	var files []DesktopFile

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		bytes, fileErr := ioutil.ReadFile(dir + "/" +fileInfo.Name())
		if fileErr != nil {
			return files, fileErr
		}
		file, fileError := Parse(string(bytes))
		if fileError != nil {
			return files, fileError
		}
		files = append(files, file)
	}

	return files, nil
}

func Parse(content string) (DesktopFile, error) {
	var d DesktopFile = DesktopFile{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// ignore header and empty lines
		if strings.Contains(line, "[Desktop Entry]") || strings.Trim(line, " \n\r\t") == "" {
			continue
		}

		// if the line contains a comment, remove the comment
		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
		}

		// ignore invalid line
		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			println("err:" , line)
			continue
		}
		d[parts[0]] = parts[1]
	}

	return d, nil
}