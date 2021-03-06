package common

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// AwsClientInfo represents aws ec2 mock api interface with some metadata.
type AwsClientInfo struct {
	SVC    ec2iface.EC2API
	Region *string
}

// IsError checks if the given err is an aws error or not or any valid error.
func IsError(err error) (bool, error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "RequestExpired":
				return true, fmt.Errorf("%v: Provided credential has expired", aerr.Code())
			default:
				return true, fmt.Errorf("%s", aerr.Error())
			}
		}

		return true, err
	}

	return false, nil
}

// IsAmiExist returns true if ami exists with the given ami filters.
// nolint:interfacer
func IsAmiExist(svc ec2iface.EC2API, filters []*ec2.Filter) (bool, error) {
	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := svc.DescribeImages(input)
	if ok, err := IsError(err); ok {
		return false, err
	}

	images := result.Images

	if len(images) == 0 {
		return false, nil
	}

	return true, nil
}

// GetAmiInfo returns the a list of aws images that matched the given ami filters.
func GetAmiInfo(ci AwsClientInfo, filters []*ec2.Filter) ([]*ec2.Image, error) {
	if ok, err := IsAmiExist(ci.SVC, filters); !ok {
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("UnableToGetAmiInfo: AMI doesnot exist in Region %s with Filters %v ", *ci.Region, filters)
	}

	input := &ec2.DescribeImagesInput{
		Filters: filters,
	}

	result, err := ci.SVC.DescribeImages(input)
	if ok, err := IsError(err); ok {
		return nil, err
	}

	images := result.Images

	return images, nil
}

// CreateEc2Tags creates tags for given ec2 resources.
//nolint:interfacer
func CreateEc2Tags(svc ec2iface.EC2API, resources []*string, tags []*ec2.Tag) error {
	input := &ec2.CreateTagsInput{
		Resources: resources,
		Tags:      tags,
	}

	_, err := svc.CreateTags(input)
	if err != nil {
		return err
	}

	return nil
}

// FormEc2Tags returns the list of ec2 tags formed from given map of strings.
func FormEc2Tags(tags map[*string]*string) []*ec2.Tag {
	ec2Tags := make([]*ec2.Tag, 0)

	for key, value := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{Key: key, Value: value})
	}

	return ec2Tags
}
