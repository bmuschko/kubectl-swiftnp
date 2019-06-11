package renderer

func BooleanIcon(flag bool) string {
	if flag {
		return "✔"
	}

	return "✖"
}
