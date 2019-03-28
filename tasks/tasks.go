package tasks

func TestTask(args ...string) (string, error) {
	ret := ""
	for _, arg := range args {
		ret += arg
	}
	return ret, nil
}
