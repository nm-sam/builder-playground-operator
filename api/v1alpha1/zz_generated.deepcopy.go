//go:build !ignore_autogenerated

/*
Copyright 2025.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuilderPlaygroundDeployment) DeepCopyInto(out *BuilderPlaygroundDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuilderPlaygroundDeployment.
func (in *BuilderPlaygroundDeployment) DeepCopy() *BuilderPlaygroundDeployment {
	if in == nil {
		return nil
	}
	out := new(BuilderPlaygroundDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BuilderPlaygroundDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuilderPlaygroundDeploymentList) DeepCopyInto(out *BuilderPlaygroundDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BuilderPlaygroundDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuilderPlaygroundDeploymentList.
func (in *BuilderPlaygroundDeploymentList) DeepCopy() *BuilderPlaygroundDeploymentList {
	if in == nil {
		return nil
	}
	out := new(BuilderPlaygroundDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BuilderPlaygroundDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuilderPlaygroundDeploymentSpec) DeepCopyInto(out *BuilderPlaygroundDeploymentSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuilderPlaygroundDeploymentSpec.
func (in *BuilderPlaygroundDeploymentSpec) DeepCopy() *BuilderPlaygroundDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(BuilderPlaygroundDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuilderPlaygroundDeploymentStatus) DeepCopyInto(out *BuilderPlaygroundDeploymentStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuilderPlaygroundDeploymentStatus.
func (in *BuilderPlaygroundDeploymentStatus) DeepCopy() *BuilderPlaygroundDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(BuilderPlaygroundDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}
