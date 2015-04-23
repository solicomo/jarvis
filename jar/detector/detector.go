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

	} else if fv.Type().NumIn() != len(params) {

		err = errors.New(fmt.Sprintf(`Params for "%v": %v expected, %v given.`,
			funcName, fv.Type().NumIn(), len(params)))

	} else {

		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v failed: %v", funcName, e))
			}
		}()

		ps := make([]reflect.Value, len(params))
		for k, p := range params {
			ps[k] = reflect.ValueOf(p)
		}

		res := fv.Call(ps)

		result = res[0].Interface().(string)
		err = res[1].Interface().(error)
	}

	return
}

func (Detector) Uptime() {

}

func (Detector) Load() {

}

func (Detector) CPUName() {

}

func (Detector) CPUCore() {

}

func (Detector) CPURate() {

}

func (Detector) MemSize() {

}

func (Detector) MemRate() {

}

func (Detector) SwapRate() {

}

func (Detector) DiskSize() {

}

func (Detector) DiskRead() {

}

func (Detector) DiskWrite() {

}

func (Detector) NetRead() {

}

func (Detector) NetWrite() {

}
