package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	endPointSTSHostnameDefault    = "sts.amazonaws.com"
	DefaultRegionSTS              = "us-east-1"
	minSegmentsForSTSRegionalHost = 4
)

func getAWSCredentials(ctx context.Context, c *config.AWSVars) (*config.Credentials, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(c.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, c.SessionToken)),
	)
	if err != nil {
		return nil, err
	}
	ep, signingRegion := ResolveSTSEndpoint(c.Endpoint, c.Region)

	stsClient := sts.NewFromConfig(cfg, func(o *sts.Options) {
		o.Region = signingRegion
		if ep != "" {
			o.EndpointResolver = sts.EndpointResolverFromURL(ep)
		}
	})

	assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, c.AssumeRoleARN)

	smCfg := cfg
	smCfg.Credentials = aws.NewCredentialsCache(assumeRoleProvider)
	smClient := secretsmanager.NewFromConfig(smCfg)

	secretString, err := secretsManagerGetSecretValue(ctx, smClient, c.SecretName)
	if err != nil {
		return nil, err
	}
	var secret config.Credentials
	err = json.Unmarshal([]byte(secretString), &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func DeriveSTSRegionFromEndpoint(ep string) string {
	if ep == "" {
		return ""
	}
	u, err := url.Parse(ep)
	if err != nil {
		return DefaultRegionSTS
	}
	host := u.Hostname() // valid values: sts.us-west-2.amazonaws.com or sts.amazonaws.com

	if host == endPointSTSHostnameDefault {
		return DefaultRegionSTS
	}

	parts := strings.Split(host, ".")
	if len(parts) >= minSegmentsForSTSRegionalHost && parts[0] == "sts" {
		return parts[1]
	}
	return DefaultRegionSTS
}

func ResolveSTSEndpoint(stsEndpoint, secretsRegion string) (ep, signingRegion string) {
	ep = stsEndpoint
	if ep == "" {
		r := secretsRegion
		if r == "" {
			r = DefaultRegionSTS
		}
		ep = fmt.Sprintf("https://sts.%s.amazonaws.com/", r)
	}
	signingRegion = DeriveSTSRegionFromEndpoint(ep)
	return ep, signingRegion
}

func secretsManagerGetSecretValue(ctx context.Context, client *secretsmanager.Client, secret string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		switch e := err.(type) {
		case *types.ResourceNotFoundException:
			log.Println("ResourceNotFoundException", e.Error())
		case *types.InvalidParameterException:
			log.Println("InvalidParameterException", e.Error())
		case *types.InvalidRequestException:
			log.Println("InvalidRequestException", e.Error())
		case *types.DecryptionFailure:
			log.Println("DecryptionFailure", e.Error())
		case *types.InternalServiceError:
			log.Println("InternalServiceError", e.Error())
		default:
			log.Println(err.Error())
		}
		return "", err
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret string is nil for secret %s", secret)
	}
	return *result.SecretString, nil
}
