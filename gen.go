package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:generate go run gen.go

func upperFirst(title string) string {
	return strings.Title(title)
}

func normalize(s string) string {
	return strings.Replace(strings.Title(strings.Replace(s, "-", " ", -1)), " ", "", -1)
}

func pathparams(s string) string {
	matches := regexp.MustCompile(`{([^}]*)}`).FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		s = strings.Replace(s, match[0], ".*", -1)
	}

	return s
}

func convfunc(path string, operation openapi3.Operation) string {

	if len(operation.Parameters) == 0 {
		return "// no parameters present"
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("func ParametersFor%s(r *http.Request) ", normalize(operation.OperationID)))

	// parameters
	builder.WriteString("(")
	paramsBuilder(operation.Parameters, &builder, true)
	builder.WriteString(") {\n")
	pathExists := false
	for _, param := range operation.Parameters {
		p := param.Value

		if p.In == "path" {
			if !pathExists {
				builder.WriteString("    pathSplit:=strings.Split(r.URL.Path, \"/\") \n")
				pathExists = true
			}
			pathSplit := strings.Split(path, "/") // definition e.g. // path: /collections/{collectionId}/items/{featureId}
			matches := regexp.MustCompile(`{([^}]*)}`).FindAllStringSubmatch(path, -1)

			for _, match := range matches {
				if p.Name == match[1] {
					for i, v := range pathSplit {
						if v == match[0] {
							builder.WriteString(fmt.Sprintf("    %s=pathSplit[%d]", p.Name, i))
						}
					}

				}
			}

		} else if p.In == "query" {
			builder.WriteString(fmt.Sprintf("    %sArray, ok := r.URL.Query()[\"%s\"]\n", p.Name, p.Name))
			builder.WriteString("    if ok {\n")
			builder.WriteString(fmt.Sprintf("        %s= %sArray[0]\n", p.Name, p.Name))
			builder.WriteString("          }\n")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("    return \n}")
	return builder.String()
}

func paramsBuilder(parameters openapi3.Parameters, builder *strings.Builder, stringOnly bool) {
	for i, v := range parameters {
		p := v.Value

		builder.WriteString(p.Name)
		builder.WriteString(" ")

		convertSchemaType(p.Schema.Value, builder, stringOnly)

		if i < len(parameters)-1 {
			builder.WriteString(", ")
		}
	}
}

func property(properties map[string]*openapi3.SchemaRef, required []string) string {
	var builder strings.Builder

	for k, schemaRef := range properties {
		if schemaRef.Value.Description != "" {
			builder.WriteString("	/* " + strings.ReplaceAll(schemaRef.Value.Description, "\n", "\n    "))
			builder.WriteString("	*/\n")
		}
		//log.Println(k)
		builder.WriteString("	")
		builder.WriteString(strings.Title(k))
		builder.WriteString(" ")
		convertSchemaRefType(schemaRef, &builder, true)
		builder.WriteString(" ")
		if contains(required, k) {
			builder.WriteString(fmt.Sprintf("`json:\"%s\"`", k))
		} else {
			builder.WriteString(fmt.Sprintf("`json:\"%s,omitempty\"`", k))
		}

		builder.WriteString("\n")

	}

	return builder.String()
}

func convertSchemaRefType(ref *openapi3.SchemaRef, builder *strings.Builder, pointer bool) {
	if ref.Ref == "" {
		schema := ref.Value
		schemaType := schema.Type
		switch schemaType {
		case "number":
			builder.WriteString("float64")
		case "integer":
			builder.WriteString("int64")
		case "object":
			builder.WriteString("interface{}")
		case "array":
			builder.WriteString("[")
			builder.WriteString("]")
			convertSchemaRefType(schema.Items, builder, false)
		default:
			builder.WriteString("string")
		}
	} else {
		println(ref.Ref)
		zplit := strings.Split(ref.Ref, "/")
		if pointer {
			builder.WriteString("*")
		}
		builder.WriteString(normalize(zplit[len(zplit)-1]))
	}
}

func responses(responses openapi3.Responses) string {
	if responses != nil {
		log.Println("Not implemented.")
	}
	return ""
}

func convertSchemaType(schema *openapi3.Schema, builder *strings.Builder, stringOnly bool) {
	schemaType := getType(stringOnly, schema.Type)
	switch schemaType {
	case "number":
		builder.WriteString("float64")
	case "array":
		builder.WriteString("[")
		builder.WriteString("]")
		convertSchemaType(schema.Items.Value, builder, stringOnly)
	default:
		builder.WriteString("string")
	}
}

func getType(stringOnly bool, schemaType string) string {
	if stringOnly {
		return "string"
	}
	return schemaType
}

func main_() {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	u, _ := url.Parse("https://raw.githubusercontent.com/opengeospatial/ogcapi-features/master/core/openapi/ogcapi-features-1.yaml")
	components, err := loader.LoadFromURI(u)

	// file based spec based upon https://github.com/opengeospatial/WFS_FES/blob/master/core/examples/openapi/ogcapi-features-1-example1.yaml
	filepath := "spec/oaf.yml"
	swagger, err := loader.LoadFromFile(filepath)
	if err != nil {
		log.Fatalf("Got error reading swagger file %v", err)
		return
	}

	// merge
	for k, v := range swagger.Components.Parameters {
		components.Components.Parameters[k] = v
	}
	swagger.Components = components.Components

	out, err := json.Marshal(swagger)

	err = ioutil.WriteFile("spec/oaf.json", out, 0644)
	if err != nil {
		log.Fatalf("Got error writing combined swagger file %v", err)
		return
	}
	// filter out geojeson stuff and non complex objects
	// geojson type are represented in extra_types file in implementation
	refs := swagger.Components.Schemas
	for k, v := range refs {

		// rest of GeoJSON is omitted
		if strings.Contains(k, "GeoJSON") {
			delete(refs, k)
			continue
		}
		// omit non complex types
		if v.Value.Type != "object" {
			delete(refs, k)
			continue
		}

		// promote complex properties to structs
		properties := v.Value.Properties
		for k, schemaRef := range properties {
			if schemaRef.Value.Type == "object" {
				refs[k] = schemaRef
				schemaRef.Ref = fmt.Sprintf("#/components/schemas/%s", k)
			}

		}
	}

	if err != nil {
		log.Fatalf("Got error reading swagger file %v", err)
		return
	}

	var templates = template.Must(template.New("templates").Funcs(
		template.FuncMap{
			"cf":         convfunc,
			"upperFirst": upperFirst,
			"normalize":  normalize,
			"pathparams": pathparams,
			"type":       property,
			"responses":  responses,
			"dict": func(values ...interface{}) (map[string]interface{}, error) {
				if len(values)%2 != 0 {
					return nil, errors.New("invalid dict call")
				}
				dict := make(map[string]interface{}, len(values)/2)
				for i := 0; i < len(values); i += 2 {
					key, ok := values[i].(string)
					if !ok {
						return nil, errors.New("dict keys must be strings")
					}
					dict[key] = values[i+1]
				}
				return dict, nil
			},
		}).ParseGlob("codegen_templates/*"))

	createFile("codegen", "provider.go", "interface.tpl", swagger, templates)
	createFile("codegen", "types.go", "types.tpl", swagger, templates)
	createFile("server", "routing.gen.go", "routing.tpl", swagger, templates)
}

func createFile(outputPath, outputFileName, templateFileName string, swagger *openapi3.T, templates *template.Template) {
	out, err := os.Create(fmt.Sprintf("%s/%s", outputPath, outputFileName))
	if err != nil {
		log.Printf("%v", err)
	}
	defer closeFile(out)
	err = templates.ExecuteTemplate(out, templateFileName, swagger)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func closeFile(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Printf("%v", err)
	}
}
