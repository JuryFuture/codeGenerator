package entity

type DocTableInfo struct {
	TableName    string
	ColumnNum    int
	TableComment string
	Columns      []DocColumnInfo
}
