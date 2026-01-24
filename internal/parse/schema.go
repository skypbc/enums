package parse

import (
	"slices"
	"strings"
)

func Schema(schema map[string]any) []EnumFile {
	var enumFiles []EnumFile
	if nums, ok := schema["nums"].(map[string]any); ok {
		enumFiles = append(enumFiles, Nums(nums, []string{"nums"})...)
	}
	if strs, ok := schema["strings"].(map[string]any); ok {
		enumFiles = append(enumFiles, Strings(strs, []string{"strings"})...)
	}
	if numStrings, ok := schema["num_strings"].(map[string]any); ok {
		enumFiles = append(enumFiles, NumStrings(numStrings, []string{"num_strings"})...)
	}
	if enums, ok := schema["enums"].(map[string]any); ok {
		if intEnums, ok := enums["int"].(map[string]any); ok {
			enumFiles = append(enumFiles, Nums(intEnums, []string{"enums", "int"})...)
		}
		if strEnums, ok := enums["string"].(map[string]any); ok {
			enumFiles = append(enumFiles, Strings(strEnums, []string{"enums", "int"})...)
		}
	}
	return normalizeEnums(enumFiles)
}

func normalizeEnums(files []EnumFile) []EnumFile {
	for index := range files {
		enumFile := &files[index]
		isNone := false
		for _, item := range enumFile.Items {
			if strings.ToUpper(item.Name) == "NONE" {
				isNone = true
				break
			}
		}
		if !isNone {
			switch enumFile.Type {
			case "int":
				enumFile.Items = append(enumFile.Items, EnumItem{Name: "NONE", Value: 0})
			case "string":
				enumFile.Items = append(enumFile.Items, EnumItem{Name: "NONE", Value: ""})
			}
		}
		slices.SortFunc(enumFile.Items, func(a, b EnumItem) int {
			switch x := a.Value.(type) {
			case int:
				if y, ok := b.Value.(int); ok {
					if x == y {
						return strings.Compare(a.Name, b.Name)
					}
					if x < y {
						return -1
					}
					return 1
				}
			case int64:
				if y, ok := b.Value.(int64); ok {
					if x == y {
						return strings.Compare(a.Name, b.Name)
					}
					if x < y {
						return -1
					}
					return 1
				}
			case string:
				if y, ok := b.Value.(string); ok {
					return strings.Compare(x, y)
				}
			}
			return strings.Compare(a.Name, b.Name)
		})
	}
	return files
}
