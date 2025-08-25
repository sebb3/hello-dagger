// A generated module for Workspace functions
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
	"dagger/workspace/internal/dagger"
)

type Workspace struct {
	Source *dagger.Directory
}

func New(
	// The source directory
	source *dagger.Directory,
) *Workspace {
	return &Workspace{Source: source}
}

// Read a file in a workspace.
func (w *Workspace) ReadFile(
	ctx context.Context,
	// The path to the file in the workspace
	path string,
) (string, error) {
	return w.Source.File(path).Contents(ctx)
}

// Write a file in a workspace.
func (w *Workspace) WriteFile(
	ctx context.Context,
	// The path to the file in the workspace
	path string,
	// The contents of the file
	content string,
) *Workspace {
	w.Source = w.Source.WithNewFile(path, content)
	return w
}

// List all of the files in the workspace.
func (w *Workspace) ListFiles(
	ctx context.Context,
) (string, error) {
	return dag.Container().
		From("alpine:3").
		WithDirectory("/src", w.Source).
		WithWorkdir("/src").
		WithExec([]string{"tree", "./src"}).
		Stdout(ctx)
}

// Return the result of running unit tests
func (w *Workspace) Test(ctx context.Context) (string, error) {
	nodeCache := dag.CacheVolume("node")
	return dag.Container().
		From("node:21-slim").
		WithDirectory("/src", w.Source).
		WithMountedCache("/root/.npm", nodeCache).
		WithWorkdir("/src").
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "test:unit", "run"}).
		Stdout(ctx)
}
