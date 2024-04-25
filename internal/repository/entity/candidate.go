package entity

type Candidate struct {
	ID    string `bson:"_id,omitempty"`
	Email string `bson:"email"`
	Code  string `bson:"code"`
}

func (e *Candidate) IsExist() bool {
	return e != nil
}
