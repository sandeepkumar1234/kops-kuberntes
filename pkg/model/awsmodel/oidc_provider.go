/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsmodel

import (
	"k8s.io/kops/pkg/model/iam"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awstasks"
)

// OIDCProviderBuilder configures IAM OIDC Provider
type OIDCProviderBuilder struct {
	*AWSModelContext
	KeyStore  fi.CAStore
	Lifecycle fi.Lifecycle
}

var _ fi.ModelBuilder = &OIDCProviderBuilder{}

const (
	defaultAudience = "amazonaws.com"
)

func (b *OIDCProviderBuilder) Build(c *fi.ModelBuilderContext) error {

	if b.Cluster.Spec.ServiceAccountIssuerDiscovery == nil ||
		!b.Cluster.Spec.ServiceAccountIssuerDiscovery.EnableAWSOIDCProvider {
		return nil
	}

	serviceAccountIssuer, err := iam.ServiceAccountIssuer(&b.Cluster.Spec)
	if err != nil {
		return err
	}

	fingerprints := getFingerprints()

	thumbprints := []*string{}

	for _, fingerprint := range fingerprints {
		thumbprints = append(thumbprints, fi.String(fingerprint))
	}

	c.AddTask(&awstasks.IAMOIDCProvider{
		Name:        fi.String(b.ClusterName()),
		Lifecycle:   b.Lifecycle,
		URL:         fi.String(serviceAccountIssuer),
		ClientIDs:   []*string{fi.String(defaultAudience)},
		Tags:        b.CloudTags(b.ClusterName(), false),
		Thumbprints: thumbprints,
	})

	return nil
}

func getFingerprints() []string {

	//These strings are the sha1 of the two possible S3 root CAs.
	return []string{
		"9e99a48a9960b14926bb7f3b02e22da2b0ab7280",
		"a9d53002e97e00e043244f3d170d6f4c414104fd",
	}

}
