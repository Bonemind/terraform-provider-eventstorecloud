package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type WaitForPeeringStateRequest struct {
	OrganizationID string
	ProjectID      string
	PeeringID      string
	State          string
}

func (c *Client) PeeringWaitForState(ctx context.Context, req *WaitForPeeringStateRequest) (*Peering, error) {
	start := time.Now()
	getRequest := &GetPeeringRequest{
		OrganizationID: req.OrganizationID,
		ProjectID:      req.ProjectID,
		PeeringID:      req.PeeringID,
	}

	for {
		resp, err := c.PeeringGet(ctx, getRequest)
		if err != nil {
			return nil, err
		}

		if req.State == "deleted" {
			return &resp.Peering, nil
		}

		if resp.Peering.Status == "defunct" {
			// Resources in a `defunct` state may not update their status right
			//away when being destroyed, so wait a bit before failing the operation.
			elapsed := time.Since(start)
			if elapsed.Seconds() > 30.0 {
				return nil, errors.Errorf("Peering entered a defunct state!")
			}
		}

		if resp.Peering.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		switch resp.Peering.Provider {
		case "aws":
			if _, has := resp.Peering.ProviderPeeringMetadata["peeringLinkId"]; !has {
				time.Sleep(5 * time.Second)
				continue
			}
		case "gcp":
			_, hasProject := resp.Peering.ProviderPeeringMetadata["projectId"]
			_, hasNetwork := resp.Peering.ProviderPeeringMetadata["networkId"]
			if !hasProject || !hasNetwork {
				time.Sleep(5 * time.Second)
				continue
			}
		}

		return &resp.Peering, nil
	}
}
