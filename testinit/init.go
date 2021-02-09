// Import this file in tests like this:
// 		import _ "github.com/BionicTeam/bionic/testinit"
// to switch directory to project root and navigate testdata easier.
// Solution source: https://brandur.org/fragments/testing-go-project-root

package testinit

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
