// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cloudwatch

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/aws/cloudwatch/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCloudWatch_GetAlarms(t *testing.T) {
	const (
		appName     = "mockApp"
		envName     = "mockEnv"
		projectName = "mockProject"
	)
	mockTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+00:00")
	mockError := errors.New("some error")
	testCases := map[string]struct {
		mockcwClient func(m *mocks.MockcloudWatchClient)
		mockrgClient func(m *mocks.MockresourceGroupClient)

		wantErr         error
		wantAlarmStatus []AlarmStatus
	}{
		"errors if failed to search resources": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(nil, mockError)
			},

			wantErr: fmt.Errorf("search CloudWatch alarm resources: some error"),
		},
		"errors if failed to get alarm names because of invalid ARN": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("badArn"),
						},
					},
				}, nil)
			},

			wantErr: fmt.Errorf("parse alarm ARN badArn: arn: invalid prefix"),
		},
		"errors if failed to get alarm names because of bad ARN resource": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:badAlarm:Names"),
						},
					},
				}, nil)
			},

			wantErr: fmt.Errorf("cannot parse alarm ARN resource alarm:badAlarm:Names"),
		},
		"errors if failed to describe CloudWatch alarms": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {
				m.EXPECT().DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
					NextToken:  nil,
					AlarmNames: aws.StringSlice([]string{"mockAlarmName"}),
				}).Return(nil, mockError)
			},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName"),
						},
					},
				}, nil)
			},

			wantErr: fmt.Errorf("describe CloudWatch alarms: some error"),
		},
		"return an empty array if no alarms found": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken:           nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{},
				}, nil)
			},

			wantAlarmStatus: []AlarmStatus{},
		},
		"success": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {
				m.EXPECT().DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
					NextToken:  nil,
					AlarmNames: aws.StringSlice([]string{"mockAlarmName"}),
				}).Return(&cloudwatch.DescribeAlarmsOutput{
					NextToken: nil,
					MetricAlarms: []*cloudwatch.MetricAlarm{
						{
							AlarmArn:              aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName"),
							AlarmName:             aws.String("mockAlarmName"),
							StateReason:           aws.String("mockReason"),
							StateValue:            aws.String("mockState"),
							StateUpdatedTimestamp: &mockTime,
						},
					},
				}, nil)
			},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName"),
						},
					},
				}, nil)
			},

			wantAlarmStatus: []AlarmStatus{
				{
					Arn:          "arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName",
					Name:         "mockAlarmName",
					Type:         "Metric",
					Reason:       "mockReason",
					Status:       "mockState",
					UpdatedTimes: mockTime.Unix(),
				},
			},
		},
		"success with pagination": {
			mockcwClient: func(m *mocks.MockcloudWatchClient) {
				m.EXPECT().DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
					NextToken:  nil,
					AlarmNames: aws.StringSlice([]string{"mockAlarmName1", "mockAlarmName2"}),
				}).Return(&cloudwatch.DescribeAlarmsOutput{
					NextToken: aws.String("mockNextToken"),
					CompositeAlarms: []*cloudwatch.CompositeAlarm{
						{
							AlarmArn:              aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName1"),
							AlarmName:             aws.String("mockAlarmName1"),
							StateReason:           aws.String("mockReason"),
							StateValue:            aws.String("mockState"),
							StateUpdatedTimestamp: &mockTime,
						},
					},
				}, nil)
				m.EXPECT().DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
					NextToken:  aws.String("mockNextToken"),
					AlarmNames: aws.StringSlice([]string{"mockAlarmName1", "mockAlarmName2"}),
				}).Return(&cloudwatch.DescribeAlarmsOutput{
					NextToken: nil,
					MetricAlarms: []*cloudwatch.MetricAlarm{
						{
							AlarmArn:              aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName2"),
							AlarmName:             aws.String("mockAlarmName2"),
							StateReason:           aws.String("mockReason"),
							StateValue:            aws.String("mockState"),
							StateUpdatedTimestamp: &mockTime,
						},
					},
				}, nil)
			},
			mockrgClient: func(m *mocks.MockresourceGroupClient) {
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: nil,
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: aws.String("mockNextToken"),
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName1"),
						},
					},
				}, nil)
				m.EXPECT().SearchResources(&resourcegroups.SearchResourcesInput{
					NextToken: aws.String("mockNextToken"),
					ResourceQuery: &resourcegroups.ResourceQuery{
						Type:  aws.String("TAG_FILTERS_1_0"),
						Query: aws.String(`{"ResourceTypeFilters":["AWS::CloudWatch::Alarm"],"TagFilters":[{"Key":"ecs-project","Values":["mockProject"]}]}`),
					},
				}).Return(&resourcegroups.SearchResourcesOutput{
					NextToken: nil,
					ResourceIdentifiers: []*resourcegroups.ResourceIdentifier{
						{
							ResourceArn: aws.String("arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName2"),
						},
					},
				}, nil)
			},

			wantAlarmStatus: []AlarmStatus{
				{
					Arn:          "arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName1",
					Name:         "mockAlarmName1",
					Type:         "Composite",
					Reason:       "mockReason",
					Status:       "mockState",
					UpdatedTimes: mockTime.Unix(),
				},
				{
					Arn:          "arn:aws:cloudwatch:us-west-2:1234567890:alarm:mockAlarmName2",
					Name:         "mockAlarmName2",
					Type:         "Metric",
					Reason:       "mockReason",
					Status:       "mockState",
					UpdatedTimes: mockTime.Unix(),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockcwClient := mocks.NewMockcloudWatchClient(ctrl)
			mockresourceGroupClient := mocks.NewMockresourceGroupClient(ctrl)
			tc.mockcwClient(mockcwClient)
			tc.mockrgClient(mockresourceGroupClient)

			cwSvc := CloudWatch{
				mockcwClient,
				mockresourceGroupClient,
			}

			gotAlarmStatus, gotErr := cwSvc.GetAlarmsWithTags(map[string]string{
				"ecs-project": projectName,
			})

			if gotErr != nil {
				require.EqualError(t, tc.wantErr, gotErr.Error())
			} else {
				require.ElementsMatch(t, tc.wantAlarmStatus, gotAlarmStatus)
			}
		})

	}
}
