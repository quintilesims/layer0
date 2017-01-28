package job

import (
	"github.com/quintilesims/layer0/api/logic"
	"sync"
)

type JobContext interface {
	Copy(string) JobContext
	Request() string
	Logic() *logic.Logic
	LoadBalancerLogic() logic.LoadBalancerLogic
	ServiceLogic() logic.ServiceLogic
	EnvironmentLogic() logic.EnvironmentLogic
}

type L0JobContext struct {
	jobID             string
	logic             *logic.Logic
	loadBalancerLogic logic.LoadBalancerLogic
	serviceLogic      logic.ServiceLogic
	environmentLogic  logic.EnvironmentLogic
	request           string
	mutex             *sync.Mutex
}

func NewL0JobContext(jobID string, lgc *logic.Logic, request string) *L0JobContext {
	return &L0JobContext{
		jobID:             jobID,
		logic:             lgc,
		loadBalancerLogic: logic.NewL0LoadBalancerLogic(*lgc),
		serviceLogic:      logic.NewL0ServiceLogic(*lgc),
		environmentLogic:  logic.NewL0EnvironmentLogic(*lgc),
		request:           request,
		mutex:             &sync.Mutex{},
	}
}

// returns a copy of the job context with a different request object
// this allows us to send different request params to many steps/actions
// while keeping the same underlying mutex and logic references
func (this *L0JobContext) Copy(request string) JobContext {
	return &L0JobContext{
		request:           request,
		jobID:             this.jobID,
		logic:             this.logic,
		loadBalancerLogic: this.loadBalancerLogic,
		serviceLogic:      this.serviceLogic,
		environmentLogic:  this.environmentLogic,
		mutex:             this.mutex,
	}
}

func (this *L0JobContext) Request() string {
	return this.request
}

func (this *L0JobContext) Logic() *logic.Logic {
	return this.logic
}

func (this *L0JobContext) LoadBalancerLogic() logic.LoadBalancerLogic {
	return this.loadBalancerLogic
}

func (this *L0JobContext) ServiceLogic() logic.ServiceLogic {
	return this.serviceLogic
}

func (this *L0JobContext) EnvironmentLogic() logic.EnvironmentLogic {
	return this.environmentLogic
}
