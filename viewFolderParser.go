package trinity

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ViewsSuffix = ".ghtml"
)

// viewFolderParser is used to parse the views folder
// Expected folder structure:
//   Controller1
//     ->  Action1.ghtml
//     ->  Action2.ghtml
type viewFolderParser struct {
	viewsFolder string
	mvcI        *MvcInfrastructure
}

func newViewFolderParser(mvcI *MvcInfrastructure) *viewFolderParser {
	parser := new(viewFolderParser)

	parser.viewsFolder = mvcI.viewsFolder
	parser.mvcI = mvcI

	return parser
}

func (parser *viewFolderParser) parse() error {
	logger.Trace("")

	controllerNames, err := parser.getControllerNames()
	if err != nil {
		return err
	}

	for _, controllerName := range controllerNames {
		controller := Controller(controllerName)
		logger.Debug("Controller: %s", controller)

		actionNames, err := parser.getActionNames(controllerName)
		if err != nil {
			return err
		}

		for _, actionName := range actionNames {
			action := Action(actionName)
			logger.Debug("Action: %s", action)

			templatePath := filepath.Join(parser.viewsFolder, controllerName, actionName+ViewsSuffix)
			logger.Debug("TemplatePath: %v", templatePath)

			parser.mvcI.bindView(controller, action, templatePath)
		}
	}

	return nil
}

func (parser *viewFolderParser) getControllerNames() ([]string, error) {
	logger.Trace("")

	result := make([]string, 0)

	viewsFolder, err := os.Open(parser.viewsFolder)
	if err != nil {
		return nil, err
	}

	viewsFolderStat, err := viewsFolder.Stat()
	if err != nil {
		return nil, err
	}

	if !viewsFolderStat.IsDir() {
		return nil, errors.New("mvc's viewsFolder isn't folder")
	}

	folderStats, err := viewsFolder.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, folderStat := range folderStats {
		logger.Debug("check %s", folderStat.Name())

		if !folderStat.IsDir() {
			continue
		}

		result = append(result, folderStat.Name())
	}

	return result, nil
}

func (parser *viewFolderParser) getActionNames(controllerName string) ([]string, error) {
	logger.Trace("")

	result := make([]string, 0)

	folder, err := os.Open(filepath.Join(parser.viewsFolder, controllerName))
	if err != nil {
		return nil, err
	}

	fileStats, err := folder.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, fileStat := range fileStats {
		logger.Debug("check %s", fileStat.Name())

		if fileStat.IsDir() {
			continue
		}

		fileName := fileStat.Name()
		if !strings.HasSuffix(fileName, ViewsSuffix) {
			continue
		}

		result = append(result, fileName[:len(fileName)-len(ViewsSuffix)])
	}

	return result, nil
}
