package typescript

import (
	"reflect"
	"regexp"
	"strings"
)

func Generate(goStructs []interface{}) Interfaces {
	tsInterfaces := Interfaces{}

	for _, goStruct := range goStructs {
		rv := reflect.ValueOf(goStruct)

		tsInterface := Interface{
			Name: rv.Type().Name(),
		}

		for i := 0; i < rv.NumField(); i++ {
			valueField := rv.Field(i)
			if !valueField.CanInterface() {
				continue
			}

			typeField := rv.Type().Field(i)
			tag := parseJSONFieldTag(typeField.Tag.Get("json"))

			if tag.Ignored {
				continue
			}

			tsInterface.Fields = append(
				tsInterface.Fields,
				genTSField(
					typeField.Name,
					typeField.Type.String(),
					tag,
				),
			)
		}

		tsInterfaces = append(tsInterfaces, tsInterface)
	}

	return tsInterfaces
}

func genTSField(goFieldName string, goFieldType string, tag jsonFieldTag) Field {
	tsField := Field{}
	if strings.HasPrefix(goFieldType, "*") {
		tsField.Null = true
		goFieldType = strings.TrimPrefix(goFieldType, "*")
	}

	tsField.Name = getTSFieldName(goFieldName, tag)
	tsField.Type = getTSFieldType(goFieldType, tag)

	if tag.Omitempty {
		tsField.Optional = true
	}

	return tsField
}

func getTSFieldName(goFieldName string, tag jsonFieldTag) string {
	if tag.NameOverride != "" {
		return tag.NameOverride
	}

	return goFieldName
}

func getTSFieldType(goFieldType string, tag jsonFieldTag) string {
	if tag.TypeOverride != "" {
		goFieldType = tag.TypeOverride
	}

	return TranslateReflectTypeString(goFieldType)
}

func TranslateReflectTypeString(reflectTypeString string) string {
	switch reflectTypeString {
	case "interface {}":
		return "any"
	case "int", "int8", "int16", "int32", "int64":
		return "number"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	}

	re := regexp.MustCompile(`(?m)^(\[\])?(\w+\.)?(\w+)$`)
	x := re.FindStringSubmatch(reflectTypeString)

	return x[3] + x[1]
}
