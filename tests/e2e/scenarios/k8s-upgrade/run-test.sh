#!/usr/bin/env bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REPO_ROOT=$(git rev-parse --show-toplevel);
source "${REPO_ROOT}"/tests/e2e/scenarios/lib/common.sh

PATH=$REPO_ROOT/bazel-bin/cmd/kops/$(go env GOOS)-$(go env GOARCH):$PATH
KOPS=${REPO_ROOT}/bazel-bin/cmd/kops/linux-amd64/kops

${KUBETEST2} --build --kops-root="${REPO_ROOT}" --kops-binary-path="${KOPS}" --stage-location="${STAGE_LOCATION:-}"

${KUBETEST2} \
	--up \
	--kops-binary-path="${KOPS}" \
	--kubernetes-version=v1.19.10 \
	--create-args="--networking calico"

kops set cluster "${CLUSTER_NAME}" cluster.spec.kubernetesVersion=v1.20.6
kops update cluster
kops update cluster --admin --yes
kops update cluster

kops rolling-update cluster
kops rolling-update cluster --yes --validation-timeout 30m

kops validate cluster

${KUBETEST2} \
		--cloud-provider="${CLOUD_PROVIDER}" \
		--kops-binary-path="${KOPS}" \
		--test=kops \
		-- \
		--test-package-version=v1.20.6 \
		--parallel 25 \
		--skip-regex="\[Slow\]|\[Serial\]|\[Disruptive\]|\[Flaky\]|\[Feature:.+\]|\[HPA\]|Dashboard|RuntimeClass|RuntimeHandler|TCP.CLOSE_WAIT|Projected.configMap.optional.updates"
