package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	endPointSTSHostnameDefault    = "sts.amazonaws.com"
	DefaultRegionSTS              = "us-east-1"
	minSegmentsForSTSRegionalHost = 4
)

func getAWSCredentials(c *config.AWSVars) (*config.Credentials, error) {
	defaultResolver := endpoints.DefaultResolver()
	stsCustResolverFn := func(service, _ string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == sts.EndpointsID {
			resolved, err := ResolveSTSEndpoint(c.Endpoint, c.Region)
			if err != nil {
				return endpoints.ResolvedEndpoint{}, err
			}
			return resolved, nil
		}
		return defaultResolver.EndpointFor(service, c.Region, optFns...)
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(c.Region),
		Credentials:      credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, c.SessionToken),
		EndpointResolver: endpoints.ResolverFunc(stsCustResolverFn),
	}))
	creds := stscreds.NewCredentials(sess, c.AssumeRoleARN)
	if _, err := sess.Config.Credentials.Get(); err != nil {
		return nil, err
	}
	if _, err := creds.Get(); err != nil {
		return nil, err
	}
	secretString, err := secretsManagerGetSecretValue(sess, &aws.Config{Credentials: creds, Region: aws.String(c.Region)}, c.SecretName)
	if err != nil {
		return nil, err
	}
	// TODO could credentials be reused removing Method?
	var secret struct {
		AccessToken  string `json:"access_token"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		PublicKey    string `json:"public_key"`
		PrivateKey   string `json:"private_key"`
	}
	err = json.Unmarshal([]byte(secretString), &secret)
	if err != nil {
		return nil, err
	}
	// TODO: how to read URLs in AWS Secrets Manager?
	return &config.Credentials{
		Method:       "AWS Secrets Manager",
		AccessToken:  secret.AccessToken,
		ClientID:     secret.ClientID,
		ClientSecret: secret.ClientSecret,
		PublicKey:    secret.PublicKey,
		PrivateKey:   secret.PrivateKey,
		BaseURL:      "", // TODO: how to read
		RealmBaseURL: "", // TODO: how to read
	}, nil
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

func ResolveSTSEndpoint(stsEndpoint, secretsRegion string) (endpoints.ResolvedEndpoint, error) {
	ep := stsEndpoint
	if ep == "" {
		r := secretsRegion
		if r == "" {
			r = DefaultRegionSTS
		}
		ep = fmt.Sprintf("https://sts.%s.amazonaws.com/", r)
	}

	signingRegion := DeriveSTSRegionFromEndpoint(ep)

	return endpoints.ResolvedEndpoint{
		URL:           ep,
		SigningRegion: signingRegion,
	}, nil
}

func secretsManagerGetSecretValue(sess *session.Session, creds *aws.Config, secret string) (string, error) {
	svc := secretsmanager.New(sess, creds)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return "", err
	}

	return *result.SecretString, err
}
