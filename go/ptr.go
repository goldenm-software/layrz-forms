package layrz

// Ptr returns a pointer to v. Form fields are pointers so that nil can model
// an absent value distinctly from a present zero value.
func Ptr[T any](v T) *T {
	return &v
}
