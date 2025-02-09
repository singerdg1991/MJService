package engine

import (
	"os"
	"regexp"
	"strings"
)

func ExportAPIDocsYaml(dest, content string) error {
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func IsIgnoredFile(filePath string) bool {
	// check in not valid slice
	for _, v := range IGNORED_FILES_TO_PARS {
		if strings.Contains(filePath, v) {
			return true
		}
	}
	return false
}

func ParseStructTagValues(tag string) OpenApiFieldTagValues {
	values := OpenApiFieldTagValues{}
	if tag != "" {
		for _, item := range strings.Split(tag, ";") {
			splitted := strings.Split(item, ":")
			if len(splitted) == 1 {
				values.Required = TerIf(splitted[0] == "required", true, false)
				values.Nullable = TerIf(splitted[0] == "nullable", true, false)
				values.Ignored = TerIf(splitted[0] == "ignored", true, false)
				continue
			}
			regexTagValueSplit := regexp.MustCompile(`(?sm)^(.*?):(.*?)$`)
			tabSplitted := regexTagValueSplit.FindStringSubmatch(item)
			if len(tabSplitted) > 2 {
				value := tabSplitted[2]
				values.In = TerIf(tabSplitted[1] == "in", value, values.In)
				values.Example = TerIf(tabSplitted[1] == "example", value, values.Example)
				values.Ref = TerIf(tabSplitted[1] == "$ref", value, values.Ref)
				values.MaxLength = TerIf(tabSplitted[1] == "maxLength", value, values.MaxLength)
				values.MinLength = TerIf(tabSplitted[1] == "minLength", value, values.MinLength)
				values.Minimum = TerIf(tabSplitted[1] == "minimum", value, values.Minimum)
				values.Maximum = TerIf(tabSplitted[1] == "maximum", value, values.Maximum)
				values.Pattern = TerIf(tabSplitted[1] == "pattern", value, values.Pattern)
				values.Format = TerIf(tabSplitted[1] == "format", value, values.Format)
			}
		}
	}
	return values
}

func SanitizeCommentLineText(str string) string {
	str = RemoveNewLines(str)
	return str
}

type MergeMapType interface {
	Schema | Operations | Response
}

func MergeMaps[T MergeMapType](src map[string]T, dest map[string]T) map[string]T {
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

func GetResponseDescription(statusCode string) string {
	if desc, ok := ResponseDescriptions[statusCode]; ok {
		return desc
	}
	return "Unknown response code"
}

func GenerateOperationId(method string, path string) string {
	return strings.Join(strings.Split(path, "/"), "_") + RestActions[strings.ToUpper(method)] + RestOperations[strings.ToUpper(method)]
}
