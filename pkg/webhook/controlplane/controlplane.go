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

package controlplane

import (
	extensionswebhook "github.com/gardener/gardener-extensions/pkg/webhook"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// WebhookName is the webhook name.
	WebhookName = "controlplane"
	// ExposureWebhookName is the exposure webhook name.
	ExposureWebhookName = "controlplaneexposure"
	// BackupWebhookName is the backup webhook name.
	BackupWebhookName = "controlplanebackup"
)

var logger = log.Log.WithName("controlplane-webhook")

// AddArgs are arguments for adding a controlplane webhook to a manager.
type AddArgs struct {
	// Kind is the kind of this webhook
	Kind string
	// Provider is the provider of this webhook.
	Provider string
	// Types is a list of resource types.
	Types []runtime.Object
	// Mutator is a mutator to be used by the admission handler.
	Mutator Mutator
}

// Add creates a new controlplane webhook and adds it to the given Manager.
func Add(mgr manager.Manager, args AddArgs) (*extensionswebhook.Webhook, error) {
	logger := logger.WithValues("kind", args.Kind, "provider", args.Provider)

	// Create handler
	handler, err := newHandler(mgr, args.Types, args.Mutator, logger)
	if err != nil {
		return nil, err
	}

	// Create webhook
	logger.Info("Creating webhook", "name", getName(args.Kind))

	return &extensionswebhook.Webhook{
		Name:     getName(args.Kind),
		Kind:     args.Kind,
		Provider: args.Provider,
		Types:    args.Types,
		Path:     getName(args.Kind),
		Webhook:  &admission.Webhook{Handler: handler},
	}, nil
}

func getName(kind string) string {
	switch kind {
	case extensionswebhook.SeedKind:
		return ExposureWebhookName
	case extensionswebhook.BackupKind:
		return BackupWebhookName
	default:
		return WebhookName
	}
}
