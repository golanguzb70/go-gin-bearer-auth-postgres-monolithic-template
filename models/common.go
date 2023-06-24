package models

type CheckIfExistsReq struct {
	Table  string
	Column string
	Value  string
}

type CheckIfExistsRes struct {
	Exists bool
}

type UpdateSingleFieldReq struct {
	Id       any
	Table    string
	Column   string
	NewValue any
}
