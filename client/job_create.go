package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateJobRequest struct {
	Data        map[string]interface{} `json:"data"`
	Description string                 `json:"description"`
	Schedule    string                 `json:"schedule"`
	Type        string                 `json:"type"`
}

type CreateJobResponse struct {
	Id string `json:"id"`
}

func (c *Client) CreateJob(ctx context.Context, organizationId string, projectId string, createJobRequest CreateJobRequest) (*CreateJobResponse, error) {
	requestBody, err := json.Marshal(createJobRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	url := *c.apiURL
	url.Path = "/organizations/{organizationId}/projects/{projectId}/jobs"
	url.Path = strings.Replace(url.Path, "{"+"organizationId"+"}", organizationId, -1)
	url.Path = strings.Replace(url.Path, "{"+"projectId"+"}", projectId, -1)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing request for CreateJob: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request for CreateJob: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, translateStatusCode(resp.StatusCode, "CreateJob", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateJobResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
