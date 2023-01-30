/*
Copyright 2022 the kube-rbac-proxy maintainers. All rights reserved.

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

package options

import (
	"github.com/brancz/kube-rbac-proxy/pkg/authn"
	"github.com/brancz/kube-rbac-proxy/pkg/authn/identityheaders"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	kubeflags "k8s.io/component-base/cli/flag"
)

// ProxyRunOptions bundles both generic server run options from upstream, the
// proxy-specific options and legacy options
type ProxyRunOptions struct {
	SecureServing *genericoptions.SecureServingOptions
	// ProxySecureServing are options for the proxy endpoints, they will be copied
	// from the above with a changed port
	ProxySecureServing *genericoptions.SecureServingOptions

	DelegatingAuthentication *genericoptions.DelegatingAuthenticationOptions
	DelegatingAuthorization  *genericoptions.DelegatingAuthorizationOptions

	ProxyOptions  *ProxyOptions
	LegacyOptions *LegacyOptions
}

func NewProxyRunOptions() *ProxyRunOptions {
	return &ProxyRunOptions{
		SecureServing:            genericoptions.NewSecureServingOptions(),
		DelegatingAuthentication: genericoptions.NewDelegatingAuthenticationOptions(),
		DelegatingAuthorization:  genericoptions.NewDelegatingAuthorizationOptions(),

		ProxyOptions: &ProxyOptions{
			UpstreamHeader: &identityheaders.AuthnHeaderConfig{},
			OIDC:           &authn.OIDCConfig{},
		},
		LegacyOptions: &LegacyOptions{
			x509Auth: &authn.X509Config{},
		},
	}
}

func (o *ProxyRunOptions) Flags() kubeflags.NamedFlagSets {
	namedFlagSets := kubeflags.NamedFlagSets{}

	o.SecureServing.AddFlags(namedFlagSets.FlagSet("secure serving"))
	o.DelegatingAuthentication.AddFlags(namedFlagSets.FlagSet("delegating authentication"))
	o.DelegatingAuthorization.AddFlags(namedFlagSets.FlagSet("delegating authorization"))
	o.ProxyOptions.AddFlags(namedFlagSets.FlagSet("proxy"))
	o.LegacyOptions.AddFlags(namedFlagSets.FlagSet("legacy kube-rbac-proxy [DEPRECATED]"))

	return namedFlagSets
}
