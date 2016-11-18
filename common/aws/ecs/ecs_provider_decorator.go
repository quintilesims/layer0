// Generated by go-decorator, DO NOT EDIT
package ecs

import ()

type ProviderDecorator struct {
	Inner     Provider
	Decorator func(name string, call func() error) error
}

func (this *ProviderDecorator) CreateCluster(p0 string) (v0 *Cluster, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.CreateCluster(p0)
		return err
	}
	err = this.Decorator("CreateCluster", call)
	return v0, err
}
func (this *ProviderDecorator) CreateService(p0 string, p1 string, p2 string, p3 int64, p4 []*LoadBalancer, p5 *string) (v0 *Service, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.CreateService(p0, p1, p2, p3, p4, p5)
		return err
	}
	err = this.Decorator("CreateService", call)
	return v0, err
}
func (this *ProviderDecorator) DeleteCluster(p0 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteCluster(p0)
		return err
	}
	err = this.Decorator("DeleteCluster", call)
	return err
}
func (this *ProviderDecorator) DeleteService(p0 string, p1 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteService(p0, p1)
		return err
	}
	err = this.Decorator("DeleteService", call)
	return err
}
func (this *ProviderDecorator) DeleteTaskDefinition(p0 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteTaskDefinition(p0)
		return err
	}
	err = this.Decorator("DeleteTaskDefinition", call)
	return err
}
func (this *ProviderDecorator) DescribeContainerInstances(p0 string, p1 []*string) (v0 []*ContainerInstance, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeContainerInstances(p0, p1)
		return err
	}
	err = this.Decorator("DescribeContainerInstances", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeCluster(p0 string) (v0 *Cluster, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeCluster(p0)
		return err
	}
	err = this.Decorator("DescribeCluster", call)
	return v0, err
}
func (this *ProviderDecorator) Helper_DescribeClusters() (v0 []*Cluster, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.Helper_DescribeClusters()
		return err
	}
	err = this.Decorator("Helper_DescribeClusters", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeService(p0 string, p1 string) (v0 *Service, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeService(p0, p1)
		return err
	}
	err = this.Decorator("DescribeService", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeServices(p0 string, p1 []*string) (v0 []*Service, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeServices(p0, p1)
		return err
	}
	err = this.Decorator("DescribeServices", call)
	return v0, err
}
func (this *ProviderDecorator) Helper_DescribeServices(p0 string) (v0 []*Service, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.Helper_DescribeServices(p0)
		return err
	}
	err = this.Decorator("Helper_DescribeServices", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeTaskDefinition(p0 string) (v0 *TaskDefinition, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeTaskDefinition(p0)
		return err
	}
	err = this.Decorator("DescribeTaskDefinition", call)
	return v0, err
}
func (this *ProviderDecorator) Helper_DescribeTaskDefinitions(p0 string) (v0 []*TaskDefinition, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.Helper_DescribeTaskDefinitions(p0)
		return err
	}
	err = this.Decorator("Helper_DescribeTaskDefinitions", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeTasks(p0 string, p1 []*string) (v0 []*Task, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeTasks(p0, p1)
		return err
	}
	err = this.Decorator("DescribeTasks", call)
	return v0, err
}
func (this *ProviderDecorator) ListClusters() (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListClusters()
		return err
	}
	err = this.Decorator("ListClusters", call)
	return v0, err
}
func (this *ProviderDecorator) ListContainerInstances(p0 string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListContainerInstances(p0)
		return err
	}
	err = this.Decorator("ListContainerInstances", call)
	return v0, err
}
func (this *ProviderDecorator) ListServices(p0 string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListServices(p0)
		return err
	}
	err = this.Decorator("ListServices", call)
	return v0, err
}
func (this *ProviderDecorator) ListTasks(p0 string, p1 *string, p2 *string, p3 *string, p4 *string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListTasks(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("ListTasks", call)
	return v0, err
}
func (this *ProviderDecorator) ListTaskDefinitions(p0 string, p1 *string) (v0 []*string, v1 *string, err error) {
	call := func() error {
		var err error
		v0, v1, err = this.Inner.ListTaskDefinitions(p0, p1)
		return err
	}
	err = this.Decorator("ListTaskDefinitions", call)
	return v0, v1, err
}
func (this *ProviderDecorator) Helper_ListTaskDefinitions(p0 string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.Helper_ListTaskDefinitions(p0)
		return err
	}
	err = this.Decorator("Helper_ListTaskDefinitions", call)
	return v0, err
}
func (this *ProviderDecorator) ListTaskDefinitionsPages(p0 string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListTaskDefinitionsPages(p0)
		return err
	}
	err = this.Decorator("ListTaskDefinitionsPages", call)
	return v0, err
}
func (this *ProviderDecorator) ListTaskDefinitionFamilies(p0 string, p1 *string) (v0 []*string, v1 *string, err error) {
	call := func() error {
		var err error
		v0, v1, err = this.Inner.ListTaskDefinitionFamilies(p0, p1)
		return err
	}
	err = this.Decorator("ListTaskDefinitionFamilies", call)
	return v0, v1, err
}
func (this *ProviderDecorator) ListTaskDefinitionFamiliesPages(p0 string) (v0 []*string, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.ListTaskDefinitionFamiliesPages(p0)
		return err
	}
	err = this.Decorator("ListTaskDefinitionFamiliesPages", call)
	return v0, err
}
func (this *ProviderDecorator) RegisterTaskDefinition(p0 string, p1 string, p2 string, p3 []*ContainerDefinition, p4 []*Volume) (v0 *TaskDefinition, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.RegisterTaskDefinition(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("RegisterTaskDefinition", call)
	return v0, err
}
func (this *ProviderDecorator) RunTask(p0 string, p1 string, p2 int64, p3 *string, p4 []*ContainerOverride) (v0 []*Task, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.RunTask(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("RunTask", call)
	return v0, err
}
func (this *ProviderDecorator) StartTask(p0 string, p1 string, p2 *TaskOverride, p3 []*string, p4 *string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.StartTask(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("StartTask", call)
	return err
}
func (this *ProviderDecorator) StopTask(p0 string, p1 string, p2 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.StopTask(p0, p1, p2)
		return err
	}
	err = this.Decorator("StopTask", call)
	return err
}
func (this *ProviderDecorator) UpdateService(p0 string, p1 string, p2 *string, p3 *int64) (err error) {
	call := func() error {
		var err error
		err = this.Inner.UpdateService(p0, p1, p2, p3)
		return err
	}
	err = this.Decorator("UpdateService", call)
	return err
}
