// Copyright 2022 Upbound Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/upbound/upjet/pkg/terraform"

	"github.com/upbound/provider-azuread/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig         = "no providerConfigRef provided"
	errGetProviderConfig        = "cannot get referenced ProviderConfig"
	errTrackUsage               = "cannot track ProviderConfig usage"
	errExtractCredentials       = "cannot extract credentials"
	errTenantIDNotSet           = "tenant ID must be set in ProviderConfig when credential source is InjectedIdentity, OIDCTokenFile or Upbound"
	errUnmarshalCredentials     = "cannot unmarshal azuread credentials as JSON"
	keySubscriptionID           = "subscription_id"
	keyUseMSI                   = "use_msi"
	keyMSIEndpoint              = "msi_endpoint"
	keyEnvironment              = "environment"
	keyTerraformClientID        = "client_id"
	keyTerraformClientSecret    = "client_secret"
	keyTerraformTenantID        = "tenant_id"
	keyTerraformFeatures        = "features"
	keySkipProviderRegistration = "skip_provider_registration"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string, scheduler terraform.ProviderScheduler) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
			Scheduler: scheduler,
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		ps.Configuration = map[string]interface{}{
			// keyTerraformFeatures: struct{}{},
			// Terraform AzureRM provider tries to register all resource providers
			// in Azure just in case if the provider of the resource you're
			// trying to create is not registered and the returned error is
			// ambiguous. However, this requires service principal to have provider
			// registration permissions which are irrelevant in most contexts.
			// For details, see https://github.com/upbound/provider-azure/issues/104
			// keySkipProviderRegistration: true,
		}

		var err = msiAuth(pc, &ps)
		return ps, err
	}
}

func msiAuth(pc *v1beta1.ProviderConfig, ps *terraform.Setup) error {
	if pc.Spec.TenantID == nil || len(*pc.Spec.TenantID) == 0 {
		return errors.New(errTenantIDNotSet)
	}
	ps.Configuration[keyTerraformTenantID] = *pc.Spec.TenantID
	ps.Configuration[keyUseMSI] = "true"
	if pc.Spec.MSIEndpoint != nil {
		ps.Configuration[keyMSIEndpoint] = *pc.Spec.MSIEndpoint
	}
	if pc.Spec.ClientID != nil {
		ps.Configuration[keyTerraformClientID] = *pc.Spec.ClientID
	}
	if pc.Spec.Environment != nil {
		ps.Configuration[keyEnvironment] = *pc.Spec.Environment
	}
	return nil
}
