package task

import "github.com/GoSimplicity/AI-CloudOps/internal/model"

type CronTaskPayload struct {
	JobID     int                    `json:"job_id"`
	JobName   string                 `json:"job_name"`
	TaskType  model.CronJobType      `json:"task_type"`
	TriggerBy string                 `json:"trigger_by,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
