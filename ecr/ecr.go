package ecr

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func LatestImage(repo, matcher string) (string, error) {
	e := newEcr()
	output, err := e.describeImages(repo)
	if err != nil {
		return "", err
	}

	if output == nil {
		return "", errors.New(fmt.Sprintf("No results found for %s", repo))
	}

	for _, id := range output.ImageDetails {
		if exists := containsMatcher(id.ImageTags, matcher); exists {
			for _, tag := range id.ImageTags {
				if *tag != matcher {
					return *tag, nil
				}
			}
		}
	}

	return "", errors.New(fmt.Sprintf("No %s tag found for %s", matcher, repo))
}

type Ecr struct {
	describer func() EcrDescriber
	cache     map[string]interface{}
}

type EcrDescriber interface {
	DescribeImages(input *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error)
}

var describerClient EcrDescriber

func newEcr() *Ecr {
	return &Ecr{
		describer: func() EcrDescriber {
			if describerClient == nil {
				describerClient = ecrClient()
			}
			return describerClient
		},
		cache: make(map[string]interface{}),
	}
}

func ecrClient() (client EcrDescriber) {
	config := aws.NewConfig()
	timeout := 500 * time.Millisecond
	config = config.WithHTTPClient(&http.Client{Timeout: timeout})
	return ecr.New(session.New(config))
}

func (e *Ecr) describeImages(repo string) (output *ecr.DescribeImagesOutput, err error) {
	e.describer()
	if cached, ok := e.cache[repo]; ok {
		output = cached.(*ecr.DescribeImagesOutput)
	} else {
		input := &ecr.DescribeImagesInput{
			RepositoryName: aws.String(repo),
		}

		output, err = e.describer().DescribeImages(input)
		if err != nil {
			return
		}

		e.cache[repo] = output
	}

	return
}

func containsMatcher(tags []*string, matcher string) bool {
	for _, tag := range tags {
		if *tag == matcher {
			return true
		}
	}
	return false
}
