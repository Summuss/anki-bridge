package common

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func MergeErrors(errList []error) error {
	var res *multierror.Error
	for _, err := range errList {
		res = multierror.Append(res, err)
	}
	return res.ErrorOrNil()
}

func Exec(prog string, arg ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.Command(prog, arg...)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf(
			"exec [%s %s]. %s. stderr:\n%s\n", prog, strings.Join(arg, " "), err.Error(),
			stderr.String(),
		)
	}
	return stdout.String(), nil
}

func CurlGetData(url string) (*[]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download file from %s failed,%s", url, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("download file from %s failed,%s", url, err.Error())
	}

	return &body, nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
