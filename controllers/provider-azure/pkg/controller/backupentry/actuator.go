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

package backupentry

import (
	"context"
	"fmt"

	azureclient "github.com/gardener/gardener-extensions/controllers/provider-azure/pkg/azure/client"
	"github.com/gardener/gardener-extensions/pkg/controller/backupentry/genericactuator"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type actuator struct {
	client client.Client
	logger logr.Logger
}

func newActuator() genericactuator.BackupEntryDelegate {
	return &actuator{
		logger: logger,
	}
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *actuator) GetETCDSecretData(ctx context.Context, be *extensionsv1alpha1.BackupEntry, backupSecretData map[string][]byte) (map[string][]byte, error) {
	return backupSecretData, nil
}

func (a *actuator) Delete(ctx context.Context, be *extensionsv1alpha1.BackupEntry) error {
	azureClient, err := azureclient.NewStorageClientFromSecretRef(ctx, a.client, be.Spec.SecretRef)
	if err != nil {
		return err
	}

	return azureClient.DeleteObjectsWithPrefix(ctx, be.Spec.BucketName, fmt.Sprintf("%s/", be.Name))
}
