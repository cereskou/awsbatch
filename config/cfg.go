package config

import (
	"fmt"
	"os"
	"path/filepath"

	"ditto.co.jp/submit/util"
	"github.com/aws/aws-sdk-go/aws/credentials"
	jsoniter "github.com/json-iterator/go"
)

var (
	setting          *Config
	awsDefaultRegion = "ap-northeast-1"
	json             = jsoniter.ConfigCompatibleWithStandardLibrary
)

//AwsConfig - for aws service
type AwsConfig struct {
	AccessKey string `json:"AWS_ACCESS_KEY_ID"`
	SecretKey string `json:"AWS_SECRET_ACCESS_KEY"`
	Region    string `json:"AWS_DEFAULT_REGION"`
	Proxy     string `json:"Proxy"`
	Timeout   int    `json:"TimeOut"`
}

//Config -
type Config struct {
	JobName       string    `json:"JobName"`
	JobQueue      string    `json:"JobQueue"`
	JobDefinition string    `json:"JobDefinition"`
	Parameters    []string  `json:"Parameters"`
	Aws           AwsConfig `json:"aws"`
}

//Load -
func Load(name string) (*Config, error) {
	if setting != nil {
		return setting, nil
	}

	setting = &Config{
		Aws: AwsConfig{},
	}

	fr, err := os.Open(name)
	if err != nil {
		return setting, err
	}
	defer fr.Close()

	err = json.NewDecoder(fr).Decode(&setting)
	if err != nil {
		return setting, err
	}

	return setting, nil
}

//GetAwsCredentials -
func GetAwsCredentials(setting *Config) {
	//Aws Config
	if setting.Aws.AccessKey == "" || setting.Aws.SecretKey == "" {
		setting.Aws.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		setting.Aws.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		setting.Aws.Region = os.Getenv("AWS_DEFAULT_REGION")
	}

	if setting.Aws.AccessKey == "" || setting.Aws.SecretKey == "" {
		homedir, err := util.DirWindows()
		if err != nil {
			homedir = "~"
		} else {
			homedir = filepath.ToSlash(homedir)
		}

		credsfile := fmt.Sprintf("%s/.aws/credentials", homedir)
		creds := credentials.NewSharedCredentials(credsfile, "default")
		credValue, err := creds.Get()
		if err == nil {
			setting.Aws.AccessKey = credValue.AccessKeyID
			setting.Aws.SecretKey = credValue.SecretAccessKey
		}
	}

	//Region
	if setting.Aws.Region == "" {
		setting.Aws.Region = awsDefaultRegion
	}
	//Timeout
	if setting.Aws.Timeout == 0 {
		setting.Aws.Timeout = 120
	}
}
