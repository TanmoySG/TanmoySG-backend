package requests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func Query(client http.Client, authToken, requestMethod, requestUrl string, requestBodyBytes []byte) ([]byte, error) {
	request, err := http.NewRequest(requestMethod, requestUrl, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("bearer %s", authToken))
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBodyBytes, nil
}
