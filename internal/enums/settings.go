package enums

type LangSettings struct {
	// Настройки файла
	File struct {
		// Заголовок который будет добавлен в начало файла
		Header string `json:"header"`
		// Содержимое файла
		Content struct {
			// Шаблон для содержимого файла
			Tmpl string `json:"tmpl"`
			// Разделитель между Content.Tmpl
			Sep string `json:"sep"`
		}
		// Расширение файла
		Extension string `json:"extension"`
		// Настройки для имени файла
		Name struct {
			// Указывает включать в имя часть пути к файлу
			// Пример:
			// - имеем следующий путь к файлу: /path/to/my/file
			// - если PathCount = 2, то имя файла бдует to_my_file
			PathCount int `json:"path_count"`
			// Переопределяет первую версию имени файла
			Value1 string `json:"value1"`
			// Переопределяет вторую версию имени файла
			Value2 string `json:"value2"`
			// Разделитель между частями имени файла для первой версии
			Sep1 string `json:"sep1"`
			// Разделитель между частями имени файла для второй версии
			Sep2 string `json:"sep2"`
			// Префикс для первой версии имени файла
			Prefix1 string `json:"prefix1"`
			// Префикс для второй версии имени файла
			Prefix2 string `json:"prefix2"`
		}
	} `json:"file"`

	// Настройки папки
	Folder struct {
		Name struct {
			// Указывает включать в имя часть пути
			// Пример:
			// - имеем следующий путь к файлу: /path/to/my/folder
			// - если PathCount = 2, то имя папки бдует path/to_my_folder
			PathCount int `json:"path_count"`
			// Переопределяет имя папки
			Value string `json:"value"`
			// Разделитель между частями имени папки
			Sep string `json:"sep"`
			// Префикс для имени папки
			Prefix string `json:"prefix"`
			// Включение в имя папки имени файла
			Append struct {
				// Включает в имя папки первую версию имени файла
				Filename1 bool `json:"filename1"`
				// Включает в имя папки вторую версию имени файла
				Filename2 bool `json:"filename2"`
			}
		}
		// Указывает сохранять файлы в плоской структуре (без вложенных папок)
		Flat bool `json:"flat"`
	}

	// Настройка для элементов
	Item struct {
		// Шаблон для элемента
		Tmpl string `json:"tmpl"`
		// Разделитель между элементами
		Sep string `json:"sep"`
		// Настройки для имени элемента
		Name struct {
			// Разделитель между частями имени элемента
			Sep string `json:"sep"`
			// Стиль имени "pascal" или "snake"
			Style string `json:"style"`
			// Использовать ли верхний регистр
			// Пример: "my_name" -> "MY_NAME"
			Upper      bool `json:"upper"`
			Lower      bool `json:"lower"`
			Capitalize bool `json:"capitalize"`
			// Если true, то имя элемента будет включать в себя имя объекта
			PrependObjectName struct {
				Allways bool `json:"always"`
				// Включать имя объекта если есть префикс
				HasPrefix  bool `json:"has_prefix"`
				HasPostfix bool `json:"has_postfix"`
			} `json:"prepend_object_name"`
		} `json:"name"`
	} `json:"item"`

	// Настройки для объекта
	Object struct {
		Type map[string]any `json:"type"`
		Name struct {
			// Указывает включать в имя часть пути к файлу
			// Пример:
			// - имя объекта: Foo
			// - имеем следующий путь к файлу: /path/to/my/file
			// - если PathCount = 2, то имя файла бдует MyFileFoo
			PathCount int `json:"path_count"`
			// Добавляет постфикс к имени объекта
			// Пример:
			// - имя объекта: Foo
			// - если Postfix = "Bar", то имя объекта будет FooBar
			Postfix string `json:"postfix"`
		}
	} `json:"object"`
}
