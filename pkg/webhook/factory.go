// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhook

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type Webhook struct {
	Name     string
	Kind     string
	Provider string
	Path     string
	Types    []runtime.Object
	Webhook  *admission.Webhook
}

// FactoryAggregator aggregates various Factory functions.
type FactoryAggregator []func(manager.Manager) (*Webhook, error)

// NewFactoryAggregator creates a new FactoryAggregator and registers the given functions.
func NewFactoryAggregator(m []func(manager.Manager) (*Webhook, error)) FactoryAggregator {
	builder := FactoryAggregator{}

	for _, f := range m {
		builder.Register(f)
	}

	return builder
}

// Register registers the given functions in this builder.
func (a *FactoryAggregator) Register(f func(manager.Manager) (*Webhook, error)) {
	*a = append(*a, f)
}

// Webhooks calls all factories with the given managers and returns all created webhooks.
// As soon as there is an error creating a webhook, the error is returned immediately.
func (a *FactoryAggregator) Webhooks(mgr manager.Manager) ([]*Webhook, error) {
	webhooks := make([]*Webhook, 0, len(*a))

	for _, f := range *a {
		wh, err := f(mgr)
		if err != nil {
			return nil, err
		}

		webhooks = append(webhooks, wh)
	}

	return webhooks, nil
}
