package integration_testing

import "os"

type MongoDBCredentials struct {
	ProjectID  string
	PublicKey  string
	PrivateKey string
}

type AWSCredentials struct {
	AccessKey         string
	SecretKey         string
	CustomerMasterKey string
	AwsRegion         string
}

func GetCredentialsFromEnv() MongoDBCredentials {
	return MongoDBCredentials{
		ProjectID:  os.Getenv("MONGODB_ATLAS_PROJECT_ID"),
		PublicKey:  os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey: os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
	}
}

func GetAWSCredentialsFromEnv() AWSCredentials {
	return AWSCredentials{
		AccessKey:         os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:         os.Getenv("AWS_SECRET_ACCESS_KEY"),
		CustomerMasterKey: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
		AwsRegion:         os.Getenv("AWS_REGION"),
	}
}
