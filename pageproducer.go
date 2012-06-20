package trinity

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

// renderPage renders a page with any dependencies (like master pages or template pages
// for inner elements).
func renderPage(vm interface{}, templateDescr *templateDescriptor) (html []byte, err error) {
	logger.Trace("")
	logger.Debug("index: %s, master: %s", templateDescr.templatePath, templateDescr.masterPage)

	var pageTemplate *template.Template = nil
	if templateDescr.masterPage != "" {
		logger.Trace("parse master page")
		pageTemplate, err = template.ParseFiles(templateDescr.masterPage)
		if err != nil {
			return nil, err
		}
		
		logger.Trace("funcs")
		pageTemplate.Funcs(template.FuncMap{"equals": equals})

		logger.Trace("parse page to master")
		_, err = pageTemplate.ParseFiles(templateDescr.templatePath)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Trace("parse page without master")
		pageTemplate, err = template.ParseFiles(templateDescr.templatePath)
		if err != nil {
			return nil, err
		}
	}
	
	for _, pagePath := range templateDescr.additionalTemplates {
		logger.Debug("dep: %s", pagePath)

		_, err = pageTemplate.ParseFiles(pagePath)
		if err != nil {
			return nil, err
		}
	}

	logger.Trace("template execute")
	var htmlBuffer bytes.Buffer
	err = pageTemplate.Execute(&htmlBuffer, vm)
	if err != nil {
		return nil, err
	}

	if err := recover(); err != nil {
		return nil, errors.New(fmt.Sprintf("Panic: %s", err))
	}

	return htmlBuffer.Bytes(), nil
}


func equals(args ...interface{}) bool {
	if len(args) != 2 {
		return false
	}
	
	return args[0] == args[1]
}