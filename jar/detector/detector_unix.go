// +build !windows

package detector

func (Detector) OSVer(params []string) (result string, err error) {
	var u syscall.Utsname
	if err = syscall.Uname(&u); err != nil {
		return
	}

	sysName := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])
	nodName := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])
	release := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])
	version := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])
	machine := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])
	domName := string(u.Sysname[:bytes.IndexByte(u.Sysname, 0)])

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
