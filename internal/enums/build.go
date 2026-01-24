package enums

import (
	"fmt"
	"github.com/skypbc/enums/internal/utils"
	"github.com/skypbc/goutils/gerrors"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func Build(enumFiles []EnumFile, settings map[string]any) (map[string]string, error) {
	files := make(map[string]string)
	for _, enumFile := range enumFiles {
		if err := build(&enumFile, settings); err != nil {
			return nil, err
		}
		for _, result := range enumFile.Result {
			fullpath := filepath.Join(result.Folder, result.Filename)
			content, ok := files[fullpath]
			if !ok {
				content = result.Header
			}
			files[fullpath] = content + result.Content
		}
	}
	return files, nil
}

func build(enumFile *EnumFile, settings map[string]any) (err error) {
	langs, ok := settings["_"].(map[string]any)
	if !ok {
		return gerrors.NewIncorrectParamsError().
			SetTemplate("Settings not found")
	}

	enumSettingsRaw := utils.GetMap(settings, strings.Join(enumFile.Settings.Path, "."))
	enumFile.Result = map[string]EnumResult{}

	for lang, langSettingsRaw := range langs {
		enumSettings, err := getEnumSettings(enumSettingsRaw[lang], langSettingsRaw)
		if err != nil {
			return err
		}

		objectName := getObjectName(enumFile, &enumSettings)
		objectTypeName := getObjectTypeName(enumFile, &enumSettings)
		lines := getLines(enumFile, objectName, &enumSettings)
		filename1 := getFilename1(enumFile, &enumSettings)
		filename2 := getFilename2(enumFile, &enumSettings)
		folder := getFolder(enumFile, lang, filename1, filename2, &enumSettings)

		log.Debug().Msgf(
			"Lang: %s, Folder: %s, Filename1: %s, Filename2: %s, Extension: %s\n",
			lang, folder, filename1, filename2, enumSettings.File.Extension,
		)

		enumFile.Result[lang] = EnumResult{
			Folder:   folder,
			Filename: getEnumFilename(filename1, &enumSettings),
			Header:   getHeader(filename2, objectTypeName, &enumSettings),
			Content:  getContent(objectName, objectTypeName, lines, &enumSettings),
		}
	}

	return nil
}

func normalizeNames(name []string) []string {
	var res []string
	for _, n := range name {
		res = append(res, strings.Split(n, "_")...)
	}
	return res
}

func getLangSettings(settings any) (res LangSettings, err error) {
	res.Folder.Name.Sep = "_"
	if err = utils.Map2Struct(settings, &res); err != nil {
		return res, err
	}
	return res, nil
}

func getEnumSettings(enumSettingsRaw any, langSettingsRaw any) (res LangSettings, err error) {
	if res, err = getLangSettings(langSettingsRaw); err != nil {
		return res, err
	}
	if err = utils.Map2Struct(enumSettingsRaw, &res); err != nil {
		return res, err
	}
	return res, nil
}

func getObjectName(enumFile *EnumFile, enumSettings *LangSettings) []string {
	objName := enumFile.Settings.Object.Name
	// Объект не был включен в enumFile при парсинге
	if len(objName) == 0 {
		return nil
	}

	nameParts := []string{}
	// При парсинге enumFile был создан префикс, который нужно добавить к имени
	if prefix := enumFile.Settings.Object.Prefix; prefix != "" {
		nameParts = append(nameParts, prefix)
	}
	nameParts = append(nameParts, objName)

	// В настройках был определен постфикс
	if postfix := enumSettings.Object.Name.Postfix; len(postfix) > 0 {
		nameParts = append(nameParts, postfix)
	}

	// При парсинге enumFile был создан постфикс, который нужно добавить к имени
	if postfix := enumFile.Settings.Object.Postfix; postfix != "" {
		nameParts = append(nameParts, postfix)
	}

	pathCound := 0
	if enumSettings.Object.Name.PathCount > -1 {
		// В настройках было определено включение пути в имя
		pathCound = enumSettings.Object.Name.PathCount
	}

	// Если указано включить элементы из пути в имя
	if pathCound > 0 {
		start := len(enumFile.Path) - pathCound
		if start < 0 {
			start = 0
		}
		if len(enumFile.Path) > 0 {
			nameParts = append(slices.Clone(enumFile.Path[start:]), nameParts...)
		}
	}

	return nameParts
}

func getLines(enumFile *EnumFile, objectName []string, enumSettings *LangSettings) []string {
	var lines []string
	for _, item := range enumFile.Items {
		if item.Name == "_" {
			continue
		}

		var name []string
		var value string

		name = append(name, item.Name)

		if enumFile.Type == "string" {
			value = fmt.Sprintf("\"%s\"", item.Value)
		} else {
			value = fmt.Sprintf("%v", item.Value)
		}

		if enumSettings.Item.Name.PrependObjectName.Allways {
			if len(objectName) > 0 {
				name = append(slices.Clone(objectName), name...)
			}
		} else {
			if len(enumFile.Settings.Object.Prefix) > 0 && enumSettings.Item.Name.PrependObjectName.HasPrefix {
				if len(objectName) > 0 {
					name = append(slices.Clone(objectName), name...)
				}
			} else if len(enumFile.Settings.Object.Postfix) > 0 && enumSettings.Item.Name.PrependObjectName.HasPostfix {
				if len(objectName) > 0 {
					name = append(slices.Clone(objectName), name...)
				}
			} else {
				if len(enumFile.Settings.Object.Prefix) > 0 {
					name = append([]string{enumFile.Settings.Object.Prefix}, name...)
				} else if len(enumFile.Settings.Object.Postfix) > 0 {
					name = append([]string{enumFile.Settings.Object.Postfix}, name...)
				}
			}
		}

		var textName string
		if enumSettings.Item.Name.Style == "pascal" {
			textName = utils.JoinPascal(name...)
		} else {
			textName = strings.Join(name, enumSettings.Item.Name.Sep)
			textName = strings.ToLower(textName)
		}

		if enumSettings.Item.Name.Upper {
			textName = strings.ToUpper(textName)
		} else if enumSettings.Item.Name.Lower {
			textName = strings.ToLower(textName)
		} else if enumSettings.Item.Name.Capitalize {
			textName = strings.ToUpper(textName[:1]) + textName[1:]
		}

		line := strings.ReplaceAll(enumSettings.Item.Tmpl, "{name}", textName)
		line = strings.ReplaceAll(line, "{value}", value)

		lines = append(lines, line)
	}
	return lines
}

func getFilename1(enumFile *EnumFile, enumSettings *LangSettings) string {
	nameParts := []string{}
	if len(enumSettings.File.Name.Prefix1) > 0 {
		nameParts = append(nameParts, enumSettings.File.Name.Prefix1)
	}

	if len(enumSettings.File.Name.Value1) > 0 {
		nameParts = append(nameParts, enumSettings.File.Name.Value1)
	} else {
		pathCount := enumSettings.File.Name.PathCount
		if enumSettings.File.Name.PathCount > -1 {
			pathCount = enumSettings.File.Name.PathCount
		}
		if len(enumFile.Path) > 0 && pathCount > 0 {
			start := len(enumFile.Path) - pathCount
			if start < 0 {
				start = 0
			}
			if len(enumFile.Path) > 0 {
				nameParts = append(nameParts, enumFile.Path[start:]...)
			}
		}
		nameParts = append(nameParts, enumFile.Names...)
	}
	return strings.Join(normalizeNames(nameParts), enumSettings.File.Name.Sep1)
}

func getFilename2(enumFile *EnumFile, enumSettings *LangSettings) string {
	nameParts := []string{}
	if len(enumSettings.File.Name.Prefix2) > 0 {
		nameParts = append(nameParts, enumSettings.File.Name.Prefix2)
	}
	if len(enumSettings.File.Name.Value2) > 0 {
		nameParts = append(nameParts, enumSettings.File.Name.Value2)
	} else {
		nameParts = append(nameParts, enumFile.Names...)
	}
	return strings.Join(normalizeNames(nameParts), enumSettings.File.Name.Sep2)
}

func getFolder(enumFile *EnumFile, lang string, filename1 string, filename2 string, enumSettings *LangSettings) string {
	nameParts := []string{lang}
	if !enumSettings.Folder.Flat {
		nameParts = append(nameParts, slices.Clone(enumFile.Path)...)
	}

	// Если задано, переопределяет имя папки
	if enumSettings.Folder.Name.Value != "" {
		nameParts = append(nameParts, enumSettings.Folder.Name.Value)
	} else {
		// Иначе, в качестве имени будет использоваться имя файла (если указано)
		if enumSettings.Folder.Name.Append.Filename1 {
			nameParts = append(nameParts, filename1)
		} else if enumSettings.Folder.Name.Append.Filename2 {
			nameParts = append(nameParts, filename2)
		}
	}

	// Добавляет к имени папки имя пути
	if enumSettings.Folder.Name.PathCount > 0 {
		start := len(nameParts) - enumSettings.Folder.Name.PathCount
		if start < 0 {
			start = 0
		}
		paths := []string{}
		if len(nameParts) > 0 {
			paths = append(paths, nameParts[start:]...)
		}
		nameParts = nameParts[:start]
		nameParts = append(nameParts, strings.Join(paths, enumSettings.Folder.Name.Sep))
	}

	// Добавление префикса к имени папки
	if enumSettings.Folder.Name.Prefix != "" {
		end := len(nameParts)
		// Если в качестве имени папки используется имя файла (filename2), то не добавляем к нему префикс, если нужно,
		// его можно настроить отдельно, в настрйоках имени файла
		if enumSettings.Folder.Name.Append.Filename2 {
			end--
		}
		// Первый элемент - это язык, к нему префикс не добавляется
		for i := 1; i < end; i++ {
			nameParts[i] = enumSettings.Folder.Name.Prefix + nameParts[i]
		}
	}

	return strings.Join(nameParts, "/")
}

func getObjectTypeName(enumFile *EnumFile, enumSettings *LangSettings) string {
	objectTypeName := enumFile.Type
	if enumSettings.Object.Type != nil {
		if typeName, ok := enumSettings.Object.Type[enumFile.Type]; ok {
			objectTypeName = fmt.Sprintf("%v", typeName)
		}
	}
	return objectTypeName
}

func getEnumFilename(filename1 string, enumSettings *LangSettings) string {
	return filename1 + enumSettings.File.Extension
}

func getHeader(filename2 string, objectTypeName string, enumSettings *LangSettings) string {
	resultHeader := strings.ReplaceAll(enumSettings.File.Header, "{filename_2}", filename2)
	return strings.ReplaceAll(resultHeader, "{object_type}", objectTypeName)
}

func getContent(objectName []string, objectTypeName string, lines []string, enumSettings *LangSettings) string {
	content := enumSettings.File.Content.Tmpl
	content = strings.ReplaceAll(content, "{object_name}", utils.JoinPascal(objectName...))
	content = strings.ReplaceAll(content, "{object_type}", objectTypeName)
	content = strings.ReplaceAll(content, "{items}", strings.Join(lines, enumSettings.Item.Sep))
	return enumSettings.File.Content.Sep + content
}
