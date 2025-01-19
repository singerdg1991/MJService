package openengine

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/hoitek/OpenEngine/engine"
)

var WgEnum = &sync.WaitGroup{}

func (p *openEngine) extractEnumNamesFromComments(enumsFilePath string) ([]string, error) {
	// Create the AST by parsing src.
	fileSet := token.NewFileSet()

	// structNames is a list of all the models in the file
	var structNames []string

	// Parse the file containing this struct definitions
	f, err := parser.ParseFile(fileSet, enumsFilePath, nil, parser.ParseComments)
	if err != nil {
		return structNames, err
	}
	fileName := filepath.Base(enumsFilePath)
	// If the file has no comments we return an error
	if len(f.Comments) == 0 {
		return structNames, engine.BuildError(fileName, "reading comment file: there is no @api declarations in the file")
	}
	// Loop through all the comments in the file
	for _, comment := range f.Comments {
		// log.Printf("comment %#v\n", comment)
		// If the comment list is empty we return an error
		if len(comment.List) == 0 {
			return structNames, engine.BuildError(fileName, "reading comment.List: there is no @api declarations in the file")
		}
		// log.Println("commentList", comment.List)
		// Loop through all the comment lines from the list
		for _, commentLine := range comment.List {
			// log.Println("commentLine", commentLine)
			// optional capture second item group with @apiDefine
			reg := regexp.MustCompile(engine.API_ENUM_DATA_REGEXP)
			// Sanitize the comment line text and remove new lines and spaces to make regexp work
			commentLineText := engine.SanitizeCommentLineText(commentLine.Text)
			// get second capture from regexp
			structName := reg.FindStringSubmatch(commentLineText)
			// If we have a second capture we append it to the modelNames list
			if structName != nil {
				// If the second capture is not empty we append it to the modelNames list
				if len(structName) > 1 {
					structNames = append(structNames, structName[2])
					break
				}
			}
		}
	}

	// If we modelNames are empty we don't have any model, so we return error
	if len(structNames) == 0 {
		return structNames, engine.BuildError(fileName, "parsing commentLine: there is no @api declarations in the file")
	}

	// return the modelNames list
	// log.Println("structNames", structNames)
	return structNames, nil
}

// TODO: refactor this function
func (p *openEngine) mapEnumFieldsToSchemaDict(list []*ast.Field, structName string, enumsDict *engine.SchemasDict) {
	var fieldNames []string
	var enumType = ""
	for _, field := range list {
		var fieldName = field.Names[0].Name
		if field.Tag != nil {
			// Get the tag value and remove start and end quotes
			var fieldTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

			// json tag parsed
			jsonFieldName := strings.Split(fieldTag.Get(engine.JSON_TAG_NAME), ",")[0]

			// json fieldName if json tag is not empty
			fieldName = engine.TerIf(jsonFieldName == "", field.Names[0].Name, strings.Split(jsonFieldName, ",")[0])
		}
		if enumType == "" {
			enumType = field.Type.(*ast.Ident).Name
		}

		fieldNames = append(fieldNames, fieldName)
	}

	enum := (*enumsDict)[structName]
	enum.Type = enumType
	enum.Enum = fieldNames
	(*enumsDict)[structName] = enum
}

func (p *openEngine) extractEnumsDictFromFile(enumsFilePath string) (engine.SchemasDict, error) {
	var enumsDict = engine.SchemasDict{}
	// Create the AST by parsing src.
	fileSet := token.NewFileSet()

	// Parse the file containing this struct definitions
	f, err := parser.ParseFile(fileSet, enumsFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// structNames, err := ExtractDataFromComments(enumsFilePath)
	structNames, err := p.extractEnumNamesFromComments(enumsFilePath)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Loop through all the declarations in the file
	for _, decl := range f.Decls {
		// If the declaration is a GenDecl (general declaration)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Loop through all the Specs and grab the one we want.
			for _, spec := range genDecl.Specs {
				// If the spec is a TypeSpec we are half way there.
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// If the Type of the TypeSpec is a StructType, we found the struct
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// check if modelNames have the typeSpec.Name.Name, If not we continue
						if !engine.StringInSlice(typeSpec.Name.Name, &structNames) {
							continue
						}
						// Get struct name
						structName := typeSpec.Name.Name
						// Create a new Enum for the struct
						enumsDict[structName] = engine.Schema{
							// This is a hack to get the type of the struct
							Type:       "string",
							Properties: engine.Properties{},
							Enum:       []string{},
						}
						// Loop through all the fields in the struct
						p.mapEnumFieldsToSchemaDict(structType.Fields.List, structName, &enumsDict)
					}
				}
			}
		}
	}
	// log.Println("SchemasDict", enumsDict)
	return enumsDict, nil
}

