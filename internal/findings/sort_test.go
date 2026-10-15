package findings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	ref := func(ns, kind, name string) ResourceRef {
		return ResourceRef{Namespace: ns, Kind: kind, Name: name}
	}
	in := []Finding{
		{RuleID: "B", Resource: ref("web", "Deployment", "b")},
		{RuleID: "A", Resource: ref("web", "Deployment", "b")},
		{RuleID: "A", Resource: ref("web", "Deployment", "a")},
		{RuleID: "A", Resource: ref("web", "DaemonSet", "z")},
		{RuleID: "A", Resource: ref("billing", "Pod", "p")},
		{RuleID: "A", Resource: ref("", "Node", "worker-01")},
	}
	want := []Finding{
		{RuleID: "A", Resource: ref("", "Node", "worker-01")},
		{RuleID: "A", Resource: ref("billing", "Pod", "p")},
		{RuleID: "A", Resource: ref("web", "DaemonSet", "z")},
		{RuleID: "A", Resource: ref("web", "Deployment", "a")},
		{RuleID: "A", Resource: ref("web", "Deployment", "b")},
		{RuleID: "B", Resource: ref("web", "Deployment", "b")},
	}

	Sort(in)
	if !reflect.DeepEqual(in, want) {
		t.Errorf("Sort() order mismatch\n got: %+v\nwant: %+v", in, want)
	}

	// Sorting must be independent of the input order.
	rev := make([]Finding, len(want))
	for i, f := range want {
		rev[len(want)-1-i] = f
	}
	Sort(rev)
	if !reflect.DeepEqual(rev, want) {
		t.Errorf("Sort() not deterministic for reversed input\n got: %+v", rev)
	}
}

func TestSortTieBreaks(t *testing.T) {
	r := ResourceRef{Namespace: "web", Kind: "Pod", Name: "p"}
	in := []Finding{
		{RuleID: "B", Severity: Risk, Resource: r},
		{RuleID: "A", Severity: Unknown, Resource: r},
		{RuleID: "A", Severity: Blocked, Resource: r},
	}
	Sort(in)

	got := []string{in[0].RuleID + "/" + string(in[0].Severity), in[1].RuleID + "/" + string(in[1].Severity), in[2].RuleID + "/" + string(in[2].Severity)}
	want := []string{"A/BLOCKED", "A/UNKNOWN", "B/RISK"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSortEmpty(t *testing.T) {
	Sort(nil)
	Sort([]Finding{})
}

func ExampleSort() {
	ref := func(ns, kind, name string) ResourceRef {
		return ResourceRef{Namespace: ns, Kind: kind, Name: name}
	}
	fs := []Finding{
		{RuleID: "pdb-missing", Resource: ref("web", "Deployment", "api")},
		{RuleID: "pdb-missing", Resource: ref("billing", "StatefulSet", "db")},
		{RuleID: "pdb-missing", Resource: ref("web", "DaemonSet", "agent")},
		{RuleID: "pdb-missing", Resource: ref("billing", "Deployment", "worker")},
	}
	Sort(fs)
	for _, f := range fs {
		r := f.Resource
		fmt.Printf("%s %s/%s\n", r.Namespace, r.Kind, r.Name)
	}
	// Output:
	// billing Deployment/worker
	// billing StatefulSet/db
	// web DaemonSet/agent
	// web Deployment/api
}
