package job

type JobType string

func (j JobType) String() string {
	return string(j)
}

const (
	CreateDeployJob       JobType = "CreateDeploy"
	CreateEnvironmentJob  JobType = "CreateEnvironment"
	CreateLoadBalancerJob JobType = "CreateLoadBalancer"
	CreateServiceJob      JobType = "CreateService"
	CreateTaskJob         JobType = "CreateTask"
	DeleteDeployJob       JobType = "DeleteDeploy"
	DeleteEnvironmentJob  JobType = "DeleteEnvironment"
	DeleteLoadBalancerJob JobType = "DeleteLoadBalancer"
	DeleteServiceJob      JobType = "DeleteService"
	DeleteTaskJob         JobType = "DeleteTask"
	UpdateDeployJob       JobType = "UpdateDeploy"
	UpdateTaskJob         JobType = "UpdateTask"
)

type JobStatus string

func (j JobStatus) String() string {
	return string(j)
}

const (
	Pending    JobStatus = "Pending"
	InProgress JobStatus = "InProgress"
	Completed  JobStatus = "Completed"
	Error      JobStatus = "Error"
)
