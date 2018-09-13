package rewriter

var slash = map[bool]string{
	false: "",
	true:  "/",
}

func join(dir string, file string, trailingSlash bool) string {
	output := dir
	if len(output) == 0 {
		output = file
	} else if len(file) != 0 {
		output += "/" + file
	}
	if len(output) == 0 {
		return slash[trailingSlash]
	}
	return output + slash[trailingSlash]
}
