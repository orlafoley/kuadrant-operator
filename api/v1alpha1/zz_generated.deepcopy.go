//go:build !ignore_autogenerated

/*
Copyright 2021.

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
	"github.com/kuadrant/authorino/api/v1beta3"
	apiv1 "github.com/kuadrant/kuadrant-operator/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Auth) DeepCopyInto(out *Auth) {
	*out = *in
	if in.TokenSource != nil {
		in, out := &in.TokenSource, &out.TokenSource
		*out = new(v1beta3.Credentials)
		(*in).DeepCopyInto(*out)
	}
	if in.Claims != nil {
		in, out := &in.Claims, &out.Claims
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Auth.
func (in *Auth) DeepCopy() *Auth {
	if in == nil {
		return nil
	}
	out := new(Auth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Limits) DeepCopyInto(out *Limits) {
	*out = *in
	if in.Daily != nil {
		in, out := &in.Daily, &out.Daily
		*out = new(int)
		**out = **in
	}
	if in.Weekly != nil {
		in, out := &in.Weekly, &out.Weekly
		*out = new(int)
		**out = **in
	}
	if in.Monthly != nil {
		in, out := &in.Monthly, &out.Monthly
		*out = new(int)
		**out = **in
	}
	if in.Yearly != nil {
		in, out := &in.Yearly, &out.Yearly
		*out = new(int)
		**out = **in
	}
	if in.Custom != nil {
		in, out := &in.Custom, &out.Custom
		*out = make([]apiv1.Rate, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Limits.
func (in *Limits) DeepCopy() *Limits {
	if in == nil {
		return nil
	}
	out := new(Limits)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MergeableTokenRateLimitPolicySpec) DeepCopyInto(out *MergeableTokenRateLimitPolicySpec) {
	*out = *in
	in.TokenRateLimitPolicySpecProper.DeepCopyInto(&out.TokenRateLimitPolicySpecProper)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MergeableTokenRateLimitPolicySpec.
func (in *MergeableTokenRateLimitPolicySpec) DeepCopy() *MergeableTokenRateLimitPolicySpec {
	if in == nil {
		return nil
	}
	out := new(MergeableTokenRateLimitPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCPolicy) DeepCopyInto(out *OIDCPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCPolicy.
func (in *OIDCPolicy) DeepCopy() *OIDCPolicy {
	if in == nil {
		return nil
	}
	out := new(OIDCPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OIDCPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCPolicyList) DeepCopyInto(out *OIDCPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OIDCPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCPolicyList.
func (in *OIDCPolicyList) DeepCopy() *OIDCPolicyList {
	if in == nil {
		return nil
	}
	out := new(OIDCPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OIDCPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCPolicySpec) DeepCopyInto(out *OIDCPolicySpec) {
	*out = *in
	in.TargetRef.DeepCopyInto(&out.TargetRef)
	in.OIDCPolicySpecProper.DeepCopyInto(&out.OIDCPolicySpecProper)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCPolicySpec.
func (in *OIDCPolicySpec) DeepCopy() *OIDCPolicySpec {
	if in == nil {
		return nil
	}
	out := new(OIDCPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCPolicySpecProper) DeepCopyInto(out *OIDCPolicySpecProper) {
	*out = *in
	if in.Provider != nil {
		in, out := &in.Provider, &out.Provider
		*out = new(Provider)
		**out = **in
	}
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(Auth)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCPolicySpecProper.
func (in *OIDCPolicySpecProper) DeepCopy() *OIDCPolicySpecProper {
	if in == nil {
		return nil
	}
	out := new(OIDCPolicySpecProper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OIDCPolicyStatus) DeepCopyInto(out *OIDCPolicyStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OIDCPolicyStatus.
func (in *OIDCPolicyStatus) DeepCopy() *OIDCPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(OIDCPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Plan) DeepCopyInto(out *Plan) {
	*out = *in
	in.Limits.DeepCopyInto(&out.Limits)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Plan.
func (in *Plan) DeepCopy() *Plan {
	if in == nil {
		return nil
	}
	out := new(Plan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlanPolicy) DeepCopyInto(out *PlanPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlanPolicy.
func (in *PlanPolicy) DeepCopy() *PlanPolicy {
	if in == nil {
		return nil
	}
	out := new(PlanPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlanPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlanPolicyList) DeepCopyInto(out *PlanPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PlanPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlanPolicyList.
func (in *PlanPolicyList) DeepCopy() *PlanPolicyList {
	if in == nil {
		return nil
	}
	out := new(PlanPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlanPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlanPolicySpec) DeepCopyInto(out *PlanPolicySpec) {
	*out = *in
	in.TargetRef.DeepCopyInto(&out.TargetRef)
	if in.Plans != nil {
		in, out := &in.Plans, &out.Plans
		*out = make([]Plan, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlanPolicySpec.
func (in *PlanPolicySpec) DeepCopy() *PlanPolicySpec {
	if in == nil {
		return nil
	}
	out := new(PlanPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlanPolicyStatus) DeepCopyInto(out *PlanPolicyStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlanPolicyStatus.
func (in *PlanPolicyStatus) DeepCopy() *PlanPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(PlanPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provider) DeepCopyInto(out *Provider) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provider.
func (in *Provider) DeepCopy() *Provider {
	if in == nil {
		return nil
	}
	out := new(Provider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenLimit) DeepCopyInto(out *TokenLimit) {
	*out = *in
	if in.When != nil {
		in, out := &in.When, &out.When
		*out = make(apiv1.WhenPredicates, len(*in))
		copy(*out, *in)
	}
	if in.Rates != nil {
		in, out := &in.Rates, &out.Rates
		*out = make([]apiv1.Rate, len(*in))
		copy(*out, *in)
	}
	if in.Counters != nil {
		in, out := &in.Counters, &out.Counters
		*out = make([]apiv1.Counter, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenLimit.
func (in *TokenLimit) DeepCopy() *TokenLimit {
	if in == nil {
		return nil
	}
	out := new(TokenLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRateLimitPolicy) DeepCopyInto(out *TokenRateLimitPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRateLimitPolicy.
func (in *TokenRateLimitPolicy) DeepCopy() *TokenRateLimitPolicy {
	if in == nil {
		return nil
	}
	out := new(TokenRateLimitPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TokenRateLimitPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRateLimitPolicyList) DeepCopyInto(out *TokenRateLimitPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TokenRateLimitPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRateLimitPolicyList.
func (in *TokenRateLimitPolicyList) DeepCopy() *TokenRateLimitPolicyList {
	if in == nil {
		return nil
	}
	out := new(TokenRateLimitPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TokenRateLimitPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRateLimitPolicySpec) DeepCopyInto(out *TokenRateLimitPolicySpec) {
	*out = *in
	in.TargetRef.DeepCopyInto(&out.TargetRef)
	if in.Defaults != nil {
		in, out := &in.Defaults, &out.Defaults
		*out = new(MergeableTokenRateLimitPolicySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Overrides != nil {
		in, out := &in.Overrides, &out.Overrides
		*out = new(MergeableTokenRateLimitPolicySpec)
		(*in).DeepCopyInto(*out)
	}
	in.TokenRateLimitPolicySpecProper.DeepCopyInto(&out.TokenRateLimitPolicySpecProper)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRateLimitPolicySpec.
func (in *TokenRateLimitPolicySpec) DeepCopy() *TokenRateLimitPolicySpec {
	if in == nil {
		return nil
	}
	out := new(TokenRateLimitPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRateLimitPolicySpecProper) DeepCopyInto(out *TokenRateLimitPolicySpecProper) {
	*out = *in
	in.MergeableWhenPredicates.DeepCopyInto(&out.MergeableWhenPredicates)
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = make(map[string]TokenLimit, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRateLimitPolicySpecProper.
func (in *TokenRateLimitPolicySpecProper) DeepCopy() *TokenRateLimitPolicySpecProper {
	if in == nil {
		return nil
	}
	out := new(TokenRateLimitPolicySpecProper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRateLimitPolicyStatus) DeepCopyInto(out *TokenRateLimitPolicyStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRateLimitPolicyStatus.
func (in *TokenRateLimitPolicyStatus) DeepCopy() *TokenRateLimitPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(TokenRateLimitPolicyStatus)
	in.DeepCopyInto(out)
	return out
}
