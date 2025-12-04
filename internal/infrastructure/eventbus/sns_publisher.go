package eventbus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"go.uber.org/zap"
)

// SNSPublisher publica eventos no AWS SNS
type SNSPublisher struct {
	snsClient SNSClient
	topicARN  string
	logger    *zap.Logger
}

// NewSNSPublisher cria um novo SNS publisher
func NewSNSPublisher(snsClient SNSClient, topicARN string, logger *zap.Logger) *SNSPublisher {
	return &SNSPublisher{
		snsClient: snsClient,
		topicARN:  topicARN,
		logger:    logger,
	}
}

// Publish publica uma mensagem no SNS Topic "Transactions"
func (p *SNSPublisher) Publish(ctx context.Context, message string, chainType string) error {
	input := &sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(message),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"chain_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(chainType),
			},
		},
	}

	result, err := p.snsClient.Publish(ctx, input)
	if err != nil {
		p.logger.Error("failed to publish to SNS",
			zap.Error(err),
			zap.String("chain_type", chainType),
		)
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	p.logger.Info("message published to SNS",
		zap.String("message_id", *result.MessageId),
		zap.String("chain_type", chainType),
	)

	return nil
}
