package ptr

/*
This is a helper package to have a pointer to the primitive type fast and not ugly.
In go you can't have a pointer to the const what a literal is (e.g. &"hello" won't work).
*/

func String(v string) *string {
	return &v
}
