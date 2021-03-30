package esc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceScheduledBackup() *schema.Resource {

	return &schema.Resource{
		Create: resourceScheduledBackupCreate,
		Exists: resourceScheduledBackupExists,
		Read:   resourceScheduledBackupRead,
		Delete: resourceScheduledBackupDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "Human readable description of the job",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"project_id": {
				Description: "ID of the project in which the backup exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"schedule": {
				Description: "Schedule for the backup, defined using restricted subset of cron",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"source_cluster_id": {
				Description: "the ID of the cluster to back up",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"backup_description": {
				Description: "backup_description",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"max_backup_count": {
				Description: "The maximum number of backups to keep for this job",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
			},
		},
	}
}

func resourceScheduledBackupCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)

	request := client.CreateJobRequest{
		Data: map[string]interface{}{
			"clusterId":      d.Get("source_cluster_id").(string),
			"description":    d.Get("backup_description").(string),
			"maxBackupCount": d.Get("max_backup_count").(int),
		},
		Description: d.Get("description").(string),
		Schedule:    d.Get("schedule").(string),
		Type:        "ScheduledBackup",
	}

	resp, err := c.client.CreateJob(context.Background(), c.organizationId, projectId, request)

	if err != nil {
		return err
	}

	d.SetId(resp.Id)

	return resourceScheduledBackupRead(d, meta)
}

func resourceScheduledBackupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	jobId := d.Id()

	cluster, err := c.client.GetJob(context.Background(), c.organizationId, projectId, jobId)
	if err != nil {
		return false, nil
	}
	if cluster.Job.Status == client.StateDeleted {
		return false, nil
	}

	return true, nil
}

func resourceScheduledBackupRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)
	projectId := d.Get("project_id").(string)
	jobId := d.Id()

	resp, err := c.client.GetJob(context.Background(), c.organizationId, projectId, jobId)
	if err != nil {
		return err
	}
	if err := d.Set("description", resp.Job.Description); err != nil {
		return err
	}
	if err := d.Set("project_id", resp.Job.ProjectId); err != nil {
		return err
	}
	if err := d.Set("schedule", resp.Job.Schedule); err != nil {
		return err
	}
	if err := d.Set("source_cluster_id", resp.Job.Data["clusterId"]); err != nil {
		return err
	}
	if err := d.Set("backup_description", resp.Job.Data["description"]); err != nil {
		return err
	}
	if err := d.Set("max_backup_count", resp.Job.Data["maxBackupCount"]); err != nil {
		return err
	}

	return nil
}

func resourceScheduledBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)
	projectId := d.Get("project_id").(string)
	jobId := d.Get("project_id").(string)

	resp, err := c.client.GetJob(context.Background(), c.organizationId, projectId, jobId)
	if err != nil {
		return fmt.Errorf("could not retrieve job %q for project %q, so cannot update", jobId, projectId)
	}

	job := resp.Job

	for _, field := range []string{"description", "project_id", "schedule", "source_cluster_id", "backup_description", "max_backup_count"} {
		if d.HasChange(field) {
			return fmt.Errorf("scheduled jobs cannot be updated after being created, so cannot change field %q of Job with ID %q (%s)", field, jobId, job.Description)
		}
	}
	return nil
}

func resourceScheduledBackupDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	jobId := d.Id()

	if err := c.client.DeleteJob(context.Background(), c.organizationId, projectId, jobId); err != nil {
		return err
	}

	start := time.Now()
	for {
		resp, err := c.client.GetJob(context.Background(), c.organizationId, projectId, jobId)
		if err != nil {
			return fmt.Errorf("error polling job %q (%q) to see if it actually got deleted", jobId, d.Get("description"))
		}
		elapsed := time.Since(start)
		if elapsed.Seconds() > 30.0 {
			return errors.Errorf("job %q (%q) does not seem to be deleting", jobId, d.Get("description"))
		}
		if resp.Job.Status == "deleted" {
			return nil
		}
		time.Sleep(1.0)
	}
}
