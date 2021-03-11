package client

import (
	"context"
	"time"
)

type WaitForManagedClusterStateRequest struct {
	OrganizationID string
	ProjectID      string
	ClusterID      string
	State          string
}

func (c *Client) ManagedClusterWaitForState(ctx context.Context, req *WaitForManagedClusterStateRequest) error {
	getRequest := &GetManagedClusterRequest{
		OrganizationID: req.OrganizationID,
		ProjectID:      req.ProjectID,
		ClusterID:      req.ClusterID,
	}

	for {
		resp, err := c.ManagedClusterGet(ctx, getRequest)
		if err != nil {
			return err
		}

		if resp.ManagedCluster.Status == "defunct" {
			return nil
		}

		if resp.ManagedCluster.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}
}
