package data

func Must(result interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return result
}
