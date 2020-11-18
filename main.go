//go:generate go run gen/buildnumber.go
//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"ditto.co.jp/submit/awss"
	"ditto.co.jp/submit/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/jessevdk/go-flags"
)

//options -
type options struct {
	Version       bool   `short:"v" long:"version" description:"show version"`
	JobName       string `long:"job-name" description:"the name of the job"`
	JobQueue      string `long:"job-queue" description:"the job queue into which the job is submitted"`
	JobDefinition string `long:"job-definition" description:"the job definition used by this job"`
	Parameters    string `long:"parameters" description:"additional parameters. "`
	InputFile     string `long:"cli-input-json" description:"performs service operation based on the JSON"`
	Wait          bool   `long:"wait" description:"wait until job finished"`
	JobID         string `long:"job-id" description:"the id of the job"`
}

func main() {
	//パラメータ取得
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(-1)
	}

	if opts.Version {
		fmt.Printf("DMP submit tool %v\n", version)
		os.Exit(0)
	}

	var err error
	var setting *config.Config

	if opts.InputFile == "" {
		setting = &config.Config{
			JobName:       opts.JobName,
			JobQueue:      opts.JobQueue,
			JobDefinition: opts.JobDefinition,
		}
	} else {
		setting, err = config.Load(opts.InputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	if len(opts.JobID) > 0 {

	} else {
		//input check
		if setting.JobName == "" {
			fmt.Println("please specify the name of the job. --job-name")
			os.Exit(-1)
		}
		if setting.JobQueue == "" {
			fmt.Println("please specify the job queue. --job-queue")
			os.Exit(-1)
		}
		if setting.JobDefinition == "" {
			fmt.Println("please specify the job definition. --job-definition")
			os.Exit(-1)
		}
	}

	config.GetAwsCredentials(setting)
	//check aws credentials
	if setting.Aws.AccessKey == "" || setting.Aws.SecretKey == "" {
		fmt.Println("no valid credentials provided")
		os.Exit(-1)
	}

	//aws service
	svc := awss.NewService(&setting.Aws)
	if svc == nil {
		fmt.Println("failed to create an aws service")
		os.Exit(-1)
	}

	//batch
	b := svc.NewBatch()

	//submit
	if len(opts.JobID) == 0 {
		//submitjob
		input := &batch.SubmitJobInput{
			JobName:       aws.String(setting.JobName),
			JobQueue:      aws.String(setting.JobQueue),
			JobDefinition: aws.String(setting.JobDefinition),
			Parameters:    make(map[string]*string),
		}
		//command parameters
		var k string
		var v string
		for _, p := range setting.Parameters {
			pp := strings.Split(p, "=")
			k = strings.TrimPrefix(pp[0], "--")
			if len(pp) == 2 {
				v = pp[1]
			} else {
				v = ""
			}

			input.Parameters[k] = aws.String(v)
		}

		//command parameters
		if len(opts.Parameters) > 0 {
			pp := strings.Split(opts.Parameters, ",")
			for _, p := range pp {
				ps := strings.Split(p, "=")
				k = strings.TrimPrefix(ps[0], "--")
				if len(ps) == 2 {
					v = ps[1]
				} else {
					v = ""
				}
				input.Parameters[k] = aws.String(v)
			}
		}

		result, err := b.SubmitJob(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case batch.ErrCodeClientException:
					fmt.Println(batch.ErrCodeClientException, aerr.Error())
				case batch.ErrCodeServerException:
					fmt.Println(batch.ErrCodeServerException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			os.Exit(-1)
		}

		fmt.Printf("\r%v: submitted", *result.JobId)
		opts.JobID = *result.JobId
	}

	//describe jobs
	input := &batch.DescribeJobsInput{
		Jobs: []*string{
			aws.String(opts.JobID),
		},
	}

	//describe jobs
	r, err := b.DescribeJobs(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case batch.ErrCodeClientException:
				fmt.Println(batch.ErrCodeClientException, aerr.Error())
			case batch.ErrCodeServerException:
				fmt.Println(batch.ErrCodeServerException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(-1)
	}

	status := *r.Jobs[0].Status
	fmt.Printf("\r%v: %9v", opts.JobID, status)

	if opts.Wait {
		//time
		for range time.Tick(5 * time.Second) {

			r, err := b.DescribeJobs(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case batch.ErrCodeClientException:
						fmt.Println(batch.ErrCodeClientException, aerr.Error())
					case batch.ErrCodeServerException:
						fmt.Println(batch.ErrCodeServerException, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					fmt.Println(err.Error())
				}
				os.Exit(-1)
			}

			status := *r.Jobs[0].Status
			if status == "SUCCEEDED" || status == "FAILED" {
				fmt.Printf("\r%v: %9v", opts.JobID, status)
				//end
				break
			}

			fmt.Printf("\r%v: %9v", opts.JobID, status)
		}
	}

	fmt.Println()
}
