// A generated module for HelloDagger functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/hello-dagger/internal/dagger"
	"fmt"
	"math"
	"math/rand/v2"
)

type HelloDagger struct{}

// Develop with a Github issue as the assignment and open a pull request
func (m *HelloDagger) DevelopIssue(
	ctx context.Context,
	// Github token with permissions to write issues and contents
	githubToken *dagger.Secret,
	// Github issue number
	issueID int,
	// Github repository url
	repository string,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	// Get the github issue
	issueClient := dag.GithubIssue(dagger.GithubIssueOpts{Token: githubToken})
	issue := issueClient.Read(repository, issueID)

	// Get information from the Github issue
	assignment, err := issue.Body(ctx)
	if err != nil {
		return "", err
	}

	// Solve the issue with the Develop Agent
	feature, err := m.Develop(ctx, assignment, source)
	if err != nil {
		return "", err
	}

	// Open a pull request
	title, err := issue.Title(ctx)
	if err != nil {
		return "", err
	}
	url, err := issue.URL(ctx)
	if err != nil {
		return "", err
	}
	body := assignment + "\n\nCloses " + url
	pr := issueClient.CreatePullRequest(repository, title, body, feature)

	return pr.URL(ctx)
}

func (m *HelloDagger) Develop(
	ctx context.Context,
	// Assignment to complete
	assignment string,
	// +defaultPath="/"
	source *dagger.Directory,
) (*dagger.Directory, error) {
	// Environment with agent inputs and outputs
	environment := dag.Env().
		WithStringInput("assignment", assignment, "the assignment to complete").
		WithWorkspaceInput(
			"workspace",
			dag.Workspace(source),
			"the workspace with tools to edit and test code",
		).
		WithWorkspaceOutput(
			"completed",
			"the workspace with the completed assignment",
		)

	promptFile := dag.CurrentModule().Source().File("develop_prompt.md")

	work := dag.LLM().
		WithEnv(environment).
		WithPromptFile(promptFile)

	completed := work.
		Env().
		Output("completed").
		AsWorkspace()
	completedDirectory := completed.Source().WithoutDirectory("node_modules")

	_, err := m.Test(ctx, completedDirectory)
	if err != nil {
		return nil, err
	}

	return completedDirectory, nil
}

func (m *HelloDagger) Publish(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.Build(source).
		Publish(ctx, fmt.Sprintf("ttl.sh/hello-dagger-%.0f", math.Floor(rand.Float64()*1000000))) //#nosec

}

func (m *HelloDagger) Build(
	// +defaultPath="/"
	source *dagger.Directory) *dagger.Container {
	build := m.BuildEnv(source).
		WithExec([]string{"npm", "run", "build"}).
		Directory("./dist")

	return dag.Container().From("nginx:1.25-alpine").
		WithDirectory("/usr/share/nginx/html", build).
		WithExposedPort(80)
}

func (m *HelloDagger) Test(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"npm", "run", "test:unit", "run"}).
		Stdout(ctx)
}

func (m *HelloDagger) BuildEnv(source *dagger.Directory) *dagger.Container {
	nodeCache := dag.CacheVolume("node")
	return dag.Container().From("node:21-slim").
		WithMountedCache("/root/.npm", nodeCache).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"npm", "install"})
}
