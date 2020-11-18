package awss

import (
	"net/http"
	"net/url"

	"ditto.co.jp/submit/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

//Service -
type Service struct {
	_config  *config.AwsConfig
	_session *session.Session
}

//NewService -
func NewService(conf *config.AwsConfig) *Service {
	svc := &Service{
		_config: conf,
	}

	sess := svc.getSession()
	if sess == nil {
		return nil
	}

	svc._session = sess

	return svc
}

//NewBatch -
func (s *Service) NewBatch() *batch.Batch {
	return batch.New(s._session)
}

func (s *Service) getSession() *session.Session {
	//Proxy
	var httpClient *http.Client
	if len(s._config.Proxy) > 0 {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(s._config.Proxy)
				},
			},
		}
	}

	//認証情報を作成します。
	cred := credentials.NewStaticCredentials(
		s._config.AccessKey,
		s._config.SecretKey,
		"")

	//セッション作成します
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(s._config.Region),
		Credentials: cred,
		HTTPClient:  httpClient,
	}))

	return sess
}
