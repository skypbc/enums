package include

type EnumItem struct {
	Name  string
	Value any
}

type EnumObjectSettings struct {
	Name    string
	Postfix string
	Prefix  string
}

type EnumFileSettings struct {
	Object EnumObjectSettings
	Path   []string
}

type EnumResult struct {
	Folder   string
	Filename string
	Header   string
	Content  string
}

type EnumFile struct {
	Path  []string
	Names []string
	Items []EnumItem
	Type  string

	Settings EnumFileSettings
	Result   map[string]EnumResult
}
