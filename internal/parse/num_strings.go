package parse

import (
	"github.com/skypbc/enums/internal/utils"
	"github.com/skypbc/goutils/gnum"
	"slices"
)

func NumStrings(data map[string]any, path []string) (res []EnumFile) {
	for key, val := range data {
		if sub, ok := val.(map[string]any); ok {
			if utils.IsLeafDict(sub) {
				exportPath := append([]string{}, path...)
				exportPath = append(exportPath, key, "_")

				var path_ []string
				if path[0] == "enums" {
					// убираем "enums" + одно под категорию (например "int" или "string")
					path_ = slices.Clone(path[2:])
				} else {
					// убираем "num_strings" из пути
					path_ = slices.Clone(path[1:])
				}

				enumIdFile := EnumFile{
					Type:  "int",
					Path:  slices.Clone(path_),
					Names: []string{key},
					Settings: EnumFileSettings{
						Path: exportPath,
					},
				}
				enumIdFile.Settings.Object.Name = key
				enumIdFile.Settings.Object.Postfix = "id"

				enumTagFile := EnumFile{
					Type:  "string",
					Path:  slices.Clone(path_),
					Names: []string{key},
					Settings: EnumFileSettings{
						Path: exportPath,
					},
				}
				enumTagFile.Settings.Object.Name = key
				enumTagFile.Settings.Object.Postfix = "tag"

				for name, raw := range sub {
					safe := utils.SanitizeName(name)
					if pair, ok := raw.([]any); ok && len(pair) >= 2 {
						id, ok := gnum.TryInt(pair[0])
						if !ok {
							continue
						}
						tag, ok := pair[1].(string)
						if !ok {
							continue
						}
						enumIdFile.Items = append(enumIdFile.Items, EnumItem{Name: safe, Value: id})
						enumTagFile.Items = append(enumTagFile.Items, EnumItem{Name: safe, Value: tag})
					}
				}

				res = append(res, enumIdFile, enumTagFile)
			} else {
				res = append(res, NumStrings(sub, append(path, key))...)
			}
		}
	}
	return res
}
