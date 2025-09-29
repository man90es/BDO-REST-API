package utils

import "net/http"

func SendDummyRequest(url string) error {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
