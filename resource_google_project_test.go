package google

import (
	"reflect"
	"sort"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

type Binding []cloudresourcemanager.Binding

func (b Binding) Len() int {
	return len(b)
}

func (b Binding) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Binding) Less(i, j int) bool {
	return b[i].Role < b[j].Role
}

func TestIamMapRolesToMembers(t *testing.T) {
	table := []struct {
		input  []cloudresourcemanager.Binding
		expect map[string]map[string]bool
	}{
		{
			input: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-1",
						"member-2",
					},
				},
			},
			expect: map[string]map[string]bool{
				"role-1": map[string]bool{
					"member-1": true,
					"member-2": true,
				},
			},
		},
		{
			input: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-1",
						"member-2",
						"member-1",
						"member-2",
					},
				},
			},
			expect: map[string]map[string]bool{
				"role-1": map[string]bool{
					"member-1": true,
					"member-2": true,
				},
			},
		},
		{
			input: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
				},
			},
			expect: map[string]map[string]bool{
				"role-1": map[string]bool{},
			},
		},
	}

	for _, test := range table {
		got := mapRolesToMembers(test.input)
		if !reflect.DeepEqual(got, test.expect) {
			t.Errorf("got %+v, expected %+v", got, test.expect)
		}
	}
}

func TestIamMergeBindings(t *testing.T) {
	table := []struct {
		input  []cloudresourcemanager.Binding
		expect []cloudresourcemanager.Binding
	}{
		{
			input: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-1",
						"member-2",
					},
				},
				{
					Role: "role-1",
					Members: []string{
						"member-3",
					},
				},
			},
			expect: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-1",
						"member-2",
						"member-3",
					},
				},
			},
		},
		{
			input: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-3",
						"member-4",
					},
				},
				{
					Role: "role-1",
					Members: []string{
						"member-2",
						"member-1",
					},
				},
				{
					Role: "role-2",
					Members: []string{
						"member-1",
					},
				},
				{
					Role: "role-1",
					Members: []string{
						"member-5",
					},
				},
				{
					Role: "role-3",
					Members: []string{
						"member-1",
					},
				},
				{
					Role: "role-2",
					Members: []string{
						"member-2",
					},
				},
			},
			expect: []cloudresourcemanager.Binding{
				{
					Role: "role-1",
					Members: []string{
						"member-1",
						"member-2",
						"member-3",
						"member-4",
						"member-5",
					},
				},
				{
					Role: "role-2",
					Members: []string{
						"member-1",
						"member-2",
					},
				},
				{
					Role: "role-3",
					Members: []string{
						"member-1",
					},
				},
			},
		},
	}

	for _, test := range table {
		got := mergeBindings(test.input)
		sort.Sort(Binding(got))
		for i, _ := range got {
			sort.Strings(got[i].Members)
		}

		if !reflect.DeepEqual(got, test.expect) {
			t.Errorf("\ngot %+v\nexpected %+v", got, test.expect)
		}
	}
}
