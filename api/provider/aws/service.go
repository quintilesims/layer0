package aws

import (
<<<<<<< 0f6b259f78b9a344646eca80c79aa47cb44bc13e
	"github.com/quintilesims/layer0/api/tag"
	awsc "github.com/quintilesims/layer0/common/aws"
=======
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	awsc "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
>>>>>>> initial change
	"github.com/quintilesims/layer0/common/models"
)

type ServiceProvider struct {
	AWS      *awsc.Client
<<<<<<< 0f6b259f78b9a344646eca80c79aa47cb44bc13e
	TagStore tag.Store
}

func NewServiceProvider(a *awsc.Client, t tag.Store) *ServiceProvider {
=======
	TagStore tag_store.TagStore
	Config   config.APIConfig
}

func NewServiceProvider(a *awsc.Client, t tag_store.TagStore, c config.APIConfig) *ServiceProvider {
>>>>>>> initial change
	return &ServiceProvider{
		AWS:      a,
		TagStore: t,
		Config:   c,
	}
}

func (s *ServiceProvider) Create(req models.CreateServiceRequest) (*models.Service, error) {

	return nil, nil
}

func (s *ServiceProvider) Read(serviceID string) (*models.Service, error) {
	model := &models.Service{}

	clusterName := addLayer0Prefix(s.Config.Instance(), model.EnvironmentID)
	service, err := s.readService(clusterName, serviceID)
	if err != nil {
		return nil, err
	}

	model.DesiredCount = aws.Int64Value(service.DesiredCount)
	model.RunningCount = aws.Int64Value(service.RunningCount)
	model.PendingCount = aws.Int64Value(service.PendingCount)

	for _, deploy := range service.Deployments {
		deployID := aws.StringValue(deploy.TaskDefinition)
		//todo: convert deployid to layer0 deploy id

		deploy := models.Deployment{
			DeploymentID: aws.StringValue(deploy.Id),
			Created:      aws.TimeValue(deploy.CreatedAt),
			Updated:      aws.TimeValue(deploy.UpdatedAt),
			Status:       aws.StringValue(deploy.Status),
			PendingCount: aws.Int64Value(deploy.PendingCount),
			RunningCount: aws.Int64Value(deploy.RunningCount),
			DesiredCount: aws.Int64Value(deploy.DesiredCount),
			DeployID:     deployID,
		}

		model.Deployments = append(model.Deployments, deploy)
	}

	if err := s.updateWithTagInfo(model, serviceID); err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ServiceProvider) readService(clusterName, serviceID string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{}
	input.SetCluster(clusterName)
	input.SetServices([]*string{
		aws.String(serviceID),
	})

	output, err := s.AWS.ECS.DescribeServices(input)
	if err != nil {
		return nil, err
	}

	for _, service := range output.Services {
		return service, nil
	}

	return nil, fmt.Errorf("ecs service '%s' in cluster '%s' does not exist", serviceID, clusterName)
}

func (s *ServiceProvider) List() ([]models.ServiceSummary, error) {
	clusterArns, err := s.listClusters()
	if err != nil {
		return nil, err
	}

	clusterServices := map[*string][]*string{}

	for _, clusterArn := range clusterArns {
		clusterName := strings.Split(aws.StringValue(clusterArn), ":")[5]
		clusterName = strings.Replace(clusterName, "cluster/", "", -1)

		if !strings.HasPrefix(clusterName, s.Config.Instance()) {
			continue
		}

		serviceArns, err := s.listServices(clusterArn)
		if err != nil {
			return nil, err
		}

		clusterServices[clusterArn] = serviceArns
	}

	serviceSummaries := []models.ServiceSummary{}
	for _, serviceArns := range clusterServices {
		for _, serviceArn := range serviceArns {
			ecsServiceID := s.serviceARNToECSServiceID(aws.StringValue(serviceArn))
			service := &models.Service{
				ServiceID: ecsServiceID,
			}
			s.updateWithTagInfo(service)

			summary := models.ServiceSummary{
				ServiceID:       service.ServiceID,
				ServiceName:     service.ServiceName,
				EnvironmentID:   service.EnvironmentID,
				EnvironmentName: service.EnvironmentName,
			}

			serviceSummaries = append(serviceSummaries, summary)
		}
	}

	return serviceSummaries, nil
}

func (s *ServiceProvider) Delete(ServiceID string) error {
	return nil
}

func (s *ServiceProvider) listServices(clusterArn *string) ([]*string, error) {
	serviceArns := []*string{}
	input := &ecs.ListServicesInput{
		Cluster: clusterArn,
	}

	for {
		output, err := s.AWS.ECS.ListServices(input)
		if err != nil {
			return nil, err
		}

		serviceArns = append(serviceArns, output.ServiceArns...)
		input.NextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}

	return serviceArns, nil
}

func (s *ServiceProvider) listClusters() ([]*string, error) {
	input := &ecs.ListClustersInput{}
	clusterArns := []*string{}

	for {
		output, err := s.AWS.ECS.ListClusters(input)
		if err != nil {
			return nil, err
		}

		clusterArns = append(clusterArns, output.ClusterArns...)
		input.NextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}

	return clusterArns, nil
}
func (s *ServiceProvider) updateWithTagInfo(service *models.Service, serviceID string) error {
	tags, err := s.TagStore.SelectByTypeAndID("service", serviceID)
	if err != nil {
		return err
	}
	service.ServiceID = serviceID

	if tag, ok := tags.WithKey("environment_id").First(); ok {
		service.EnvironmentID = tag.Value
	}

	if tag, ok := tags.WithKey("load_balancer_id").First(); ok {
		service.LoadBalancerID = tag.Value
	}

	if tag, ok := tags.WithKey("name").First(); ok {
		service.ServiceName = tag.Value
	}

	if service.EnvironmentID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("environment", service.EnvironmentID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			service.EnvironmentName = tag.Value
		}
	}

	if service.LoadBalancerID != "" {
		tags, err := s.TagStore.SelectByTypeAndID("load_balancer", service.LoadBalancerID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			service.LoadBalancerName = tag.Value
		}
	}

	deployments := []models.Deployment{}
	for _, deploy := range service.Deployments {
		tags, err := s.TagStore.SelectByTypeAndID("deploy", deploy.DeployID)
		if err != nil {
			return err
		}

		if tag, ok := tags.WithKey("name").First(); ok {
			deploy.DeployName = tag.Value
		}

		if tag, ok := tags.WithKey("version").First(); ok {
			deploy.DeployVersion = tag.Value
		}

		deployments = append(deployments, deploy)
	}

	service.Deployments = deployments

	return nil
}

func (s *ServiceProvider) serviceARNToECSServiceID(arn string) string {
	split := strings.SplitN(arn, "/", -1)
	serviceName := split[len(split)-1]
	return serviceName
}
