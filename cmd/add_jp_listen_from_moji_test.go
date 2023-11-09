package cmd

import (
	"fmt"
	"testing"
)

func Test_fetchDataFromMoji(t *testing.T) {
	got, err := fetchWordsFromMoji()
	if err != nil {
		panic(err)
	}
	fmt.Println(got)

}

func Test_fetchTTSFromMoji(t *testing.T) {
	got, err := fetchTTSFromMoji("198990643", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(got)
}

func Test_buildJpListenModel(t *testing.T) {
	res, err := buildJpListenModel()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func Test_deleteFromMojiFolder(t *testing.T) {
	err := deleteFromMojiFolder("nJCCZsM9tL", []string{"198924326"})
	print(err)
}
