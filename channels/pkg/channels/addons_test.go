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

package channels

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
	fakecertmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakekubernetes "k8s.io/client-go/kubernetes/fake"
	"k8s.io/kops/channels/pkg/api"

	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/utils"
)

func Test_Filtering(t *testing.T) {
	grid := []struct {
		Input             api.AddonSpec
		KubernetesVersion string
		Expected          bool
	}{
		{
			Input: api.AddonSpec{
				KubernetesVersion: ">=1.6.0",
			},
			KubernetesVersion: "1.6.0",
			Expected:          true,
		},
		{
			Input: api.AddonSpec{
				KubernetesVersion: "<1.6.0",
			},
			KubernetesVersion: "1.6.0",
			Expected:          false,
		},
		{
			Input: api.AddonSpec{
				KubernetesVersion: ">=1.6.0",
			},
			KubernetesVersion: "1.5.9",
			Expected:          false,
		},
		{
			Input: api.AddonSpec{
				KubernetesVersion: ">=1.4.0 <1.6.0",
			},
			KubernetesVersion: "1.5.9",
			Expected:          true,
		},
		{
			Input: api.AddonSpec{
				KubernetesVersion: ">=1.4.0 <1.6.0",
			},
			KubernetesVersion: "1.6.0",
			Expected:          false,
		},
	}
	for _, g := range grid {
		k8sVersion := semver.MustParse(g.KubernetesVersion)
		addon := &Addon{
			Spec: &g.Input,
		}
		actual := addon.matches(k8sVersion)
		if actual != g.Expected {
			t.Errorf("unexpected result from %v, %s.  got %v", g.Input.KubernetesVersion, g.KubernetesVersion, actual)
		}
	}
}

func Test_Replacement(t *testing.T) {
	grid := []struct {
		Old      *ChannelVersion
		New      *ChannelVersion
		Replaces bool
	}{
		// With no id, update if and only if newer semver
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "", ManifestHash: ""},
			Replaces: false,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.1"), Id: "", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.1"), Id: "", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "", ManifestHash: ""},
			Replaces: false,
		},
		{
			Old:      &ChannelVersion{Version: s("1.1.0"), Id: "", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.1.1"), Id: "", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.1.1"), Id: "", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.1.0"), Id: "", ManifestHash: ""},
			Replaces: false,
		},

		// With id, update if different id and same version, otherwise follow semver
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			Replaces: false,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "b", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "b", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.1"), Id: "a", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.1"), Id: "a", ManifestHash: ""},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.1"), Id: "a", ManifestHash: ""},
			Replaces: true,
		},
		//Test ManifestHash Changes
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: "3544de6578b2b582c0323b15b7b05a28c60b9430"},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: "3544de6578b2b582c0323b15b7b05a28c60b9430"},
			Replaces: false,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: ""},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: "3544de6578b2b582c0323b15b7b05a28c60b9430"},
			Replaces: true,
		},
		{
			Old:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: "3544de6578b2b582c0323b15b7b05a28c60b9430"},
			New:      &ChannelVersion{Version: s("1.0.0"), Id: "a", ManifestHash: "ea9e79bf29adda450446487d65a8fc6b3fdf8c2b"},
			Replaces: true,
		},
	}
	for _, g := range grid {
		actual := g.New.replaces(g.Old)
		if actual != g.Replaces {
			t.Errorf("unexpected result from %v -> %v, expect %t.  actual %v", g.Old, g.New, g.Replaces, actual)
		}
	}
}

