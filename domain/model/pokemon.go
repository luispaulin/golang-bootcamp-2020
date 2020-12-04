package model

// Pokemon model
// TODO: Can I use csv notations in model?
type Pokemon struct {
	Name string `csv:"name"`
	Url  string `csv:"url"`
}
