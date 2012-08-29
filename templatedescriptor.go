package trinity

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	allCharExceptCloseBrace = "[^}]*"

	// Name of define with options
	ViewOptionsTemplate = "ViewOptions"

	regexpOptionsTepmplate = regexp.MustCompile("{{" + allCharExceptCloseBrace + ViewOptionsTemplate + allCharExceptCloseBrace + "}}")
	regexpEnd              = regexp.MustCompile("{{\\s*end\\s*}}")

	MasterPageOption         = "MasterPage"
	AdditionalTemplateOption = "AdditionalTemplate"
)

// templateDescriptor stores the options section for a template
type templateDescriptor struct {
	templatePath string

	additionalTemplates []string
	masterPage          string
}

// newTemplateDescriptor creates a new templateDescriptor object
func newTemplateDescriptor(viewsFolder string, templatePath string) (*templateDescriptor, error) {
	logger.Trace("")

	template := new(templateDescriptor)

	template.templatePath = templatePath
	template.additionalTemplates = make([]string, 0)

	err := template.parseOptions(viewsFolder)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (template *templateDescriptor) parseOptions(viewsFolder string) error {
	file, err := os.Open(template.templatePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	optionLines := template.extractOptionLines(string(bytes))

	for _, line := range optionLines {
		logger.Debugf("line: %s", line)

		if strings.HasPrefix(line, "{") {
			break
		}

		vals := strings.Split(line, "=")
		if len(vals) != 2 {
			continue
		}

		switch vals[0] {
		case MasterPageOption:
			{
				template.masterPage = filepath.Join(viewsFolder, vals[1])

				break
			}
		case AdditionalTemplateOption:
			{
				logger.Debugf("add template: %s", vals[1])
				template.additionalTemplates = append(template.additionalTemplates, filepath.Join(viewsFolder, vals[1]))
				break
			}
		}
	}

	return nil
}

func (template *templateDescriptor) extractOptionLines(text string) []string {
	logger.Trace("")

	emptyResult := make([]string, 0)

	indexes := regexpOptionsTepmplate.FindAllStringIndex(text, -1)
	logger.Debugf("indexes: %v", indexes)

	if len(indexes) == 0 {
		return emptyResult
	}

	indexesEnds := regexpEnd.FindAllStringIndex(text, -1)
	logger.Debugf("indexesEnds: %v", indexesEnds)

	if len(indexes) == 0 {
		return emptyResult
	}

	contentStartsAt := indexes[0][1]
	contentEndsAt := 0
	for _, indexEnd := range indexesEnds {
		if indexEnd[0] < contentStartsAt {
			continue
		}

		if (contentEndsAt == 0) || (contentEndsAt > indexEnd[0]) {
			contentEndsAt = indexEnd[0]
		}
	}

	if contentEndsAt == 0 {
		return emptyResult
	}

	result := make([]string, 0)
	for _, line := range strings.Split(text[contentStartsAt:contentEndsAt], "\n") {
		result = append(result, strings.TrimSpace(line))
	}

	return result
}
