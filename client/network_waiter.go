package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type WaitForNetworkStateRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkID      string
	State          string
}

func (c *Client) NetworkWaitForState(ctx context.Context, req *WaitForNetworkStateRequest) error {
	start := time.Now()

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
			// Resources in a `defunct` state may not update their status right
			//away when being destroyed, so wait a bit before failing the operation.
			elapsed := time.Since(start)
			if elapsed.Seconds() > 30.0 {
				return errors.Errorf("Network entered a defunct state!")
			}
		}

		if resp.Network.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}
}
