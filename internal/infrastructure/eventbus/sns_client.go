package eventbus

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSClient interface para o cliente SNS (permite mocking)
type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
