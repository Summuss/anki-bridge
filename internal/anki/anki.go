package anki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BasicTwoSide struct {
}

type Card struct {
	ID         int64
	Collection string
	ModelID    string

	// need to manually set
	Front string
	Back  string
	Desk  string
}

// AddCard write id back to card after successful adding
func AddCard(card *Card) error {
	return nil
}

func requestAnki(action string, params map[string]interface{}) (map[string]interface{}, error) {
	req := map[string]interface{}{"action": action, "params": params, "version": 6}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	resp, err := http.Post(
		"http://192.168.162.1:8766", "application/json", bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	return res, nil

}
