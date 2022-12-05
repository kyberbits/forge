package forgeutils

func UnreachableError(err error) {
	if err != nil {
		panic(err)
	}
}
