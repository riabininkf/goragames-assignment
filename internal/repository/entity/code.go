package entity

type Code struct {
	ID     string `bson:"_id,omitempty"`
	Name   string `bson:"name"`
	Usages uint64 `bson:"usages"`
}

func (e *Code) IsExist() bool {
	return e != nil
}
