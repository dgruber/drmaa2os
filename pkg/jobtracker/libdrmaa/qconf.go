package libdrmaa

import (
	"fmt"
	"os/exec"
	"strings"
)

func QconfSQL() ([]string, error) {
	out, err := Qconf([]string{"-sql"})
	if err != nil {
		return nil, err
	}
	return ParseLines(string(out)), nil
}

func Qconf(args []string) (string, error) {
	out, err := exec.Command("qconf", args...).Output()
	if err != nil {
		return "", fmt.Errorf("error executing qconf with args %v: %v", args, err)
	}
	return string(out), nil
}

func ParseLines(out string) []string {
	res := make([]string, 0, 1)
	for _, line := range strings.Split(out, "\n") {
		res = append(res, line)
	}
	return res
}
