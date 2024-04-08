package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/summuss/anki-bridge/internal/anki"
	"github.com/summuss/anki-bridge/internal/common"
	"github.com/summuss/anki-bridge/internal/config"
	"github.com/summuss/anki-bridge/internal/dto"
	"github.com/summuss/anki-bridge/internal/model"
	"net/http"
	"time"
)

func init() {
	rootCmd.AddCommand(addJPListenFromMojiCMD)
}

var addJPListenFromMojiCMD = &cobra.Command{
	Use:  "add_jp_listen_from_moji",
	Args: cobra.MaximumNArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := anki.GetAllDecks()
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		jlms, err := buildJpListenModel()
		if err != nil {
			return err
		}
		ms := lo.Map(
			*jlms, func(item *model.JpListen[dto.MojiBookItem], _ int) model.IModel {
				return item
			},
		)
		err = addModels(&ms)
		if err != nil {
			return err
		}
		tarIds := lo.Map(
			*jlms, func(item *model.JpListen[dto.MojiBookItem], _ int) string {
				return item.ExtInfo.ObjectId
			},
		)
		err = deleteFromMojiFolder("dM9FRXFvAG", tarIds)
		return err
	},
}

func buildJpListenModel() (*[]*model.JpListen[dto.MojiBookItem], error) {
	wordItems, err := fetchWordsFromMoji()
	if err != nil {
		return nil, err
	}
	noteInfo := config.Conf.GetNoteInfoByName(common.NoteType_JpListenMoji_Name)
	ms := lo.Map(
		*wordItems, func(word *dto.MojiBookItem, index int) *model.JpListen[dto.MojiBookItem] {
			m := &model.JpListen[dto.MojiBookItem]{
				JpListenKey: word.Spell,
				ExtInfo:     word,
			}
			m.SetNoteInfo(noteInfo)
			m.SetNoteTypeName(noteInfo.Name)
			return m
		},
	)

	for _, t := range ms {
		data, err := fetchTTSFromMoji((*t).ExtInfo.ObjectId, 0)
		if err != nil {
			return nil, err
		}
		resource := model.Resource{
			Metadata: model.ResourceMetadata{
				FileName: (*t).ExtInfo.ObjectId + ".mp3", ResourceType: model.Sound,
				ExtName: ".mp3",
			},
		}
		resource.SetData(*data)
		resources := []model.Resource{resource}
		(*t).SetResources(&resources)
	}
	return &ms, nil

}

func fetchTTSFromMoji(objectId string, retryTimes int) (*[]byte, error) {
	req := map[string]interface{}{
		"tarId":           objectId,
		"tarType":         102,
		"voiceId":         "f002",
		"_ApplicationId":  "E62VyFVLMiW7kvbtVq3p",
		"_ClientVersion":  "js3.4.1",
		"_InstallationId": config.Conf.MojiInstallationId,
		"_SessionToken":   config.Conf.MojiSessionToken,
	}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	resp, err := http.Post(
		"https://api.mojidict.com/parse/functions/tts-fetch",
		"application/json", bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		if retryTimes > 5 {
			return nil, fmt.Errorf("requestMoji: %s", resp.Status)
		} else {
			time.Sleep(100 * time.Millisecond)
			return fetchTTSFromMoji(objectId, retryTimes+1)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requestMoji: %s", resp.Status)
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	ttlUrl := res["result"].(map[string]interface{})["result"].(map[string]interface{})["url"].(string)
	data, err := common.CurlGetData(ttlUrl)
	if err != nil {
		return nil, fmt.Errorf("download tts from %s failed,error:\n%s", ttlUrl, err.Error())
	}
	return data, nil
}

func fetchWordsFromMoji() (*[]*dto.MojiBookItem, error) {
	req := map[string]interface{}{
		"count":           50,
		"fid":             "dM9FRXFvAG",
		"pageIndex":       1,
		"sortType":        0,
		"_ApplicationId":  "E62VyFVLMiW7kvbtVq3p",
		"_ClientVersion":  "js3.4.1",
		"_InstallationId": config.Conf.MojiInstallationId,
		"_SessionToken":   config.Conf.MojiSessionToken,
	}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	resp, err := http.Post(
		"https://api.mojidict.com/parse/functions/folder-fetchContentWithRelatives",
		"application/json", bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requestMoji: %s", resp.Status)
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("requestMoji: %s", err.Error())
	}
	items := res["result"].(map[string]interface{})["result"].([]interface{})
	result := lo.Map(
		items, func(e interface{}, i int) *dto.MojiBookItem {
			item := e.(map[string]interface{})
			target := item["target"].(map[string]interface{})
			accent, ok := target["accent"]
			if !ok {
				accent = ""
			}
			return &dto.MojiBookItem{
				ObjectId: target["objectId"].(string),
				Spell:    target["spell"].(string),
				Excerpt:  target["excerpt"].(string),
				Pron:     target["pron"].(string),
				Accent:   accent.(string),
			}

		},
	)
	return &result, nil
}

func deleteFromMojiFolder(fid string, tarIds []string) error {
	req := map[string]interface{}{
		"functions": []map[string]interface{}{
			{
				"name": "deleteItems",
				"params": map[string]interface{}{
					"pfid":   fid,
					"tarIds": tarIds,
				},
			},
		},
		"_ApplicationId":  "E62VyFVLMiW7kvbtVq3p",
		"_ClientVersion":  "js3.4.1",
		"_InstallationId": config.Conf.MojiInstallationId,
		"_SessionToken":   config.Conf.MojiSessionToken,
	}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("requestMoji: %s", err.Error())
	}
	resp, err := http.Post(
		"https://api.mojidict.com/parse/functions/union-api",
		"application/json", bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return fmt.Errorf("requestMoji: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("requestMoji: %s", resp.Status)
	}
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return fmt.Errorf("requestMoji: %s", err.Error())
	}
	code := res["result"].(map[string]interface{})["code"].(float64)
	if code != 200 {
		return fmt.Errorf("requestMoji: delete item failed, resultCode:%v", code)
	}
	return nil
}
