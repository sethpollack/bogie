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

var ecrCache *ecrClient

func LatestImage(skip bool) func(string, string) (string, error) {
	return func(repo, matcher string) (string, error) {
		if skip {
			return matcher, nil
		}

		e := getClient()
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
}

type ecrClient struct {
	describer ecrDescriber
	cache     map[string]interface{}
}

type ecrDescriber interface {
	DescribeImagesPages(input *ecr.DescribeImagesInput, fn func(*ecr.DescribeImagesOutput, bool) bool) error
}

func getClient() *ecrClient {
	if ecrCache == nil {
		ecrCache = &ecrClient{
			describer: newClient(),
			cache:     make(map[string]interface{}),
		}
	}

	return ecrCache
}

func newClient() (client ecrDescriber) {
	config := aws.NewConfig()
	timeout := 500 * time.Millisecond
	config = config.WithHTTPClient(&http.Client{Timeout: timeout})
	return ecr.New(session.New(config))
}

func (e *ecrClient) describeImages(repo string) (*ecr.DescribeImagesOutput, error) {
	if cached, ok := e.cache[repo]; ok {
		return cached.(*ecr.DescribeImagesOutput), nil
	}

	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repo),
		Filter: &ecr.DescribeImagesFilter{
			TagStatus: aws.String("TAGGED"),
		},
	}

	output := &ecr.DescribeImagesOutput{}
	err := e.describer.DescribeImagesPages(input,
		func(page *ecr.DescribeImagesOutput, lastPage bool) bool {
			output.ImageDetails = append(output.ImageDetails, page.ImageDetails...)
			return !lastPage
		},
	)
	if err != nil {
		return nil, err
	}

	e.cache[repo] = output
	return output, nil
}

func containsMatcher(tags []*string, matcher string) bool {
	for _, tag := range tags {
		if *tag == matcher {
			return true
		}
	}
	return false
}
