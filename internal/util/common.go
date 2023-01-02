package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func MergeErrors(errList []error) error {
	errList = lo.Filter(
		errList, func(item error, _ int) bool {
			return item != nil
		},
	)
	if len(errList) > 0 {
		msg := lo.Reduce(
			errList, func(agg string, err error, i int) string {
				return agg + "\n" + err.Error()
			}, "",
		)
		return errors.New(msg)
	}
	return nil
}

func Exec(prog string, arg ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.Command(prog, arg...)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("%s. stderr:\n%s\n", err.Error(), stderr.String())
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

	tf, err := os.CreateTemp("", "tts.*.mp3")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tf.Name())
	//curl '{url}' --output '{tmp_file}

	log.Printf("download file from %s\n", url)
	_, err = Exec("curl", url, "--output", tf.Name())
	if err != nil {
		return nil, fmt.Errorf("download file from %s failed,%s", url, err.Error())
	}
	file, err := os.Open(tf.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return nil, err
	}
	return &bs, nil

}
