// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stack

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/addons"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/deploy"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/manifest"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/template"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	lbFargateAppTemplateName              = "lb-web-app"
	lbFargateAppParamsPath                = "applications/params.json.tmpl"
	lbFargateAppRulePriorityGeneratorPath = "custom-resources/alb-rule-priority-generator.js"
)

// Parameter logical IDs for a load balanced Fargate service.
const (
	LBFargateProjectNameParamKey       = "ProjectName"
	LBFargateHTTPSParamKey             = "HTTPSEnabled"
	LBFargateEnvNameParamKey           = "EnvName"
	LBFargateAppNameParamKey           = "AppName"
	LBFargateContainerImageParamKey    = "ContainerImage"
	LBFargateContainerPortParamKey     = "ContainerPort"
	LBFargateRulePathParamKey          = "RulePath"
	LBFargateHealthCheckPathParamKey   = "HealthCheckPath"
	LBFargateTaskCPUParamKey           = "TaskCPU"
	LBFargateTaskMemoryParamKey        = "TaskMemory"
	LBFargateTaskCountParamKey         = "TaskCount"
	LBFargateLogRetentionParamKey      = "LogRetention"
	LBFargateAddonsTemplateURLParamKey = "AddonsTemplateURL"
)

type templater interface {
	Template() (string, error)
}

// LBFargateStackConfig represents the configuration needed to create a CloudFormation stack from a
// load balanced Fargate application.
type LBFargateStackConfig struct {
	*deploy.CreateLBFargateAppInput
	httpsEnabled bool
	parser       template.AppTemplateReadParser
	addons       templater
}

// NewLBFargateStack creates a new LBFargateStackConfig from a load-balanced AWS Fargate application.
func NewLBFargateStack(in *deploy.CreateLBFargateAppInput) (*LBFargateStackConfig, error) {
	addons, err := addons.New(in.App.Name)
	if err != nil {
		return nil, fmt.Errorf("new addons: %w", err)
	}
	return &LBFargateStackConfig{
		CreateLBFargateAppInput: in,
		httpsEnabled:            false,
		parser:                  template.New(),
		addons:                  addons,
	}, nil
}

// NewHTTPSLBFargateStack creates a new LBFargateStackConfig from a load-balanced AWS Fargate application. It
// creates an HTTPS listener and assumes that the environment it's being deployed into has an HTTPS configured
// listener.
func NewHTTPSLBFargateStack(in *deploy.CreateLBFargateAppInput) (*LBFargateStackConfig, error) {
	addons, err := addons.New(in.App.Name)
	if err != nil {
		return nil, fmt.Errorf("new addons: %w", err)
	}
	return &LBFargateStackConfig{
		CreateLBFargateAppInput: in,
		httpsEnabled:            true,
		parser:                  template.New(),
		addons:                  addons,
	}, nil
}

// StackName returns the name of the stack.
func (c *LBFargateStackConfig) StackName() string {
	return NameForApp(c.Env.Project, c.Env.Name, c.App.Name)
}

// Template returns the CloudFormation template for the application parametrized for the environment.
func (c *LBFargateStackConfig) Template() (string, error) {
	rulePriorityLambda, err := c.parser.Read(lbFargateAppRulePriorityGeneratorPath)
	if err != nil {
		return "", err
	}
	outputs, err := c.addonsOutputs()
	if err != nil {
		return "", err
	}
	content, err := c.parser.ParseAppTemplate(lbFargateAppTemplateName, struct {
		RulePriorityLambda string
		AddonsOutputs      []addons.Output
		*lbFargateTemplateParams
	}{
		RulePriorityLambda:      rulePriorityLambda.String(),
		AddonsOutputs:           outputs,
		lbFargateTemplateParams: c.toTemplateParams(),
	}, template.WithFuncs(map[string]interface{}{
		"toSnakeCase":           toSnakeCase,
		"filterSecrets":         filterSecrets,
		"filterManagedPolicies": filterManagedPolicies,
	}))
	if err != nil {
		return "", err
	}
	return content.String(), nil
}

