package index

import (
	"fmt"
	"strings"
)

// resolves a single string that has a reference in it
func resolveString(valString string, depthLeft int) interface{} {
	key := strings.Replace(valString, "REF::", "", 1)
	file, ok := I.Lookup(key)

	// if key is found, get content
	if ok {
		// change bytes into a map
		jsonMap, err := file.ToMap()
		if err != nil {
			errMessage := fmt.Sprintf("REF::ERR key '%s' cannot be parsed into json: %s", key, err.Error())
			return errMessage
		}

		return ResolveReferences(jsonMap, depthLeft-1)
	}

	// if key is not found
	return fmt.Sprintf("REF::ERR key '%s' not found", key)
}

func ResolveReferences(jsonVal interface{}, depthLen int) interface{} {

}
