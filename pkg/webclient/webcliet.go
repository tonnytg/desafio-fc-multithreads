package webclient

import (
	"context"
	"io"
	"net/http"
)

func Request(ctx context.Context, method string, url string, data io.Reader) ([]byte, error) {

	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return []byte{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	return body, nil
}
