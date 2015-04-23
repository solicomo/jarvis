// +build !windows

package detector

import (
	"strings"
	"syscall"
	"unsafe"
)

func int8ToString(bs []int8) string {
	return strings.TrimRight(string(*(*[]byte)(unsafe.Pointer(&bs))), "\x00")
}

func (Detector) OSVer(params ...string) (result string, err error) {
	var u syscall.Utsname
	if err = syscall.Uname(&u); err != nil {
		return
	}

	sysName := int8ToString(u.Sysname[:])
	nodName := int8ToString(u.Nodename[:])
	release := int8ToString(u.Release[:])
	version := int8ToString(u.Version[:])
	machine := int8ToString(u.Machine[:])
	domName := int8ToString(u.Domainname[:])

	if len(params) == 0 {
		result = sysName + " " + nodName + " " + release + " " + version + " " + machine + " " + domName
	}

	for _, p := range params {
		switch p {
		case "s":
			result += sysName + " "
		case "n":
			result += nodName + " "
		case "r":
			result += release + " "
		case "v":
			result += version + " "
		case "m":
			result += machine + " "
		case "o":
			result += domName + " "
		}
	}
	return
}
