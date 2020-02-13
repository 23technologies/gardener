// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package validation

import (
	"github.com/gardener/gardener/pkg/apis/core"

	rbacv1 "k8s.io/api/rbac/v1"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateProject validates a Project object.
func ValidateProject(project *core.Project) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(&project.ObjectMeta, false, ValidateName, field.NewPath("metadata"))...)
	maxProjectNameLength := 10
	if len(project.Name) > maxProjectNameLength {
		allErrs = append(allErrs, field.TooLong(field.NewPath("metadata", "name"), project.Name, maxProjectNameLength))
	}
	allErrs = append(allErrs, validateNameConsecutiveHyphens(project.Name, field.NewPath("metadata", "name"))...)
	allErrs = append(allErrs, ValidateProjectSpec(&project.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateProjectUpdate validates a Project object before an update.
func ValidateProjectUpdate(newProject, oldProject *core.Project) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaUpdate(&newProject.ObjectMeta, &oldProject.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateProject(newProject)...)

	if oldProject.Spec.CreatedBy != nil {
		allErrs = append(allErrs, apivalidation.ValidateImmutableField(newProject.Spec.CreatedBy, oldProject.Spec.CreatedBy, field.NewPath("spec", "createdBy"))...)
	}
	if oldProject.Spec.Namespace != nil {
		allErrs = append(allErrs, apivalidation.ValidateImmutableField(newProject.Spec.Namespace, oldProject.Spec.Namespace, field.NewPath("spec", "namespace"))...)
	}
	if oldProject.Spec.Owner != nil && newProject.Spec.Owner == nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "owner"), newProject.Spec.Owner, "owner cannot be reset"))
	}

	return allErrs
}

// ValidateProjectSpec validates the specification of a Project object.
func ValidateProjectSpec(projectSpec *core.ProjectSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, member := range projectSpec.Members {
		allErrs = append(allErrs, ValidateSubject(member.Subject, fldPath.Child("members").Index(i))...)
	}
	if createdBy := projectSpec.CreatedBy; createdBy != nil {
		allErrs = append(allErrs, ValidateSubject(*createdBy, fldPath.Child("createdBy"))...)
	}
	if owner := projectSpec.Owner; owner != nil {
		allErrs = append(allErrs, ValidateSubject(*owner, fldPath.Child("owner"))...)
	}
	if description := projectSpec.Description; description != nil && len(*description) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("description"), "must provide a description when key is present"))
	}
	if purpose := projectSpec.Purpose; purpose != nil && len(*purpose) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("purpose"), "must provide a purpose when key is present"))
	}

	return allErrs
}

// ValidateSubject validates the subject representing the owner.
func ValidateSubject(subject rbacv1.Subject, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(subject.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	switch subject.Kind {
	case rbacv1.ServiceAccountKind:
		if len(subject.Name) > 0 {
			for _, msg := range apivalidation.ValidateServiceAccountName(subject.Name, false) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), subject.Name, msg))
			}
		}
		if len(subject.APIGroup) > 0 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("apiGroup"), subject.APIGroup, []string{""}))
		}
		if len(subject.Namespace) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), ""))
		}

	case rbacv1.UserKind, rbacv1.GroupKind:
		if subject.APIGroup != rbacv1.GroupName {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("apiGroup"), subject.APIGroup, []string{rbacv1.GroupName}))
		}

	default:
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("kind"), subject.Kind, []string{rbacv1.ServiceAccountKind, rbacv1.UserKind, rbacv1.GroupKind}))
	}

	return allErrs
}

// ValidateProjectStatusUpdate validates the status field of a Project object.
func ValidateProjectStatusUpdate(newProject, oldProject *core.Project) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(oldProject.Status.Phase) > 0 && len(newProject.Status.Phase) == 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("status").Child("phase"), "phase cannot be updated to an empty string"))
	}

	return allErrs
}
