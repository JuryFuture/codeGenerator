package entity

type ClassInfo struct {
	PackageName  string
	ClassName    string
	DateTime     string
	Year         string
	TableComment string
	Author       string
	Date         string
	TableName    string
	Fields       []FieldInfo
	Methods      []MethodInfo
}
