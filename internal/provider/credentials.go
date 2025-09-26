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
	endPointSTSHostnameDefault = "sts.amazonaws.com"
	defaultRegionSTS   = "us-east-1"
)

func configureCredentialsSTS(cfg *config.Config, secret, region, awsAccessKeyID, awsSecretAccessKey, awsSessionToken, endpoint string) (config.Config, error) {
	defaultResolver := endpoints.DefaultResolver()
	stsCustResolverFn := func(service, _ string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == sts.EndpointsID {
			resolved, err := resolveSTSEndpoint(endpoint, region)
			if err != nil {
				return endpoints.ResolvedEndpoint{}, err
			}
			return resolved, nil
		}
		return defaultResolver.EndpointFor(service, region, optFns...)
	}
	

	sess := session.Must(session.NewSession(&aws.Config{
		Region:              aws.String(region),
		Credentials:         credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsSessionToken),
		EndpointResolver:    endpoints.ResolverFunc(stsCustResolverFn),
	}))

	creds := stscreds.NewCredentials(sess, cfg.AssumeRole.RoleARN)

	_, err := sess.Config.Credentials.Get()
	if err != nil {
		log.Printf("Session get credentials error: %s", err)
		return *cfg, err
	}
	_, err = creds.Get()
	if err != nil {
		log.Printf("STS get credentials error: %s", err)
		return *cfg, err
	}
	secretString, err := secretsManagerGetSecretValue(sess, &aws.Config{Credentials: creds, Region: aws.String(region)}, secret)
	if err != nil {
		log.Printf("Get Secrets error: %s", err)
		return *cfg, err
	}

	var secretData SecretData
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		return *cfg, err
	}
	if secretData.PrivateKey == "" {
		return *cfg, fmt.Errorf("secret missing value for credential PrivateKey")
	}

	if secretData.PublicKey == "" {
		return *cfg, fmt.Errorf("secret missing value for credential PublicKey")
	}

	cfg.PublicKey = secretData.PublicKey
	cfg.PrivateKey = secretData.PrivateKey
	return *cfg, nil
}


func deriveSTSRegionFromEndpoint(ep string) string {
    if ep == "" {
        return ""
    }
    u, err := url.Parse(ep)
    if err != nil {
        return defaultRegionSTS
    }
    host := u.Hostname() // valid values: sts.us-west-2.amazonaws.com or sts.amazonaws.com

    if host == endPointSTSHostnameDefault {
        return defaultRegionSTS
    }

    parts := strings.Split(host, ".")
    if len(parts) >= 4 && parts[0] == "sts" {
        return parts[1]
    }
    return defaultRegionSTS
}

func resolveSTSEndpoint(stsEndpoint, secretsRegion string) (endpoints.ResolvedEndpoint, error) {
    ep := stsEndpoint
    if ep == "" {
        r := secretsRegion
        if r == "" {
            r = defaultRegionSTS
        }
        ep = fmt.Sprintf("https://sts.%s.amazonaws.com/", r)
    }

    signingRegion := deriveSTSRegionFromEndpoint(ep)

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
