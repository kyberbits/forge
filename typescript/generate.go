package typescript

import (
	"fmt"
	"reflect"
	"strings"
)

func Generate(goStructs map[string]interface{}) Interfaces {
	tsInterfaces := Interfaces{}

	for nameFromMap, goStruct := range goStructs {
		rv := reflect.ValueOf(goStruct)

		tsInterface := Interface{
			Name: nameFromMap,
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

			// Maps need additional handling
			if valueField.Kind() == reflect.Map {
				tsInterface.Fields = append(
					tsInterface.Fields,
					genTSMapField(typeField.Name, valueField, tag))

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

func genTSMapField(goFieldName string, mapField reflect.Value, tag jsonFieldTag) Field {
	tsField := Field{}
	tsField.Name = getTSFieldName(goFieldName, tag)

	mapIndex := TranslateReflectTypeString(mapField.Type().Key().String())
	mapValue := TranslateReflectTypeString(mapField.Type().Elem().String())

	if strings.HasPrefix(mapValue, "*") {
		mapValue = strings.TrimPrefix(mapValue, "*")
		mapValue += " | null"
	}

	tsField.Type = fmt.Sprintf("Record<%s, %s>", mapIndex, mapValue)

	if tag.Omitempty {
		tsField.Optional = true
	}

	return tsField
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
	isSlice := strings.HasPrefix(reflectTypeString, "[]")
	if isSlice {
		reflectTypeString = strings.TrimPrefix(reflectTypeString, "[]")
	}

	tsType := translateReflectTypeString(reflectTypeString)
	if isSlice {
		tsType += "[]"
	}

	return tsType
}

func translateReflectTypeString(reflectTypeString string) string {
	switch reflectTypeString {
	case "interface {}":
		return "any"
	case "int", "int8", "int16", "int32", "int64":
		return "number"
	case "uint", "uint8", "uint16", "uint32", "uint64", "float", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	default:
		// If not detected above, assume this is a custom struct. Trim the package name
		if i := strings.LastIndex(reflectTypeString, "."); i >= 0 {
			return reflectTypeString[i+1:]
		}
	}

	return reflectTypeString
}
