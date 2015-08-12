package data

// scanMap does basically what is described in the Map.Get documentation.
// The value found at p is written to v.
func scanMap(m Map, p string, v *Value) (err error) {
	path, err := CompilePath(p)
	if err != nil {
		return err
	}
	val, err := path.Evaluate(m)
	if err != nil {
		return err
	}
	*v = val
	return nil
}

// setInMap does basically what is described in the Map.Set documentation.
// The value at v is written to m at the given path.
func setInMap(m Map, p string, v Value) (err error) {
	path, err := CompilePath(p)
	if err != nil {
		return err
	}
	return path.Set(m, v)
}
