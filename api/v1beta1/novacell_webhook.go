/*
Copyright 2023.

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

//
// Generated by:
//
// operator-sdk create webhook --group nova --version v1beta1 --kind NovaCell --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// NovaCellDefaults -
type NovaCellDefaults struct {
	ConductorContainerImageURL         string
	MetadataContainerImageURL          string
	NoVNCContainerImageURL             string
	NovaIronicComputeContainerImageURL string
}

var novaCellDefaults NovaCellDefaults

// log is for logging in this package.
var novacelllog = logf.Log.WithName("novacell-resource")

// SetupNovaCellDefaults - initialize NovaCell spec defaults for use with either internal or external webhooks
func SetupNovaCellDefaults(defaults NovaCellDefaults) {
	novaCellDefaults = defaults
	novacelllog.Info("NovaCell defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *NovaCell) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-nova-openstack-org-v1beta1-novacell,mutating=true,failurePolicy=fail,sideEffects=None,groups=nova.openstack.org,resources=novacells,verbs=create;update,versions=v1beta1,name=mnovacell.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &NovaCell{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NovaCell) Default() {
	novacelllog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this NovaCell spec
func (spec *NovaCellSpec) Default() {
	if spec.ConductorServiceTemplate.ContainerImage == "" {
		spec.ConductorServiceTemplate.ContainerImage = novaCellDefaults.ConductorContainerImageURL
	}

	if spec.MetadataServiceTemplate.ContainerImage == "" {
		spec.MetadataServiceTemplate.ContainerImage = novaCellDefaults.MetadataContainerImageURL
	}
	if spec.MetadataServiceTemplate.Enabled == nil {
		spec.MetadataServiceTemplate.Enabled = ptr.To(false)
	}

	if spec.NoVNCProxyServiceTemplate.ContainerImage == "" {
		spec.NoVNCProxyServiceTemplate.ContainerImage = novaCellDefaults.NoVNCContainerImageURL
	}

	if spec.CellName == Cell0Name {
		// in cell0 disable VNC by default
		if spec.NoVNCProxyServiceTemplate.Enabled == nil {
			spec.NoVNCProxyServiceTemplate.Enabled = ptr.To(false)
		}
	} else {
		// in other cells enable VNC  by default
		if spec.NoVNCProxyServiceTemplate.Enabled == nil {
			spec.NoVNCProxyServiceTemplate.Enabled = ptr.To(true)
		}
	}
	if spec.NovaComputeIronicServiceTemplate.ContainerImage == "" {
		spec.NovaComputeIronicServiceTemplate.ContainerImage = novaCellDefaults.NovaIronicComputeContainerImageURL
	}
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-nova-openstack-org-v1beta1-novacell,mutating=false,failurePolicy=fail,sideEffects=None,groups=nova.openstack.org,resources=novacells,verbs=create;update,versions=v1beta1,name=vnovacell.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &NovaCell{}

func (r *NovaCellSpec) validate(basePath *field.Path) field.ErrorList {
	var errors field.ErrorList

	if r.CellName == Cell0Name {
		errors = append(
			errors, r.MetadataServiceTemplate.ValidateCell0(
				basePath.Child("metadataServiceTemplate"))...,
		)
		errors = append(
			errors, r.NoVNCProxyServiceTemplate.ValidateCell0(
				basePath.Child("noVNCProxyServiceTemplate"))...,
		)
	}

	errors = append(
		errors, ValidateCellName(
			basePath.Child("cellName"), r.CellName)...,
	)
	return errors
}

func (r *NovaCellSpec) ValidateCreate(basePath *field.Path) field.ErrorList {
	return r.validate(basePath)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NovaCell) ValidateCreate() error {
	novacelllog.Info("validate create", "name", r.Name)

	errors := r.Spec.ValidateCreate(field.NewPath("spec"))

	if len(errors) != 0 {
		novacelllog.Info("validation failed", "name", r.Name)
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "nova.openstack.org", Kind: "NovaCell"},
			r.Name, errors)
	}
	return nil
}

func (r *NovaCellSpec) ValidateUpdate(old NovaCellSpec, basePath *field.Path) field.ErrorList {
	return r.validate(basePath)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NovaCell) ValidateUpdate(old runtime.Object) error {
	novacelllog.Info("validate update", "name", r.Name)
	oldCell, ok := old.(*NovaCell)
	if !ok || oldCell == nil {
		return apierrors.NewInternalError(fmt.Errorf("unable to convert existing object"))
	}

	errors := r.Spec.ValidateUpdate(oldCell.Spec, field.NewPath("spec"))

	if len(errors) != 0 {
		novacelllog.Info("validation failed", "name", r.Name)
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "nova.openstack.org", Kind: "NovaCell"},
			r.Name, errors)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NovaCell) ValidateDelete() error {
	novacelllog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

// ValidateCellName validates the cell name. It is expected to be called
// from various webhooks.
func ValidateCellName(path *field.Path, cellName string) field.ErrorList {
	var errors field.ErrorList
	if len(cellName) > 35 {
		errors = append(
			errors,
			field.Invalid(
				path, cellName, "should be shorter than 36 characters"),
		)
	}
	return errors
}
