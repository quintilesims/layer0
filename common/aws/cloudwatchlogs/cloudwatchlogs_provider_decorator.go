// Generated by go-decorator, DO NOT EDIT
package cloudwatchlogs

import ()

type ProviderDecorator struct {
	Inner     Provider
	Decorator func(name string, call func() error) error
}

func (this *ProviderDecorator) CreateLogGroup(p0 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.CreateLogGroup(p0)
		return err
	}
	err = this.Decorator("CreateLogGroup", call)
	return err
}
func (this *ProviderDecorator) DeleteLogGroup(p0 string) (err error) {
	call := func() error {
		var err error
		err = this.Inner.DeleteLogGroup(p0)
		return err
	}
	err = this.Decorator("DeleteLogGroup", call)
	return err
}
func (this *ProviderDecorator) DescribeLogGroups(p0 string, p1 *string) (v0 []*LogGroup, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeLogGroups(p0, p1)
		return err
	}
	err = this.Decorator("DescribeLogGroups", call)
	return v0, err
}
func (this *ProviderDecorator) DescribeLogStreams(p0 string, p1 string) (v0 []*LogStream, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.DescribeLogStreams(p0, p1)
		return err
	}
	err = this.Decorator("DescribeLogStreams", call)
	return v0, err
}
func (this *ProviderDecorator) GetLogEvents(p0 string, p1 string, p2 string, p3 string, p4 int64) (v0 []*OutputLogEvent, err error) {
	call := func() error {
		var err error
		v0, err = this.Inner.GetLogEvents(p0, p1, p2, p3, p4)
		return err
	}
	err = this.Decorator("GetLogEvents", call)
	return v0, err
}
func (this *ProviderDecorator) FilterLogEvents(p0 *string, p1 *string, p2 *string, p3 []*string, p4 *int64, p5 *int64, p6 *bool) (v0 []*FilteredLogEvent, v1 []*SearchedLogStream, err error) {
	call := func() error {
		var err error
		v0, v1, err = this.Inner.FilterLogEvents(p0, p1, p2, p3, p4, p5, p6)
		return err
	}
	err = this.Decorator("FilterLogEvents", call)
	return v0, v1, err
}

