package parse

import (
	"github.com/skypbc/enums/internal/utils"
	"github.com/skypbc/goutils/gnum"
	"slices"
)

func Nums(data map[string]any, path []string) (res []EnumFile) {
	for key, val := range data {
		if sub, ok := val.(map[string]any); ok {
			if utils.IsFlatIntMap(sub) {
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
					// убираем "nums" из пути
					path_ = slices.Clone(path[1:])
				}

				enumFile := EnumFile{
					Type:     "int",
					Path:     path_,
					Names:    []string{key},
					Settings: settings,
				}
				enumFile.Settings.Object.Name = key
				for name, raw := range sub {
					if name == "_" {
						continue
					}
					safe := utils.SanitizeName(name)
					id, ok := gnum.TryInt(raw)
					if !ok {
						continue
					}
					enumFile.Items = append(enumFile.Items, EnumItem{Name: safe, Value: id})
				}
				res = append(res, enumFile)
			} else {
				res = append(res, Nums(sub, append(path, key))...)
			}
		}
	}
	return res
}
