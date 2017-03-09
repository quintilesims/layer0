package scheduler

type EnvironmentResourceManager struct {
	RequiredServiceResources func() ([]ServiceResource, error)
	RequiredTaskResources    func() ([]TaskResource, error)
	EnvironmentScalar        EnvironmentScalar
}

type EnvironmentScalar interface {
	ScaleTo(int) error
}

type ContainerResource struct {
	InstancePorts []int
	Memory        bytesize.Bytesize
}

/*
	Problems: The calculations are theoretical - they are based on the assumption that every ec2 instance in the
	cluster is the same size. The math for adding up all the current memory doesn't work well, because if one ec2 instance
	has 1GB and another has 2GB, we won't have room for a task that needs 3GB. We would need to theoretically place
	the task first.
	Another issue with memory is even if we do the 'sum all required' method, this only works assuming ecs will move around
	services that are already running to better optimize the cluster. I don't think it actually does this.
	So what can we do? Well, we can do the theoretical approach:
	currentInstance := getCurrentResources() // this includes Desired ec2 instances that will have full memory, ports
	for _, pendingResource := range pendingResources {
		var instanceFound bool
		for _, instance := range currentInstances {
			if instance.HasAavailableResources(pendingResource){
				// now, next time we calculate it will have port X taken up and memory X taken up
				instance.UseUpResources(pendingResources)
				instanceFound = true
				break
			}
			if !instanceFound{
				requiredInstances++
				currentInstances = append(currentInstances, NewInstance())
			}
		}
	}
*/

func (e *EnvironmentResourceManager) Approach2() error {
	pendingServiceResources, err := r.PendingServiceResources()
	if err != nil {
		return err
	}

	pendingTaskResources, err := r.PendingTaskResources()
	if err != nil {
		return err
	}

	// this is calculated by curent + desired. the desired ones that don't exist yet
	// will have full ports, moeomroy open
	resourceProviders, err := r.ResourceProviders()
	if err != nil {
		return err
	}

	pendingResourceConsumers = append(pendingServiceResources, pendingTaskResources...)

	for _, consumer := range pendingResourceConsumers {
		hasRoom := false

		for _, provider := range resourceProviders {
			if provider.HasResourcesFor(consumer) {
				hasRoom = true
				provider.SubtractResourcesFor(consumer)
				break
			}
		}

		if !hasRoom {
			newProvider := scalar.NewProvider()
			if !newProvider.HasResourcesFor(consumer) {
				log.Error("This resource consumer is too large for our current instance size!")
				continue
			}

			resourcesProviders = append(resourceProviders, newProvider)
			newProvider.SubtractResourcesFor(consumer)
		}
	}

	if newProvidersRequired == 0 {
		return nil
	}
	
	newScale := len(resourceProviders)
	return e.Scalar.ScaleTo(newScale)

	// todo: what about scaling down? different function perhaps
	// maybe the first-approach one can work
	// our current sclae-down function works pretty well - just go through each ec2 instanec and terminate unused ones
}

// this would run probably every 5 minutes or so
func ScaleDown() error {
	asg := getASG()
	
	// we can't scale down anymore or we'd hit the min
	if asg.Desired == asg.Min {
		return nil
	}

	// is this too simple? seems do-able
	// you could have a service/task that will never get its required resources
	// prevent this from ever running
	// but on the other hand, if this runs during a service/task creation (which it will)
	// you will likely have unused ec2 instances for some time while the instance spins up
	// ...this seems the safer option
	pendingResources := getAllPendingResources()
	if len(pendingResources) > 0 {
		return nil
	}

	for _, instance := range asg.Instances {
		// if the instnace isn't being used _and_ we can scale the asg down
		if instance.NotUsed() and asg.Desired-1 >= asg.Min {
			instance.terminate()
			asg.SetDesired(-1)
			asg.SetMax(-1)
		}
	}

	// todo: do we want to honor a max? 

	return nil
}

func (e *EnvironmentResourceManager) Run() error {
	requiredServiceResources, err := r.GetServiceResources()
	if err != nil {
		return err
	}

	requiredTaskResources, err := r.GetTaskResources()
	if err != nil {
		return err
	}

	// calculate the minimum number of ec2 instances needed if all we worried about was the ports:
	// the minimum number would be the largest set of overlapping ports.
	// Example:
	//	- 80, 8000, 5000
	// 	- 80, 3000
	//	- 5000
	//	- 80, 8000, 22
	//	- 80, 8000
	//
	// in this case, the overlapping port sets are: [80, 80, 80, 80] [8000, 8000, 8000] [5000, 5000] [3000] [22]
	// the largest port set is [80, 80, 80, 80]
	// this means we could have a minimum of 4 ec2 instances in our environment

	allRequiredPorts := getAllRequiredPorts(requiredServiceResources, requiredTaskResources)
	minInstancesForPorts := calculateMinInstancesForPorts(allRequiredPorts)

	// now we hit some issues - What if the users have a custom 'PlacementStrategy'?
	// we could not allow that and always assume binPacking strategy

	// assume we always use the binpacking strategy
	// The, calculate the minimum number of ec2 instances needed if all we worried about was memory:
	// Example:
	//      - 10GB: 80, 8000, 5000
	//      - 8GB: 80, 3000
	//      - 5GB: 5000
	//      - 10GB: 80, 8000, 22
	//      - 10GB: 80, 8000
	//
	// just add up the total required memory: 10+8+5+10+10 = 43GB
	// assume each ec2 instance gives 10GB, we get: 43/10 = 4.6. We round up so we get 5 instances minimum
	// we would also need to throw an error if a task requires more memory than provided in a memoryPerInstance

	allRequiredMemory := getAllRequiredMemory(requiredServiceResources, requiredTaskResources)
	minInstancesForMemory := calculateMinInstancesForMemory(allRequiredMemory, memoryPerInstance)

	// now, tell the EnvironmentScalar to scale to the minimum required size
	minInstances := max(minInstancesForPorts, minInstancesForMemory)
	return e.EnvironmentScalar.ScaleTo(minInstances)
}

// todo: make sure when calculating task resoources we don't add stopped tasks
// also, add a resource for each instance of a task/service just for sanity
// for i := 0; i < task.Copies
// for i := 0; i < service.Scale

func (e EnvironmentScalar) ScaleTo(count int) error {
	// todo: don't scale under minimum count
	// todo: run cleanup async

	// todo: optimize scaledown process:
	// if we just lower the asg count, it will terminate a random ec2 instance
	// ideally, we terminate an unused ec2 instance if it exists
	// the process could look something like:

	asg := e.getASG()

	if count < asg.MinCount {
		log.Warn("trying to go lower than min count!")
		return nil
	}

	if asg.MaxCount != count {
		asg.SetMaxCount(count)
	}

	if asg.DesiredCount != count {
		asg.SetDesiredCount(count)
	}

	// if we are scaling down, terminate up to N unused instances
	if numToTerminate := asg.CurrentCount - count; numToTerminate > 0 {
		instances := e.GetInstances()

		numTerminated := 0
		for _, instance := range instances {
			if instance.NotInUse() {
				e.TerminateInstanceInASG(instance)
				numTerminated++
			}

			if numTerminated == numToTeminate {
				break
			}
		}
	}
