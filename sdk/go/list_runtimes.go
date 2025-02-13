package pesto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Runtime struct {
	Language string   `json:"language"`
	Version  string   `json:"version"`
	Aliases  []string `json:"aliases"`
	Compiled bool     `json:"compiled"`
}

type RuntimeResponse struct {
	Runtime []Runtime `json:"runtime"`
}

// ListRuntimes calls the list-runtimes endpoint. The Language and Version item from the response struct
// can be used to create an execute code request.
func (c *Client) ListRuntimes(ctx context.Context) (RuntimeResponse, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL.JoinPath("/api/list-runtimes").String(), nil)
	if err != nil {
		return RuntimeResponse{}, fmt.Errorf("creating request: %w", err)
	}

	response, err := c.sendRequest(ctx, request)
	if err != nil {
		return RuntimeResponse{}, fmt.Errorf("sending request: %w", err)
	}

	if response.StatusCode != 200 {
		var errResponse errorResponse
		// HACK: the error is intentionally not handled, we wanted to leave the empty errorResponse struct
		// if there is any non-json response being sent from the server
		json.NewDecoder(response.Body).Decode(&errResponse)

		err = response.Body.Close()
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) && !errors.Is(err, http.ErrBodyReadAfterClose) {
			return RuntimeResponse{}, fmt.Errorf("closing response body: %w", err)
		}

		return RuntimeResponse{}, c.handleErrorCode(response.StatusCode, errResponse)
	}

	var runtimes RuntimeResponse
	err = json.NewDecoder(response.Body).Decode(&runtimes)
	if err != nil {
		return RuntimeResponse{}, fmt.Errorf("reading json body: %w", err)
	}

	err = response.Body.Close()
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) && !errors.Is(err, http.ErrBodyReadAfterClose) {
		return RuntimeResponse{}, fmt.Errorf("closing response body: %w", err)
	}

	return runtimes, nil
}