func Test_UnparseableVersion(t *testing.T) {
	addons := api.Addons{
		TypeMeta: metav1.TypeMeta{
			Kind: "Addons",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: api.AddonsSpec{
			Addons: []*api.AddonSpec{
				{
					Name:    s("testaddon"),
					Version: s("1.0-kops"),
				},
			},
		},
	}
	bytes, err := utils.YamlMarshal(addons)
	require.NoError(t, err, "marshalling test addons struct")
	location, err := url.Parse("file://testfile")
	require.NoError(t, err, "parsing file url")

	_, err = ParseAddons("test", location, bytes)
	assert.EqualError(t, err, "addon \"testaddon\" has unparseable version \"1.0-kops\": Short version cannot contain PreRelease/Build meta data", "detected invalid version")
}

func Test_MergeAddons(t *testing.T) {
	merges := []struct {
		LeftSide           *AddonMenu
		RightSide          *AddonMenu
		ExpectedAfterMerge *AddonMenu
	}{
		{
			LeftSide:           addonMenu(addon(t, "a", "1.0.0", ">=1.18.0", "k8s-1.18")),
			RightSide:          addonMenu(addon(t, "a", "1.0.1", ">=1.18.0", "k8s-1.18")),
			ExpectedAfterMerge: addonMenu(addon(t, "a", "1.0.1", ">=1.18.0", "k8s-1.18")),
		},
		{
			LeftSide:           addonMenu(addon(t, "a", "1.0.1", ">=1.18.0", "k8s-1.18")),
			RightSide:          addonMenu(addon(t, "a", "1.0.0", ">=1.18.0", "k8s-1.18")),
			ExpectedAfterMerge: addonMenu(addon(t, "a", "1.0.1", ">=1.18.0", "k8s-1.18")),
		},
	}

	for _, m := range merges {
		m.LeftSide.MergeAddons(m.RightSide)
		if !reflect.DeepEqual(m.LeftSide, m.ExpectedAfterMerge) {
			t.Errorf("Unexpected AddonMenu merge result,\nMerged:\n%s\nExpected:\n%s\n", addonMenuString(m.LeftSide), addonMenuString(m.ExpectedAfterMerge))
		}
	}
}

func Test_GetRequiredUpdates(t *testing.T) {
	ctx := context.Background()
	kubeSystem := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-system",
		},
	}
	fakek8s := fakekubernetes.NewSimpleClientset(kubeSystem)
	fakecm := fakecertmanager.NewSimpleClientset()
	addon := &Addon{
		Name: "test",
		Spec: &api.AddonSpec{
			Name:     fi.String("test"),
			NeedsPKI: true,
		},
	}
	addonUpdate, err := addon.GetRequiredUpdates(ctx, fakek8s, fakecm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if addonUpdate == nil {
		t.Fatal("expected addon update, got nil")
	}
	if !addonUpdate.InstallPKI {
		t.Errorf("expected addon to require install")
	}
}