// Parameters returns the list of CloudFormation parameters used by the template.
func (c *LBFargateStackConfig) Parameters() []*cloudformation.Parameter {
	templateParams := c.toTemplateParams()
	return []*cloudformation.Parameter{
		{
			ParameterKey:   aws.String(LBFargateProjectNameParamKey),
			ParameterValue: aws.String(templateParams.Env.Project),
		},
		{
			ParameterKey:   aws.String(LBFargateEnvNameParamKey),
			ParameterValue: aws.String(templateParams.Env.Name),
		},
		{
			ParameterKey:   aws.String(LBFargateAppNameParamKey),
			ParameterValue: aws.String(templateParams.App.Name),
		},
		{
			ParameterKey:   aws.String(LBFargateContainerImageParamKey),
			ParameterValue: aws.String(templateParams.Image.URL),
		},
		{
			ParameterKey:   aws.String(LBFargateContainerPortParamKey),
			ParameterValue: aws.String(strconv.FormatUint(uint64(templateParams.Image.Port), 10)),
		},
		{
			ParameterKey:   aws.String(LBFargateRulePathParamKey),
			ParameterValue: aws.String(templateParams.App.Path),
		},
		{
			ParameterKey:   aws.String(LBFargateHealthCheckPathParamKey),
			ParameterValue: aws.String(templateParams.App.HealthCheckPath),
		},
		{
			ParameterKey:   aws.String(LBFargateTaskCPUParamKey),
			ParameterValue: aws.String(strconv.Itoa(templateParams.App.CPU)),
		},
		{
			ParameterKey:   aws.String(LBFargateTaskMemoryParamKey),
			ParameterValue: aws.String(strconv.Itoa(templateParams.App.Memory)),
		},
		{
			ParameterKey:   aws.String(LBFargateTaskCountParamKey),
			ParameterValue: aws.String(strconv.Itoa(templateParams.App.Count)),
		},
		{
			ParameterKey:   aws.String(LBFargateHTTPSParamKey),
			ParameterValue: aws.String(strconv.FormatBool(c.httpsEnabled)),
		},
		{
			ParameterKey:   aws.String(LBFargateLogRetentionParamKey),
			ParameterValue: aws.String("30"),
		},
		{
			ParameterKey:   aws.String(LBFargateAddonsTemplateURLParamKey),
			ParameterValue: aws.String(""),
		},
	}
}

// SerializedParameters returns the CloudFormation stack's parameters serialized
// to a YAML document annotated with comments for readability to users.
func (c *LBFargateStackConfig) SerializedParameters() (string, error) {
	doc, err := c.parser.Parse(lbFargateAppParamsPath, struct {
		Parameters []*cloudformation.Parameter
		Tags       []*cloudformation.Tag
	}{
		Parameters: c.Parameters(),
		Tags:       c.Tags(),
	}, template.WithFuncs(map[string]interface{}{
		"inc": func(i int) int { return i + 1 },
	}))
	if err != nil {
		return "", err
	}
	return doc.String(), nil
}

// Tags returns the list of tags to apply to the CloudFormation stack.
func (c *LBFargateStackConfig) Tags() []*cloudformation.Tag {
	return mergeAndFlattenTags(c.AdditionalTags, map[string]string{
		ProjectTagKey: c.Env.Project,
		EnvTagKey:     c.Env.Name,
		AppTagKey:     c.App.Name,
	})
}

func (c *LBFargateStackConfig) addonsOutputs() ([]addons.Output, error) {
	stack, err := c.addons.Template()
	if err == nil {
		return addons.Outputs(stack)
	}

	var noAddonsErr *addons.ErrDirNotExist
	if !errors.As(err, &noAddonsErr) {
		return nil, fmt.Errorf("generate addons template for application %s: %w", c.App.Name, err)
	}
	return nil, nil // Addons directory does not exist, so there are no outputs and error.
}

// lbFargateTemplateParams holds the data to render the CloudFormation template for an application.
type lbFargateTemplateParams struct {
	*deploy.CreateLBFargateAppInput

	HTTPSEnabled string
	// Field types to override.
	Image struct {
		URL  string
		Port uint16
	}
}

func (c *LBFargateStackConfig) toTemplateParams() *lbFargateTemplateParams {
	url := fmt.Sprintf("%s:%s", c.ImageRepoURL, c.ImageTag)
	return &lbFargateTemplateParams{
		CreateLBFargateAppInput: &deploy.CreateLBFargateAppInput{
			App: &manifest.LoadBalancedWebApp{
				App:                      c.App.App,
				LoadBalancedWebAppConfig: c.CreateLBFargateAppInput.App.ApplyEnv(c.Env.Name), // Get environment specific app configuration.
			},
			Env: c.Env,
		},
		HTTPSEnabled: strconv.FormatBool(c.httpsEnabled),
		Image: struct {
			URL  string
			Port uint16
		}{
			URL:  url,
			Port: c.App.Image.Port,
		},
	}
}
