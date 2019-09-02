package config

import (
	"fmt"

	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/factory"
	// this blank import is used to register the S3 driver with the storage driver factory
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
)

// S3 is the Config implementation for the S3 client
type S3 struct {
	AccessKeyFile string `envconfig:"ACCESS_KEY_FILE" default:"/var/run/secrets/deis/objectstore/creds/accesskey"`
	SecretKeyFile string `envconfig:"SECRET_KEY_FILE" default:"/var/run/secrets/deis/objectstore/creds/secretkey"`
	RegionFile    string `envconfig:"REGION_FILE" default:"/var/run/secrets/deis/objectstore/creds/region"`
	BucketFile    string `envconfig:"BUCKET_FILE" default:"/var/run/secrets/deis/objectstore/creds/bucket"`
	S3Host        string `envconfig:"S3_HOST"`
}

// CreateDriver is the Config interface implementation
func (s S3) CreateDriver() (driver.StorageDriver, error) {
	files, err := readFiles(true, s.AccessKeyFile, s.SecretKeyFile, s.RegionFile, s.BucketFile)
	if err != nil {
		return nil, err
	}
	key, secret, region, bucket := files[0], files[1], files[2], files[3]
	host := parseEnvVar(s.S3Host)
	params := map[string]interface{}{
		"accesskey": key,
		"secretkey": secret,
		"region":    region,
		"bucket":    bucket,
	}
	if host != "" {
		params["regionendpoint"] = fmt.Sprintf("https://%s", host)
	}
	return factory.Create("s3aws", params)
}

// String is the fmt.Stringer interface implementation
func (s S3) String() string {
	return S3StorageType.String()
}
