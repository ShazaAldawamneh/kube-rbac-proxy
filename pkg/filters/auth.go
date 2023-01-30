/*
Copyright 2022 the kube-rbac-proxy maintainers All rights reserved.

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
package filters

import (
	"net/http"
	"strings"

	"github.com/brancz/kube-rbac-proxy/pkg/authn"

	"k8s.io/apiserver/pkg/endpoints/request"
)

func UnauthorizedHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// WithAuthHeaders adds identity information to the headers.
// Must not be used, if connection is not encrypted with TLS.
func WithAuthHeaders(handler http.Handler, cfg *authn.AuthnHeaderConfig) http.Handler {
	upstreamHeadersEnabled := len(cfg.GroupsFieldName) > 0 || len(cfg.UserFieldName) > 0
	if !upstreamHeadersEnabled {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		u, ok := request.UserFrom(req.Context())
		if ok {
			// Seemingly well-known headers to tell the upstream about user's identity
			// so that the upstream can achieve the original goal of delegating RBAC authn/authz to kube-rbac-proxy
			req.Header.Set(cfg.UserFieldName, u.GetName())
			req.Header.Set(cfg.GroupsFieldName, strings.Join(u.GetGroups(), cfg.GroupSeparator))
		}

		handler.ServeHTTP(w, req)
	})
}