func (p *openEngine) extractEnumsFromDirectory(structsDirPath string, chanEnums chan engine.ChanSchemas) (engine.SchemasDict, error) {
	// Done waitgroup for this goroutine
	defer WgEnum.Done()

	// Create AllSchemasDict
	var AllSchemasDict = engine.SchemasDict{}

	// Get all the files in the models directory
	files, err := os.ReadDir(structsDirPath)

	// If we have an error, we return it
	if err != nil {
		// Send ChanSchemas to channel if we have channel
		if chanEnums != nil {
			chanEnums <- engine.ChanSchemas{
				Items: AllSchemasDict,
				Err:   err,
			}
		}

		// Return AllSchemasDict
		return AllSchemasDict, err
	}

	// Loop through all the files
	for _, file := range files {

		// If the file is a directory or is an ignored file we continue
		if file.IsDir() || engine.IsIgnoredFile(file.Name()) {
			continue
		}

		// Extract the enums from the file
		enumsDict, err := p.extractEnumsDictFromFile(structsDirPath + "/" + file.Name())

		// If we have an error, we return it
		if err != nil {
			// TODO: we should have verbose param to log more
			// log.Println(err)
			continue
		}

		// Append the enums to the global enums map
		AllSchemasDict = engine.MergeMaps(AllSchemasDict, enumsDict)
	}
	// Send ChanSchemas to channel if we have channel
	if chanEnums != nil {
		chanEnums <- engine.ChanSchemas{
			Items: AllSchemasDict,
			Err:   nil,
		}
	}

	// Return AllSchemasDict
	return AllSchemasDict, nil
}

func (p *openEngine) AddEnums(enumsDict engine.SchemasDict) OpenEngine {
	p.Components.Schemas = engine.MergeMaps(p.Components.Schemas, enumsDict)
	return p
}

func (p *openEngine) ParseEnums(baseDirectory string, allIgnoredPaths ...[]string) OpenEngine {
	// Create Mutex for SchemasDict
	mx := &sync.Mutex{}

	// ReInitialize WaitGroup for goroutines
	WgEnum = &sync.WaitGroup{}

	// Create SchemasDict
	AllSchemasDict := p.Components.Schemas
	var ignoredPaths []string
	if len(allIgnoredPaths) > 0 {
		ignoredPaths = allIgnoredPaths[0]
	}
	mergedIgnoredPaths := append(p.GeneralIgnoredPaths, ignoredPaths...)
	// Find all the directories in the baseDirPath
	structsDirectoryPaths, err := engine.FindAllDirectoriesInPath(baseDirectory, &mergedIgnoredPaths)

	// If we have an error, we return it
	if err != nil {
		p.err = err
		return p
	}

	// Create a channel for the enums
	var chanEnums = make(chan engine.ChanSchemas, len(structsDirectoryPaths))

	// Close the channel when we are done with it in goroutine
	go func() {
		defer close(chanEnums)
		WgEnum.Wait()
	}()

	// Loop through all the models paths
	for _, structsDirectoryPath := range structsDirectoryPaths {
		// Add 1 to WaitGroup
		WgEnum.Add(1)

		// Extract the enums from the directory in a goroutine
		go p.extractEnumsFromDirectory(structsDirectoryPath, chanEnums)
	}

	// Loop through all the enums in the channel
	for result := range chanEnums {
		// Lock the mutex for the SchemasDict
		mx.Lock()

		// Check if we have an error
		if result.Err != nil {
			p.err = result.Err
			return p
		}

		// Merge the enums to the global enums map
		AllSchemasDict = engine.MergeMaps(AllSchemasDict, result.Items)

		// Unlock the mutex for the SchemasDict
		mx.Unlock()
	}

	// Return the global enums map
	p.Components.Schemas = AllSchemasDict
	return p
}
