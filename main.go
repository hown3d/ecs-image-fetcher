package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {

	// If the AWS_SDK_LOAD_CONFIG environment variable is set,
	// or SharedConfigEnable option is used to create the Session
	// the full shared config values will be loaded. This includes credentials,
	// region, and support for assume role. In addition the Session will load
	// its configuration from both the shared config file (~/.aws/config) and
	// shared credentials file (~/.aws/credentials). Both files have the same format.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:                        aws.String("eu-central-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	))
	svc := ecs.New(sess)
	clusters, err := svc.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("Can't get clusters: %v", err)
	}

	for _, cluster := range clusters.ClusterArns {
		services, err := svc.ListServices(&ecs.ListServicesInput{Cluster: cluster})
		if err != nil {
			log.Fatalf("Can't list services: %v", err)
		}
		for _, service := range services.ServiceArns {
			out, err := svc.DescribeServices(&ecs.DescribeServicesInput{
				Cluster:  cluster,
				Services: []*string{service},
			})
			if err != nil {
				log.Fatalf("Can't describe service %v: %v", service, err)
			}
			for _, output := range out.Services {
				out, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: output.TaskDefinition})
				if err != nil {
					log.Fatalf("Can't describe task %v: %v", output.TaskDefinition, err)
				}
				for _, definition := range out.TaskDefinition.ContainerDefinitions {
					fmt.Println(*definition.Image)
				}
			}
		}
	}
}
