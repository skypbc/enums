package parse

import (
	"github.com/skypbc/enums/internal/utils"
	"slices"
)

func Strings(data map[string]any, path []string) (res []EnumFile) {
	for key, val := range data {
		if sub, ok := val.(map[string]any); ok {
			exportPath := append([]string{}, path...)
			exportPath = append(exportPath, key, "_")
			settings := EnumFileSettings{
				Path: exportPath,
			}

			var path_ []string
			if path[0] == "enums" {
				// убираем "enums" + одно под категорию (например "int" или "string")
				path_ = slices.Clone(path[2:])
			} else {
				// убираем "strings" из пути
				path_ = slices.Clone(path[1:])
			}

			enumFile := EnumFile{
				Type:     "string",
				Path:     path_,
				Names:    []string{key},
				Settings: settings,
			}
			enumFile.Settings.Object.Name = key

			sub_ := map[string]any{}
			for name, raw := range sub {
				if name == "_" {
					continue
				}

				var value string
				var exist bool

				switch x := raw.(type) {
				case map[string]any:
					if raw, ok := x["_"]; ok {
						value, exist = raw.(string)
					}
					sub_[name] = x
				case string:
					value = x
					exist = true
				default:
					sub_[name] = x
				}

				if !exist {
					continue
				}

				safe := utils.SanitizeName(name)
				enumFile.Items = append(enumFile.Items, EnumItem{Name: safe, Value: value})
			}

			if len(enumFile.Items) > 0 {
				res = append(res, enumFile)
			}

			if len(sub_) > 0 {
				res = append(res, Strings(sub_, append(path, key))...)
			}
		}
	}

	return res
}
