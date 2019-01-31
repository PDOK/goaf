package main

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
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

	limitParam := operation.Parameters.GetByInAndName("query", "limit")
	if limitParam != nil {
		offsetParameter := &openapi3.ParameterRef{
			Value: &openapi3.Parameter{Name: "offset", In: "query", Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}}},
		}
		operation.Parameters = append(operation.Parameters, offsetParameter)
	}

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
		println(ref.Ref)
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
		zplit := strings.Split(ref.Ref, "/")
		if pointer {
			builder.WriteString("*")
		}
		builder.WriteString(normalize(zplit[len(zplit)-1]))
	}
}

func responses(responses openapi3.Responses) string {
	log.Println("Not implemented.")
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

func downloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	filepath := "spec/wfs3.0.yml"
	downloadFile(filepath, "https://raw.githubusercontent.com/opengeospatial/WFS_FES/master/openapi.yaml")
	yaml, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatalf("Cannot find file %s", filepath)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromYAMLData(yaml)

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

func createFile(outputPath, outputFileName, templateFileName string, swagger *openapi3.Swagger, templates *template.Template) {
	out, err := os.Create(fmt.Sprintf("%s/%s", outputPath, outputFileName))
	if err != nil {
		log.Println("%v", err)
	}
	defer out.Close()
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
