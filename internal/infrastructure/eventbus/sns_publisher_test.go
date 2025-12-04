package eventbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/eventbus"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSNSClient mock do SNS Client
type MockSNSClient struct {
	mock.Mock
}

func (m *MockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestSNSPublisher_Publish_Success(t *testing.T) {
	mockClient := new(MockSNSClient)
	log := logger.NewNopLogger()
	topicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

	publisher := eventbus.NewSNSPublisher(mockClient, topicARN, log)

	messageID := "message-id-123"
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(
		&sns.PublishOutput{
			MessageId: &messageID,
		},
		nil,
	)

	err := publisher.Publish(context.Background(), "test message", "EVM")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockClient.AssertCalled(t, "Publish", mock.Anything, mock.Anything)
}

func TestSNSPublisher_Publish_Error(t *testing.T) {
	mockClient := new(MockSNSClient)
	log := logger.NewNopLogger()
	topicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

	publisher := eventbus.NewSNSPublisher(mockClient, topicARN, log)

	expectedError := errors.New("SNS service error")
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(nil, expectedError)

	err := publisher.Publish(context.Background(), "test message", "EVM")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish to SNS")
	mockClient.AssertExpectations(t)
}

func TestSNSPublisher_Publish_DifferentChainTypes(t *testing.T) {
	chainTypes := []string{"EVM", "BTC", "TRON", "SOL"}

	for _, chainType := range chainTypes {
		t.Run(chainType, func(t *testing.T) {
			mockClient := new(MockSNSClient)
			log := logger.NewNopLogger()
			topicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

			publisher := eventbus.NewSNSPublisher(mockClient, topicARN, log)

			messageID := "message-id-123"
			mockClient.On("Publish", mock.Anything, mock.MatchedBy(func(input *sns.PublishInput) bool {
				// Verify that chain_type message attribute is set correctly
				if attr, ok := input.MessageAttributes["chain_type"]; ok {
					return *attr.StringValue == chainType
				}
				return false
			})).Return(
				&sns.PublishOutput{
					MessageId: &messageID,
				},
				nil,
			)

			err := publisher.Publish(context.Background(), "test message", chainType)

			assert.NoError(t, err)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestSNSPublisher_Publish_VerifyTopicARN(t *testing.T) {
	mockClient := new(MockSNSClient)
	log := logger.NewNopLogger()
	expectedTopicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

	publisher := eventbus.NewSNSPublisher(mockClient, expectedTopicARN, log)

	messageID := "message-id-123"
	mockClient.On("Publish", mock.Anything, mock.MatchedBy(func(input *sns.PublishInput) bool {
		return *input.TopicArn == expectedTopicARN
	})).Return(
		&sns.PublishOutput{
			MessageId: &messageID,
		},
		nil,
	)

	err := publisher.Publish(context.Background(), "test message", "EVM")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSNSPublisher_Publish_LargeMessage(t *testing.T) {
	mockClient := new(MockSNSClient)
	log := logger.NewNopLogger()
	topicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

	publisher := eventbus.NewSNSPublisher(mockClient, topicARN, log)

	// Create a large message (but still under SNS 256KB limit)
	largePayload := make([]byte, 100000) // 100KB
	for i := range largePayload {
		largePayload[i] = 'A'
	}
	message := string(largePayload)

	messageID := "message-id-123"
	mockClient.On("Publish", mock.Anything, mock.Anything).Return(
		&sns.PublishOutput{
			MessageId: &messageID,
		},
		nil,
	)

	err := publisher.Publish(context.Background(), message, "EVM")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestNewSNSPublisher(t *testing.T) {
	mockClient := new(MockSNSClient)
	log := logger.NewNopLogger()
	topicARN := "arn:aws:sns:us-east-1:123456789012:Transactions"

	publisher := eventbus.NewSNSPublisher(mockClient, topicARN, log)

	assert.NotNil(t, publisher)
}
