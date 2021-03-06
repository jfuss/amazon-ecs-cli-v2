// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/aws/codepipeline"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/aws/session"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/manifest"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/store"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/term/color"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/term/log"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/workspace"
	"github.com/spf13/cobra"
)

const (
	pipelineShowProjectNamePrompt     = "Which project's pipelines would you like to show?"
	pipelineShowProjectNameHelpPrompt = "A project groups all of your pipelines together."
)

type showPipelineVars struct {
	*GlobalOpts
	shouldOutputJSON      bool
	shouldOutputResources bool
	pipelineName          string
}

type showPipelineOpts struct {
	showPipelineVars

	// Interfaces to dependencies
	ws          wsPipelineReader
	store       storeReader
	pipelineSvc pipelineGetter
}

func newShowPipelineOpts(vars showPipelineVars) (*showPipelineOpts, error) {
	ssmStore, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("connect to environment datastore: %w", err)
	}

	ws, err := workspace.New()
	if err != nil {
		return nil, fmt.Errorf("workspace cannot be created: %w", err)
	}

	p := session.NewProvider()
	defaultSession, err := p.Default()
	if err != nil {
		return nil, err
	}

	opts := &showPipelineOpts{
		showPipelineVars: vars,
		ws:               ws,
		store:            ssmStore,
		pipelineSvc:      codepipeline.New(defaultSession),
	}

	return opts, nil
}

// Validate returns an error if the flag values passed by the user are invalid.
func (o *showPipelineOpts) Validate() error {
	if o.ProjectName() != "" {
		if _, err := o.store.GetProject(o.ProjectName()); err != nil {
			return err
		}
	}
	if o.pipelineName != "" {
		if _, err := o.pipelineSvc.GetPipeline(o.pipelineName); err != nil {
			return err
		}
	}

	return nil
}

// Ask prompts for fields that are required but not passed in.
func (o *showPipelineOpts) Ask() error {
	if err := o.askProject(); err != nil {
		return err
	}

	return o.askPipelineName()
}

func (o *showPipelineOpts) askProject() error {
	if o.ProjectName() != "" {
		return nil
	}
	projNames, err := o.retrieveProjects()
	if err != nil {
		return err
	}
	if len(projNames) == 0 {
		return fmt.Errorf("no project found: run %s please", color.HighlightCode("project init"))
	}
	proj, err := o.prompt.SelectOne(
		pipelineShowProjectNamePrompt,
		pipelineShowProjectNameHelpPrompt,
		projNames,
	)
	if err != nil {
		return fmt.Errorf("select projects: %w", err)
	}
	o.projectName = proj

	return nil
}

func (o *showPipelineOpts) retrieveProjects() ([]string, error) {
	projs, err := o.store.ListProjects()
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	projNames := make([]string, len(projs))
	for ind, proj := range projs {
		projNames[ind] = proj.Name
	}
	return projNames, nil
}

func (o *showPipelineOpts) askPipelineName() error {
	// return if pipeline name is set by flag
	if o.pipelineName != "" {
		return nil
	}

	pipelineName, err := o.getPipelineNameFromManifest()
	if err != nil {
		return err
	}

	o.pipelineName = pipelineName
	// if pipelineName == "" {
	// get pipeline name from remote
	// }

	return nil
}

func (o *showPipelineOpts) getPipelineNameFromManifest() (string, error) {
	// try to read pipelinename from from manifest

	data, err := o.ws.ReadPipelineManifest()
	if err != nil {
		log.Infof("No pipeline manifest in workspace for project %s\n.", color.HighlightUserInput(o.ProjectName()))
		return "", nil
	}

	pipeline, err := manifest.UnmarshalPipeline(data)
	if err != nil {
		return "", fmt.Errorf("unmarshal pipeline manifest: %w", err)
	}

	return pipeline.Name, nil
}

// Execute writes the pipeline manifest file.
func (o *showPipelineOpts) Execute() error {
	// TODO Placeholder
	return nil
}

// BuildPipelineShowCmd build the command for deploying a new pipeline or updating an existing pipeline.
func BuildPipelineShowCmd() *cobra.Command {
	vars := showPipelineVars{
		GlobalOpts: NewGlobalOpts(),
	}
	cmd := &cobra.Command{
		Hidden: true, // TODO remove when ready for production!
		Use:    "show",
		Short:  "Shows info about a deployed pipeline for a project.",
		Long:   "Shows info about a deployed pipeline for a project, including information about each stage.",
		Example: `
  Shows info about the pipeline pipeline-myproject-mycompany-myrepo"
  /code $ ecs-preview pipeline show --project myproject --resources`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newShowPipelineOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			if err := opts.Ask(); err != nil {
				return err
			}
			return opts.Execute()
		}),
	}
	cmd.Flags().StringVarP(&vars.pipelineName, nameFlag, nameFlagShort, "", pipelineFlagDescription)
	cmd.Flags().StringVarP(&vars.projectName, projectFlag, projectFlagShort, "", projectFlagDescription)
	cmd.Flags().BoolVar(&vars.shouldOutputJSON, jsonFlag, false, jsonFlagDescription)
	cmd.Flags().BoolVar(&vars.shouldOutputResources, resourcesFlag, false, resourcesFlagDescription)

	return cmd
}
