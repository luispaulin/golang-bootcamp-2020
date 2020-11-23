package model

// Pokemon model
// TODO: Can I use csv notations in model?
type Pokemon struct {
	ID    uint16 `csv:"id"`
	Name  string `csv:"name"`
	Type1 string `csv:"type1"`
	Type2 string `csv:"type2"`
}
