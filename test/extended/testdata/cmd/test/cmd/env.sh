#!/bin/bash
source "$(dirname "${BASH_SOURCE}")/../../hack/lib/init.sh"
trap os::test::junit::reconcile_output EXIT

os::test::junit::declare_suite_start "cmd/set-env"
# This test validates the value of --image for oc run
os::cmd::expect_success 'oc create deploymentconfig testdc --image=busybox'
os::cmd::expect_failure_and_text 'oc set env dc/testdc' 'error: at least one environment variable must be provided'
os::cmd::expect_success_and_text 'oc set env dc/testdc key=value' 'deploymentconfig.apps.openshift.io/testdc updated'
os::cmd::expect_success_and_text 'oc set env dc/testdc --list' 'deploymentconfigs/testdc, container default-container'
os::cmd::expect_success_and_text 'oc set env dc --all --containers="default-container" key-' 'deploymentconfig.apps.openshift.io/testdc updated'
os::cmd::expect_failure_and_text 'oc set env dc --all --containers="default-container"' 'error: at least one environment variable must be provided'
os::cmd::expect_failure_and_not_text 'oc set env --from=secret/mysecret dc/testdc' 'error: at least one environment variable must be provided'
os::cmd::expect_failure_and_text 'oc set env dc/testdc test#abc=1234' 'environment variable test#abc=1234 is invalid, a valid environment variable name must consist of alphabetic characters'

# ensure deleting a var through --env does not result in an error message
os::cmd::expect_success_and_text 'oc set env dc/testdc key=value' 'deploymentconfig.apps.openshift.io/testdc updated'
os::cmd::expect_success_and_text 'oc set env dc/testdc dots.in.a.key=dots.in.a.value' 'deploymentconfig.apps.openshift.io/testdc updated'
os::cmd::expect_success_and_text 'oc set env dc --all --containers="default-container" --env=key-' 'deploymentconfig.apps.openshift.io/testdc updated'
# ensure deleting a var through --env actually deletes the env var
os::cmd::expect_success_and_not_text "oc get dc/testdc -o jsonpath='{ .spec.template.spec.containers[?(@.name==\"default-container\")].env }'" 'name.?\:.?key'
os::cmd::expect_success_and_text "oc get dc/testdc -o jsonpath='{ .spec.template.spec.containers[?(@.name==\"default-container\")].env }'" 'name.?\:.?dots.in.a.key'
os::cmd::expect_success_and_text 'oc set env dc --all --containers="default-container" --env=dots.in.a.key-' 'deploymentconfig.apps.openshift.io/testdc updated'
os::cmd::expect_success_and_not_text "oc get dc/testdc -o jsonpath='{ .spec.template.spec.containers[?(@.name==\"default-container\")].env }'" 'name.?\:.?dots.in.a.key'

# check that env vars are not split at commas
os::cmd::expect_success_and_text 'oc set env -o yaml dc/testdc PASS=x,y=z' 'value: x,y=z'
os::cmd::expect_success_and_text 'oc set env -o yaml dc/testdc --env PASS=x,y=z' 'value: x,y=z'
# warning is printed when --env has comma in it
os::cmd::expect_success_and_text 'oc set env dc/testdc --env PASS=x,y=z' 'no longer accepts comma-separated list'
# warning is not printed for variables passed as positional arguments
os::cmd::expect_success_and_not_text 'oc set env dc/testdc PASS=x,y=z' 'no longer accepts comma-separated list'

# create a build-config object with the JenkinsPipeline strategy
os::cmd::expect_success 'oc process -p NAMESPACE=openshift -f ${TEST_DATA}/jenkins/jenkins-ephemeral-template.json | oc create -f -'
os::cmd::expect_success "echo 'apiVersion: build.openshift.io/v1
kind: BuildConfig
metadata:
  name: fake-pipeline
spec:
  source:
    git:
      uri: https://github.com/openshift/ruby-hello-world.git
  strategy:
    jenkinsPipelineStrategy: {}
' | oc create -f -"

# ensure build-config has been created and that its type is "JenkinsPipeline"
os::cmd::expect_success_and_text "oc get bc fake-pipeline -o jsonpath='{ .spec.strategy.type }'" 'JenkinsPipeline'
# attempt to set an environment variable
os::cmd::expect_success_and_text 'oc set env bc/fake-pipeline FOO=BAR' 'buildconfig.build.openshift.io/fake\-pipeline updated'
# ensure environment variable was set
os::cmd::expect_success_and_text "oc get bc fake-pipeline -o jsonpath='{ .spec.strategy.jenkinsPipelineStrategy.env }'" 'name.?\:.?FOO'
os::cmd::expect_success 'oc delete bc fake-pipeline'

echo "oc set env: ok"
os::test::junit::declare_suite_end