func Test_NeedsRollingUpdate(t *testing.T) {
	grid := []struct {
		newAddon            *Addon
		originalAnnotations map[string]string
		updateRequired      bool
		installRequired     bool
		expectedNodeUpdates int
	}{
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					ManifestHash:       "originalHash",
					Version:            fi.String("1"),
					NeedsRollingUpdate: "all",
				},
			},
		},
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					ManifestHash:       "originalHash",
					Version:            fi.String("2"),
					NeedsRollingUpdate: "all",
				},
			},
			updateRequired:      true,
			expectedNodeUpdates: 2,
		},
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					Version:            fi.String("1"),
					ManifestHash:       "newHash",
					NeedsRollingUpdate: "all",
				},
			},
			updateRequired:      true,
			expectedNodeUpdates: 2,
		},
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					Version:            fi.String("1"),
					ManifestHash:       "newHash",
					NeedsRollingUpdate: "worker",
				},
			},
			updateRequired:      true,
			expectedNodeUpdates: 1,
		},
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					Version:            fi.String("1"),
					ManifestHash:       "newHash",
					NeedsRollingUpdate: "control-plane",
				},
			},
			updateRequired:      true,
			expectedNodeUpdates: 1,
		},
		{
			newAddon: &Addon{
				Name: "test",
				Spec: &api.AddonSpec{
					Name:               fi.String("test"),
					Version:            fi.String("1"),
					ManifestHash:       "newHash",
					NeedsRollingUpdate: "all",
				},
			},
			originalAnnotations: map[string]string{
				"addons.k8s.io/placeholder": "{\"version\":\"1\",\"manifestHash\":\"originalHash\"}",
			},
			installRequired:     true,
			expectedNodeUpdates: 0,
		},
	}

	for _, g := range grid {
		ctx := context.Background()

		annotations := map[string]string{
			"addons.k8s.io/test": "{\"version\":\"1\",\"manifestHash\":\"originalHash\"}",
		}
		if len(g.originalAnnotations) > 0 {
			annotations = g.originalAnnotations
		}

		objects := []runtime.Object{
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "kube-system",
					Annotations: annotations,
				},
			},
			&corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cp",
					Labels: map[string]string{
						"node-role.kubernetes.io/master": "",
					},
				},
			},
			&corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node",
					Labels: map[string]string{
						"node-role.kubernetes.io/node": "",
					},
				},
			},
		}
		fakek8s := fakekubernetes.NewSimpleClientset(objects...)
		fakecm := fakecertmanager.NewSimpleClientset()

		addon := g.newAddon
		required, err := addon.GetRequiredUpdates(ctx, fakek8s, fakecm)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !g.updateRequired && !g.installRequired {
			if required == nil {
				continue
			} else {
				t.Fatalf("did not expect update, but required was not nil")
			}
		}

		if required == nil {
			t.Fatalf("expected required update, got nil")
		}

		if required.NewVersion == nil {
			t.Errorf("updating or installing addon, but NewVersion was nil")
		}

		if required.ExistingVersion != nil {
			if g.installRequired {
				t.Errorf("new install of addon, but ExistingVersion was not nil")
			}
		} else {
			if g.updateRequired {
				t.Errorf("update of addon, but ExistingVersion was nil")
			}
		}

		if err := addon.AddNeedsUpdateLabel(ctx, fakek8s, required); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		nodes, _ := fakek8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		nodeUpdates := 0

		for _, node := range nodes.Items {
			if _, exists := node.Annotations["kops.k8s.io/needs-update"]; exists {
				nodeUpdates++
			}
		}

		if nodeUpdates != g.expectedNodeUpdates {
			t.Errorf("expected %d node updates, but got %d", g.expectedNodeUpdates, nodeUpdates)
		}

	}

}

func Test_InstallPKI(t *testing.T) {
	ctx := context.Background()
	kubeSystem := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-system",
		},
	}
	fakek8s := fakekubernetes.NewSimpleClientset(kubeSystem)
	fakecm := fakecertmanager.NewSimpleClientset()
	addon := &Addon{
		Name: "test",
		Spec: &api.AddonSpec{
			Name:     fi.String("test"),
			NeedsPKI: true,
		},
	}
	err := addon.installPKI(ctx, fakek8s, fakecm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = fakek8s.CoreV1().Secrets("kube-system").Get(ctx, "test-ca", metav1.GetOptions{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	//Two consecutive calls should work since multiple CP nodes can update at the same time
	err = addon.installPKI(ctx, fakek8s, fakecm)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = fakecm.CertmanagerV1().Issuers("kube-system").Get(ctx, "test", metav1.GetOptions{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func s(v string) *string {
	return &v
}

func addon(t *testing.T, name, version, kubernetesVersion, Id string) *Addon {
	location, err := url.Parse("file://x/y/z")
	require.NoError(t, err, "parsing file url")
	return &Addon{
		Name:            name,
		ChannelName:     "test",
		ChannelLocation: *location,
		Spec: &api.AddonSpec{
			Name:              &name,
			Version:           &version,
			KubernetesVersion: kubernetesVersion,
			Id:                Id,
		},
	}
}

func addonMenu(addons ...*Addon) *AddonMenu {
	addonMenu := NewAddonMenu()
	for _, addon := range addons {
		addonMenu.Addons[addon.Name] = addon
	}
	return addonMenu
}

func addonString(addon *Addon) string {
	return fmt.Sprintf("  Addon{Name: %s, Version: %s, KubernetesVersion: %s, Id: %s},\n", *addon.Spec.Name, *addon.Spec.Version, addon.Spec.KubernetesVersion, addon.Spec.Id)
}

func addonMenuString(addonMenu *AddonMenu) string {
	addonMenuString := "AddonMenu{\n"
	for _, addon := range addonMenu.Addons {
		addonMenuString += addonString(addon)
	}
	return addonMenuString + "}\n"
}
