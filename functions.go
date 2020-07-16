package configtypes

// StringOrElse is a shortcut for calling String() on a value if it is defined, or returning an
// alternative string if it is empty.
//
//     StringOrElse(NewOptBool(true), "undefined") // == "true"
//     StringOrElse(OptBool{}, "undefined")        // == "undefined"
func StringOrElse(value SingleValue, orElseString string) string {
	if value.IsDefined() {
		return value.String()
	}
	return orElseString
}
