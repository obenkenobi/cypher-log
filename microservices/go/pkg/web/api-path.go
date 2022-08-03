package web

import "fmt"

// apiVersionFormat formats the api version of sub-path of an api endpoint such that if
// the version number is 11, "v11" is returned
func apiVersionFormat(versionNum int) string { return fmt.Sprintf("v%v", versionNum) }

// APIPath formats the url sub-path in the form "/v{versionNum}/{path}"
func APIPath(versionNum int, path string) string {
	return fmt.Sprintf("/%v/%v", apiVersionFormat(versionNum), path)
}
