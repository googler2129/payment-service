package fetcher

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	"github.com/mercor/payment-service/constants"
	"github.com/mercor/payment-service/pkg/config/model"
)

type CloudFetcher struct {
	configData *ConfigData
	client     *appconfigdata.Client
}

type ConfigData struct {
	sessionToken *string
	profileName  string
	yamlContent  string
}

func NewCloudFetcher(ctx context.Context, profileName, applicationName string) (Fetcher, error) {
	awsConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := appconfigdata.NewFromConfig(awsConfig)

	fetcher := &CloudFetcher{
		client:     client,
		configData: &ConfigData{profileName: profileName},
	}

	if err = fetcher.initiateSession(ctx, applicationName); err != nil {
		return nil, err
	}

	return fetcher, nil
}

func (fetcher *CloudFetcher) GetConfig(ctx context.Context) (*model.Config, error) {
	if err := fetcher.updateConfigData(ctx); err != nil {
		return nil, err
	}

	return fetcher.constructConfig()
}

func (fetcher *CloudFetcher) initiateSession(ctx context.Context, applicationName string) error {
	environment := constants.DeployedEnv
	output, err := fetcher.client.StartConfigurationSession(ctx, &appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:          &applicationName,
		ConfigurationProfileIdentifier: &fetcher.configData.profileName,
		EnvironmentIdentifier:          &environment,
	})
	if err != nil {
		return fmt.Errorf("failed to initiate configuration session: %w", err)
	}

	fetcher.configData.sessionToken = output.InitialConfigurationToken
	return nil
}

func (fetcher *CloudFetcher) updateConfigData(ctx context.Context) error {
	configuration, err := fetcher.fetchLatestConfiguration(ctx)
	if err != nil {
		return err
	}

	if configuration != "" {
		fetcher.configData.yamlContent = configuration
	}

	return nil
}

func (fetcher *CloudFetcher) fetchLatestConfiguration(ctx context.Context) (string, error) {
	output, err := fetcher.client.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: fetcher.configData.sessionToken,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest configuration: %w", err)
	}

	fetcher.configData.sessionToken = output.NextPollConfigurationToken
	return string(output.Configuration), nil
}

func (fetcher *CloudFetcher) constructConfig() (*model.Config, error) {
	configMap, err := ParseYAMLToConfigMap(fetcher.configData.yamlContent)
	if err != nil {
		return nil, err
	}

	return model.NewConfig(configMap), nil
}
