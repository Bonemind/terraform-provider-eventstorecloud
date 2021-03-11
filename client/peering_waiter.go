package client

import (
	"context"
	"time"
)

type WaitForPeeringStateRequest struct {
	OrganizationID string
	ProjectID      string
	PeeringID      string
	State          string
}

func (c *Client) PeeringWaitForState(ctx context.Context, req *WaitForPeeringStateRequest) (*Peering, error) {
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

		if req.State == "deleted" || req.State == "defunct" {
			return &resp.Peering, nil
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
