package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var describerClient EcrDescriber

type EcrDescriber interface {
	DescribeImages(input *ecr.DescribeImagesInput) (*ecr.DescribeImagesOutput, error)
}

type Ecr struct {
	describer func() EcrDescriber
	cache     map[string]interface{}
}

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

func (e *Ecr) describeImages(repo string) (output *ecr.DescribeImagesOutput) {
	e.describer()
	if cached, ok := e.cache[repo]; ok {
		output = cached.(*ecr.DescribeImagesOutput)
	} else {
		input := &ecr.DescribeImagesInput{
			RepositoryName: aws.String(repo),
		}
		var err error
		output, err = e.describer().DescribeImages(input)
		if err != nil {
			log.Panicln(err.Error())
			return nil
		}
		e.cache[repo] = output
	}
	return
}

func (e *Ecr) LatestImage(repo, matcher string) string {
	output := e.describeImages(repo)
	if output == nil {
		log.Fatalf("No results found for %s", repo)
	}
	for _, id := range output.ImageDetails {
		if exists := containsMatcher(id.ImageTags, matcher); exists {
			for _, tag := range id.ImageTags {
				if *tag != matcher {
					return *tag
				}
			}
		}
	}
	log.Fatalf("No latest tag found for %s", repo)
	return ""
}

func containsMatcher(tags []*string, matcher string) bool {
	for _, tag := range tags {
		if *tag == matcher {
			return true
		}
	}
	return false
}

type EcrInit struct {
	client  *Ecr
	ecrInit sync.Once
}

func (e *EcrInit) initEcr() {
	if e.client == nil {
		e.client = newEcr()
	}
}
