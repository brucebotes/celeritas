package celeritas

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

func (c *Celeritas) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	//pc -> for program caller
	pc, _, _, _ := runtime.Caller(1) // only one level up
	funcObj := runtime.FuncForPC(pc) // get the function which called the pc, we can get the name from it
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	c.InfoLog.Println(fmt.Sprintf("Load Time: %s took %s", name, elapsed))
}
