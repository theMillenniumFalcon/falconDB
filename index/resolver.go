package index

import (
	"fmt"
	"reflect"
	"strings"
)

// resolves a single string that has a reference in it
func ResolveString(valString string, depthLeft int) interface{} {
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

// ResolveReferences tries to find key references and
// if found, replace the references with their corresponding value
func ResolveReferences(jsonVal interface{}, depthLeft int) interface{} {
	// if max recursive depth is exceeded, return as is
	if depthLeft < 1 {
		return jsonVal
	}

	val := reflect.ValueOf(jsonVal)

	switch val.Kind() {
	case reflect.String:
		valString := val.String()

		// if value is reference to another key
		if strings.Contains(valString, "REF::") {
			resolvedString := ResolveString(valString, depthLeft)
			return resolvedString
		}

		// if not a ref, keep as it is
		return valString

	case reflect.Slice:
		numberOfValues := val.Len()
		newSlice := make([]interface{}, numberOfValues)

		// for each value in the slice, try to resolve it recursively
		for i := 0; i < numberOfValues; i++ {
			pointer := val.Index(i)
			newSlice[i] = ResolveReferences(pointer.Interface(), depthLeft)
		}

		return newSlice

	case reflect.Map:
		newMap := make(map[string]interface{})

		// for each value in the map, try to resolve it recursively
		for _, key := range val.MapKeys() {
			nestedVal := val.MapIndex(key).Interface()
			newMap[key.String()] = ResolveReferences(nestedVal, depthLeft)
		}

		return newMap

	default:
		return jsonVal
	}
}
