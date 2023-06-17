package wdb

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TanmoySG/TanmoySG-backend/jobs/update-projects/internal/requests"
	wdbr "github.com/TanmoySG/wdb-migrate/pkg/wdb/retro"
)

type WdbAdapter struct {
	httpClient    http.Client
	connectionURL string
}

func NewClient(baseURL string, cluster string, token string) WdbAdapter {
	return WdbAdapter{
		httpClient:    *http.DefaultClient,
		connectionURL: fmt.Sprintf("%s/connect?cluster=%s&token=%s", baseURL, cluster, token),
	}
}

func getError(responseBytes []byte) error {
	var resp map[string]interface{}
	err := json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return nil
	}

	switch resp["status_code"].(string) {
	case "0":
		return fmt.Errorf(resp["response"].(string))
	case "1":
		return nil
	default:
		return nil
	}
}

func (w WdbAdapter) GetData(database, collection string) (*wdbr.GetDataResponse, error) {
	requestBody := wdbr.RequestBody{
		Action: "get-data",
		Payload: wdbr.Payload{
			Database:   database,
			Collection: collection,
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	responseBytes, err := requests.Query(w.httpClient, "", http.MethodPost, w.connectionURL, requestBodyBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(responseBytes))

	err = getError(responseBytes)
	if err != nil {
		return nil, err
	}

	var getDataResponse wdbr.GetDataResponse
	err = json.Unmarshal(responseBytes, &getDataResponse)
	if err != nil {
		return nil, err
	}

	return &getDataResponse, nil
}

func (w WdbAdapter) AddData(database, collection string, data map[string]interface{}) error {
	requestBody := wdbr.RequestBody{
		Action: "add-data",
		Payload: wdbr.Payload{
			Database:   database,
			Collection: collection,
			Data:       (*wdbr.Data)(&data),
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	responseBytes, err := requests.Query(w.httpClient, "", http.MethodPost, w.connectionURL, requestBodyBytes)
	if err != nil {
		return err
	}

	err = getError(responseBytes)
	if err != nil {
		return err
	}

	return nil
}

func (w WdbAdapter) DeleteData(database, collection string, key, value string) error {
	marker := fmt.Sprintf("%s : %s", key, value)
	requestBody := wdbr.RequestBody{
		Action: "delete-data",
		Payload: wdbr.Payload{
			Database:   database,
			Collection: collection,
			Marker:     &marker,
		},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	responseBytes, err := requests.Query(w.httpClient, "", http.MethodPost, w.connectionURL, requestBodyBytes)
	if err != nil {
		return err
	}

	err = getError(responseBytes)
	if err != nil {
		return err
	}

	return nil
}
