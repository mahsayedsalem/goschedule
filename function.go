package goschedule

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type function struct {
	implementation      interface{}
	params              []interface{}
	latestExecutionTime time.Duration
	identifier          string
	res                 []reflect.Value
}

type functionInfo struct {
	LatestExecutionTime time.Duration
	Results             []interface{}
}

func newFunction(implementation interface{}, params []interface{}, identifier string) *function {
	return &function{
		implementation: implementation,
		params:         params,
		identifier:     identifier,
	}
}

func (f *function) runFunc(wg *sync.WaitGroup) {
	log.Info("Running function inside job ", f.identifier)
	defer wg.Done()
	now := time.Now()
	fValue := reflect.ValueOf(f.implementation)
	if len(f.params) != fValue.Type().NumIn() {
		fmt.Println(fValue.Type().NumIn(), len(f.params))
		return
	}
	var params []reflect.Value
	for _, param := range f.params {
		params = append(params, reflect.ValueOf(param))
	}
	f.res = fValue.Call(params)
	f.latestExecutionTime = time.Since(now)
	log.Info("Function inside job ", f.identifier, " execution time is ", f.latestExecutionTime, " results are ", f.res)
}

func (f *function) getLatestExecutionTime() time.Duration {
	return f.latestExecutionTime
}

func (f *function) getLatestResults() []interface{} {
	var results []interface{}
	for _, res := range f.res {
		results = append(results, res.Interface())
	}
	return results
}

func (f *function) GetFuncInfo() *functionInfo {
	return &functionInfo{
		LatestExecutionTime: f.getLatestExecutionTime(),
		Results:             f.getLatestResults(),
	}
}
