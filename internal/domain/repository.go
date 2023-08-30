package domain

// RepositoryBundle is a structure that bundles the repositories of jobs and plugins.
type RepositoryBundle struct {
	Jobs    IJobRepository
	Plugins IPluginRepository
}
