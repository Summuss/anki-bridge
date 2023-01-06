package anki

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/model"
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
	_ = `
{
    "action": "addNote",
    "version": 6,
    "params": {
        "note": {
            "deckName": "%s",
            "modelName": "%s",
            "fields": {
                "Front": "front content",
                "Back": "back content"
            },
            "options": {
                "allowDuplicate":true 
            }
        }
    }
}
`
	var fields = map[string]interface{}{
		"front":      card.Front,
		"back":       card.Back,
		"db_id":      card.ModelID,
		"collection": card.Collection,
	}
	var options = map[string]interface{}{"allowDuplicate": true}
	var note = map[string]interface{}{
		"deckName":  card.Desk,
		"modelName": "BasicTwoSide",
		"fields":    fields,
		"options":   options,
	}

	var params = map[string]interface{}{"note": note}

	res, err := requestAnki("addNote", params)
	if err != nil {
		return err
	}
	id := res["result"].(float64)
	card.ID = int64(id)
	return nil
}

func StoreMedia(resource *model.Resource) error {

	sEnc := b64.StdEncoding.EncodeToString(resource.GetData())
	var params = map[string]interface{}{
		"filename": resource.Metadata.FileName,
		"data":     sEnc,
	}
	_, err := requestAnki("storeMediaFile", params)
	return err
}

func GetAllDecks() ([]string, error) {
	res, err := requestAnki("deckNames", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("fetch desk names failed, %s", err.Error())
	}
	datas := res["result"].([]interface{})
	deskNames := lo.Map(
		datas, func(item interface{}, _ int) string {
			return item.(string)
		},
	)
	return deskNames, nil
}

var ankiRequestCtrlCh = make(chan interface{}, 1)

func requestAnki(action string, params map[string]interface{}) (map[string]interface{}, error) {
	req := map[string]interface{}{"action": action, "params": params, "version": 6}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	ankiRequestCtrlCh <- struct{}{}
	resp, err := http.Post(
		config.Conf.AnkiAPIURL, "application/json", bytes.NewBuffer(jsonStr),
	)
	<-ankiRequestCtrlCh
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requestAnki: %s", resp.Status)
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("requestAnki: %s", err.Error())
	}
	errInfo := res["error"]
	if errInfo != nil {
		errStr := errInfo.(string)
		if len(errStr) > 0 {
			return nil, fmt.Errorf("requestAnki: %s,\n%s", errStr, string(jsonStr))
		}
	}
	return res, nil

}
