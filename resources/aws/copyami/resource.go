package copyami

import (
	"log"
	"os"

	clog "example.com/proffer/common/clogger"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/mapstructure"
)

var (
	clogger = clog.New(os.Stdout, "aws-copyami | ", log.Lmsgprefix)
)

type Source struct {
	Environment string              `yaml:"environment"`
	Region      *string             `yaml:"region"`
	AmiFilters  map[*string]*string `yaml:"amiFilters"`
}

type Target struct {
	Regions      []*string         `yaml:"regions"`
	AddExtraTags map[*string]*string `yaml:"addExtraTags"`
}

type Config struct {
	Source Source `yaml:"source"`
	Target Target `yaml:"target"`
}

type Resource struct {
	Config Config `yaml:"config"`
}

func (r *Resource) Prepare(rawConfig map[string]interface{}) error {

	var c Config

	if err := mapstructure.Decode(rawConfig, &c); err != nil {
		return err
	}

	r.Config = c

	return nil
}

func (r *Resource) Run() error {

	source := r.Config.Source
	target := r.Config.Target
	var amiFilters []*ec2.Filter

	for filterName, filterValue := range source.AmiFilters {
		f := &ec2.Filter{
			Name:   filterName,
			Values: []*string{filterValue},
		}
		amiFilters = append(amiFilters, f)
	}

	srcAmiInfo := SrcAmiInfo{
		Region:  source.Region,
		Filters: amiFilters,
	}

	targetInfo := TargetInfo{
		Regions: target.Regions,
		ExtraTags: formEc2Tags(target.AddExtraTags),
	}

	if err := copyAmi(srcAmiInfo, targetInfo); err != nil {
		return err
	}

	return nil
}
