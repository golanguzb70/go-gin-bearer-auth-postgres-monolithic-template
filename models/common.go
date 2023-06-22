package models

type CheckIfExistsReq struct {
	Table  string
	Column string
	Value  string
}

type CheckIfExistsRes struct {
	Exists bool
}
