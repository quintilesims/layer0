package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/quintilesims/layer0/common/aws/provider"
	"time"
)

type Provider interface {
	GetMetricStatistics(namespace, metricName string, period int64, statistics []string, dimensions []*cloudwatch.Dimension, startTime, endTime time.Time) ([]cloudwatch.Datapoint, error)
	ListMetrics(namespace, metricName string, dimensionFilters []*cloudwatch.DimensionFilter) ([]cloudwatch.Metric, error)
}

type CloudWatch struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (CloudWatchInternal, error)
}

type CloudWatchInternal interface {
	GetMetricStatistics(input *cloudwatch.GetMetricStatisticsInput) (output *cloudwatch.GetMetricStatisticsOutput, err error)
	ListMetrics(input *cloudwatch.ListMetricsInput) (output *cloudwatch.ListMetricsOutput, err error)
}

func NewCloudWatch(credProvider provider.CredProvider, region string) (Provider, error) {
	cloudwatch := CloudWatch{
		credProvider,
		region,
		func() (CloudWatchInternal, error) {
			return Connect(credProvider, region)
		},
	}
	_, err := cloudwatch.Connect()
	if err != nil {
		return nil, err
	}
	return &cloudwatch, nil
}

func Connect(credProvider provider.CredProvider, region string) (CloudWatchInternal, error) {
	connection, err := provider.GetCloudWatchConnection(credProvider, region)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (this *CloudWatch) GetMetricStatistics(namespace, metricName string, period int64, statistics []string, dimensions []*cloudwatch.Dimension, startTime, endTime time.Time) ([]cloudwatch.Datapoint, error) {
	//panic("This is broken until provider.ToStrPtrList is reimplemented")
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(period),
		StartTime:  aws.Time(endTime),
		EndTime:    aws.Time(startTime),
		Statistics: []*string{},
		//Statistics: provider.ToStrPtrList(statistics),
		Dimensions: dimensions,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}
	resp, err := connection.GetMetricStatistics(input)
	if err != nil {
		return nil, err
	}

	datapoints := []cloudwatch.Datapoint{}
	for _, dp := range resp.Datapoints {
		datapoints = append(datapoints, *dp)
	}

	return datapoints, nil
}

func (this *CloudWatch) ListMetrics(namespace, metricName string, dimensionFilters []*cloudwatch.DimensionFilter) ([]cloudwatch.Metric, error) {
	input := &cloudwatch.ListMetricsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Dimensions: dimensionFilters,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}
	resp, err := connection.ListMetrics(input)
	if err != nil {
		return nil, err
	}

	metrics := []cloudwatch.Metric{}
	for _, metric := range resp.Metrics {
		metrics = append(metrics, *metric)
	}

	return metrics, nil
}
