package main

//DescribeJobResult -
type DescribeJobResult struct {
	Jobs []struct {
		JobName  string `json:"jobName"`
		JobID    string `json:"jobId"`
		JobQueue string `json:"jobQueue"`
		Status   string `json:"status"`
		Attempts []struct {
			Container struct {
				ContainerInstanceArn string        `json:"containerInstanceArn"`
				TaskArn              string        `json:"taskArn"`
				ExitCode             int           `json:"exitCode"`
				LogStreamName        string        `json:"logStreamName"`
				NetworkInterfaces    []interface{} `json:"networkInterfaces"`
			} `json:"container"`
			StartedAt    int64  `json:"startedAt"`
			StoppedAt    int64  `json:"stoppedAt"`
			StatusReason string `json:"statusReason"`
		} `json:"attempts"`
		StatusReason  string        `json:"statusReason"`
		CreatedAt     int64         `json:"createdAt"`
		StartedAt     int64         `json:"startedAt"`
		StoppedAt     int64         `json:"stoppedAt"`
		DependsOn     []interface{} `json:"dependsOn"`
		JobDefinition string        `json:"jobDefinition"`
		Parameters    struct{}      `json:"parameters"`
		Container     struct {
			Image       string        `json:"image"`
			Vcpus       int           `json:"vcpus"`
			Memory      int           `json:"memory"`
			Command     []string      `json:"command"`
			JobRoleArn  string        `json:"jobRoleArn"`
			Volumes     []interface{} `json:"volumes"`
			Environment []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"environment"`
			MountPoints          []interface{} `json:"mountPoints"`
			Ulimits              []interface{} `json:"ulimits"`
			ExitCode             int           `json:"exitCode"`
			ContainerInstanceArn string        `json:"containerInstanceArn"`
			TaskArn              string        `json:"taskArn"`
			LogStreamName        string        `json:"logStreamName"`
			NetworkInterfaces    []interface{} `json:"networkInterfaces"`
			ResourceRequirements []interface{} `json:"resourceRequirements"`
		} `json:"container"`
	} `json:"jobs"`
}
