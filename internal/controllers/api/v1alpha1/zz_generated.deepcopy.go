//go:build !ignore_autogenerated

/*
Copyright 2023 The Stefan Prodan

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
func (in *ClusterProject) DeepCopyInto(out *ClusterProject) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProject.
func (in *ClusterProject) DeepCopy() *ClusterProject {
	if in == nil {
		return nil
	}
	out := new(ClusterProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterProject) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectAccess) DeepCopyInto(out *ClusterProjectAccess) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectAccess.
func (in *ClusterProjectAccess) DeepCopy() *ClusterProjectAccess {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectAccess)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectList) DeepCopyInto(out *ClusterProjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterProject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectList.
func (in *ClusterProjectList) DeepCopy() *ClusterProjectList {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterProjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectSource) DeepCopyInto(out *ClusterProjectSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectSource.
func (in *ClusterProjectSource) DeepCopy() *ClusterProjectSource {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectSpec) DeepCopyInto(out *ClusterProjectSpec) {
	*out = *in
	in.Access.DeepCopyInto(&out.Access)
	out.Source = in.Source
	if in.Watches != nil {
		in, out := &in.Watches, &out.Watches
		*out = make([]ClusterProjectWatch, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectSpec.
func (in *ClusterProjectSpec) DeepCopy() *ClusterProjectSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectStatus) DeepCopyInto(out *ClusterProjectStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectStatus.
func (in *ClusterProjectStatus) DeepCopy() *ClusterProjectStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterProjectWatch) DeepCopyInto(out *ClusterProjectWatch) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterProjectWatch.
func (in *ClusterProjectWatch) DeepCopy() *ClusterProjectWatch {
	if in == nil {
		return nil
	}
	out := new(ClusterProjectWatch)
	in.DeepCopyInto(out)
	return out
}
