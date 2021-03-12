package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type WaitForManagedClusterStateRequest struct {
	OrganizationID string
	ProjectID      string
	ClusterID      string
	State          string
}

func (c *Client) ManagedClusterWaitForState(ctx context.Context, req *WaitForManagedClusterStateRequest) error {
	start := time.Now()

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
			// Resources in a `defunct` state may not update their status right
			//away when being destroyed, so wait a bit before failing the operation.
			elapsed := time.Since(start)
			if elapsed.Seconds() > 30.0 {
				return errors.Errorf("Cluster entered a defunct state!")
			}
		}

		if resp.ManagedCluster.Status != req.State {
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}
}
