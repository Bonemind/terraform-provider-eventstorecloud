package client

import (
	"context"
	"time"
)

type WaitForNetworkStateRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkID      string
	State          string
}

func (c *Client) NetworkWaitForState(ctx context.Context, req *WaitForNetworkStateRequest) error {
	getRequest := &GetNetworkRequest{
		OrganizationID: req.OrganizationID,
		ProjectID:      req.ProjectID,
		NetworkID:      req.NetworkID,
	}

	for {
		resp, err := c.NetworkGet(ctx, getRequest)
		if err != nil {
			return err
		}

		if resp.Network.Status == "defunct" {
			return nil
		}

		if resp.Network.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}
}
