package rollout

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/cmd/rollout"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	kcmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	rolloutLong = templates.LongDesc(`
		Start a new rollout, view its status or history, rollback to a previous revision of your app

		This command allows you to control a deployment config. Each individual rollout is exposed
		as a replication controller, and the deployment process manages scaling down old replication
		controllers and scaling up new ones.

		There are several deployment strategies defined:

		* Rolling (default) - scales up the new replication controller in stages, gradually reducing the
			number of old pods. If one of the new deployed pods never becomes "ready", the new rollout
			will be rolled back (scaled down to zero). Use when your application can tolerate two versions
			of code running at the same time (many web applications, scalable databases)
		* Recreate - scales the old replication controller down to zero, then scales the new replication
			controller up to full. Use when your application cannot tolerate two versions of code running
			at the same time
		* Custom - run your own deployment process inside a Docker container using your own scripts.`)
)

// NewCmdRollout facilitates kubectl rollout subcommands
func NewCmdRollout(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollout SUBCOMMAND",
		Short: "Manage a Kubernetes deployment or OpenShift deployment config",
		Long:  rolloutLong,
		Run:   kcmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	// subcommands
	cmd.AddCommand(NewCmdRolloutHistory(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutPause(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutResume(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutUndo(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutLatest(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutStatus(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutCancel(fullName, f, streams))
	cmd.AddCommand(NewCmdRolloutRetry(fullName, f, streams))

	return cmd
}

var (
	rolloutHistoryLong = templates.LongDesc(`
		View the history of rollouts for a specific deployment config

		You can also view more detailed information for a specific revision
		by using the --revision flag.`)

	rolloutHistoryExample = templates.Examples(`
		# View the rollout history of a deployment
	  %[1]s rollout history dc/nginx

	  # View the details of deployment revision 3
	  %[1]s rollout history dc/nginx --revision=3`)
)

// NewCmdRolloutHistory is a wrapper for the Kubernetes cli rollout history command
func NewCmdRolloutHistory(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := rollout.NewCmdRolloutHistory(f, streams.Out)
	cmd.Long = rolloutHistoryLong
	cmd.Example = fmt.Sprintf(rolloutHistoryExample, fullName)
	cmd.ValidArgs = append(cmd.ValidArgs, "deploymentconfig")
	return cmd
}

var (
	rolloutPauseLong = templates.LongDesc(`
    Mark the provided resource as paused

    Paused resources will not be reconciled by a controller.
    Use \"%[1]s rollout resume\" to resume a paused resource.`)

	rolloutPauseExample = templates.Examples(`
    # Mark the nginx deployment as paused. Any current state of
    # the deployment will continue its function, new updates to the deployment will not
    # have an effect as long as the deployment is paused.
    %[1]s rollout pause dc/nginx`)
)

// NewCmdRolloutPause is a wrapper for the Kubernetes cli rollout pause command
func NewCmdRolloutPause(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := rollout.NewCmdRolloutPause(f, streams)
	cmd.Long = rolloutPauseLong
	cmd.Example = fmt.Sprintf(rolloutPauseExample, fullName)
	cmd.ValidArgs = append(cmd.ValidArgs, "deploymentconfig")
	return cmd
}

var (
	rolloutResumeLong = templates.LongDesc(`
    Resume a paused resource

    Paused resources will not be reconciled by a controller. By resuming a
    resource, we allow it to be reconciled again.`)

	rolloutResumeExample = templates.Examples(`
    # Resume an already paused deployment
    %[1]s rollout resume dc/nginx`)
)

// NewCmdRolloutResume is a wrapper for the Kubernetes cli rollout resume command
func NewCmdRolloutResume(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := rollout.NewCmdRolloutResume(f, streams)
	cmd.Long = rolloutResumeLong
	cmd.Example = fmt.Sprintf(rolloutResumeExample, fullName)
	cmd.ValidArgs = append(cmd.ValidArgs, "deploymentconfig")
	return cmd
}

var (
	rolloutUndoLong = templates.LongDesc(`
    Revert an application back to a previous deployment

    When you run this command your deployment configuration will be updated to
    match a previous deployment. By default only the pod and container
    configuration will be changed and scaling or trigger settings will be left as-
    is. Note that environment variables and volumes are included in rollbacks, so
    if you've recently updated security credentials in your environment your
    previous deployment may not have the correct values.

    Any image triggers present in the rolled back configuration will be disabled
    with a warning. This is to help prevent your rolled back deployment from being
    replaced by a triggered deployment soon after your rollback. To re-enable the
    triggers, use the 'oc set triggers --auto' command.

    If you would like to review the outcome of the rollback, pass '--dry-run' to print
    a human-readable representation of the updated deployment configuration instead of
    executing the rollback. This is useful if you're not quite sure what the outcome
    will be.`)

	rolloutUndoExample = templates.Examples(`
    # Rollback to the previous deployment
    %[1]s rollout undo dc/nginx

    # Rollback to deployment revision 3. The replication controller for that version must exist.
    %[1]s rollout undo dc/nginx --to-revision=3`)
)

// NewCmdRolloutUndo is a wrapper for the Kubernetes cli rollout undo command
func NewCmdRolloutUndo(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := rollout.NewCmdRolloutUndo(f, streams.Out)
	cmd.Long = rolloutUndoLong
	cmd.Example = fmt.Sprintf(rolloutUndoExample, fullName)
	cmd.ValidArgs = append(cmd.ValidArgs, "deploymentconfig")
	return cmd
}

var (
	rolloutStatusLong = templates.LongDesc(`
		Watch the status of the latest rollout, until it's done.`)

	rolloutStatusExample = templates.Examples(`
		# Watch the status of the latest rollout
  	%[1]s rollout status dc/nginx`)
)

// NewCmdRolloutStatus is a wrapper for the Kubernetes cli rollout status command
func NewCmdRolloutStatus(fullName string, f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := rollout.NewCmdRolloutStatus(f, streams)
	cmd.Long = rolloutStatusLong
	cmd.Example = fmt.Sprintf(rolloutStatusExample, fullName)
	return cmd
}