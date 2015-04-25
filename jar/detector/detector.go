package detector

import (
	"errors"
	"fmt"
	"reflect"
)

type Detector struct {
	// Nothing
}

func Call(funcName string, params []interface{}) (result string, err error) {

	var d Detector

	fv := reflect.ValueOf(d).MethodByName(funcName)

	if !fv.IsValid() {

		err = errors.New(`"` + funcName + `" is not found.`)

	} else {

		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf(`"%v" failed: %v`, funcName, e))
			}
		}()

		ps := make([]reflect.Value, len(params))
		for k, p := range params {
			ps[k] = reflect.ValueOf(p)
		}

		res := fv.Call(ps)

		if len(res) < 2 || !res[0].IsValid() || !res[1].IsValid() {
			err = errors.New(fmt.Sprintf(`Results of "%v": 2 expected, %v given or some invalid.`, funcName, len(res)))
			return
		}

		if r, ok := res[0].Interface().(string); ok {
			result = r
		} else {
			err = errors.New(fmt.Sprintf(`First result of "%v": string expected, %v given.`, funcName, res[0].Kind().String()))
			return
		}

		if e, ok := res[1].Interface().(error); ok {
			err = e
		} else if res[1].Interface() != nil {
			err = errors.New(fmt.Sprintf(`Second result of "%v": error or nil expected, %v given.`, funcName, res[1].Kind().String()))
		}

	}

	return
}
