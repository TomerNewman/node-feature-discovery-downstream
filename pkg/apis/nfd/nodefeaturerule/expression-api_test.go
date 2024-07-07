/*
Copyright 2021 The Kubernetes Authors.

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

package nodefeaturerule_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	nfdv1alpha1 "github.com/openshift/node-feature-discovery/api/nfd/v1alpha1"
	api "github.com/openshift/node-feature-discovery/pkg/apis/nfd/nodefeaturerule"
)

type BoolAssertionFunc func(assert.TestingT, bool, ...interface{}) bool

type ValueAssertionFunc func(assert.TestingT, interface{}, ...interface{}) bool

func TestMatchKeys(t *testing.T) {
	type I = map[string]nfdv1alpha1.Nil
	type O = []api.MatchedElement
	type TC struct {
		name   string
		mes    string
		input  I
		output O
		result BoolAssertionFunc
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{
			name:   "empty expression and nil input",
			output: O{},
			result: assert.True,
			err:    assert.Nil,
		},
		{
			name:   "empty expression and empty input",
			input:  I{},
			output: O{},
			result: assert.True,
			err:    assert.Nil,
		},
		{
			name:   "empty expression with non-empty input",
			input:  I{"foo": {}},
			output: O{},
			result: assert.True,
			err:    assert.Nil,
		},
		{
			name: "expressions match",
			mes: `
foo: { op: DoesNotExist }
bar: { op: Exists }
`,
			input:  I{"bar": {}, "baz": {}, "buzz": {}},
			output: O{{"Name": "bar"}, {"Name": "foo"}},
			result: assert.True,
			err:    assert.Nil,
		},
		{
			name: "expression does not match",
			mes: `
foo: { op: DoesNotExist }
bar: { op: Exists }
`,
			input:  I{"foo": {}, "bar": {}, "baz": {}},
			output: nil,
			result: assert.False,
			err:    assert.Nil,
		},
		{
			name: "op that never matches",
			mes: `
foo: { op: In, value: ["bar"] }
bar: { op: Exists }
`,
			input:  I{"bar": {}, "baz": {}},
			output: nil,
			result: assert.False,
			err:    assert.Nil,
		},
		{
			name: "error in expression",
			mes: `
foo: { op: Exists, value: ["bar"] }
bar: { op: Exists }
`,
			input:  I{"bar": {}},
			output: nil,
			result: assert.False,
			err:    assert.NotNil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mes := &nfdv1alpha1.MatchExpressionSet{}
			if err := yaml.Unmarshal([]byte(tc.mes), mes); err != nil {
				t.Fatal("failed to parse data of test case")
			}

			res, out, err := api.MatchGetKeys(mes, tc.input)
			tc.result(t, res)
			assert.Equal(t, tc.output, out)
			tc.err(t, err)

			res, err = api.MatchKeys(mes, tc.input)
			tc.result(t, res)
			tc.err(t, err)
		})
	}
}

func TestMatchValues(t *testing.T) {
	type I = map[string]string
	type O = []api.MatchedElement
	type TC struct {
		name   string
		mes    string
		input  I
		output O
		result BoolAssertionFunc
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{name: "1", output: O{}, result: assert.True, err: assert.Nil},

		{name: "2", input: I{}, output: O{}, result: assert.True, err: assert.Nil},

		{name: "3", input: I{"foo": "bar"}, output: O{}, result: assert.True, err: assert.Nil},

		{name: "4",
			mes: `
foo: { op: Exists }
bar: { op: In, value: ["val", "wal"] }
baz: { op: Gt, value: ["10"] }
`,
			input:  I{"bar": "val"},
			result: assert.False, err: assert.Nil},

		{name: "5",
			mes: `
foo: { op: Exists }
bar: { op: In, value: ["val", "wal"] }
baz: { op: Gt, value: ["10"] }
`,
			input:  I{"foo": "1", "bar": "val", "baz": "123", "buzz": "light"},
			output: O{{"Name": "bar", "Value": "val"}, {"Name": "baz", "Value": "123"}, {"Name": "foo", "Value": "1"}},
			result: assert.True, err: assert.Nil},

		{name: "5",
			mes: `
foo: { op: Exists }
bar: { op: In, value: ["val"] }
baz: { op: Gt, value: ["10"] }
`,
			input:  I{"foo": "1", "bar": "val", "baz": "123.0"},
			result: assert.False, err: assert.NotNil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mes := &nfdv1alpha1.MatchExpressionSet{}
			if err := yaml.Unmarshal([]byte(tc.mes), mes); err != nil {
				t.Fatal("failed to parse data of test case")
			}

			res, out, err := api.MatchGetValues(mes, tc.input)
			tc.result(t, res)
			assert.Equal(t, tc.output, out)
			tc.err(t, err)

			res, err = api.MatchValues(mes, tc.input)
			tc.result(t, res)
			tc.err(t, err)
		})
	}
}

func TestMatchInstances(t *testing.T) {
	type I = nfdv1alpha1.InstanceFeature
	type O = []api.MatchedElement
	type A = map[string]string
	type TC struct {
		name   string
		mes    string
		input  []I
		output O
		result BoolAssertionFunc
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{name: "1", output: O{}, result: assert.False, err: assert.Nil}, // nil instances -> false

		{name: "2", input: []I{}, output: O{}, result: assert.False, err: assert.Nil}, // zero instances -> false

		{name: "3", input: []I{I{Attributes: A{}}}, output: O{A{}}, result: assert.True, err: assert.Nil}, // one "empty" instance

		{name: "4",
			mes: `
foo: { op: Exists }
bar: { op: Lt, value: ["10"] }
`,
			input:  []I{I{Attributes: A{"foo": "1"}}, I{Attributes: A{"bar": "1"}}},
			output: O{},
			result: assert.False, err: assert.Nil},

		{name: "5",
			mes: `
foo: { op: Exists }
bar: { op: Lt, value: ["10"] }
`,
			input:  []I{I{Attributes: A{"foo": "1"}}, I{Attributes: A{"foo": "2", "bar": "1"}}},
			output: O{A{"foo": "2", "bar": "1"}},
			result: assert.True, err: assert.Nil},

		{name: "6",
			mes: `
bar: { op: Lt, value: ["10"] }
`,
			input:  []I{I{Attributes: A{"foo": "1"}}, I{Attributes: A{"bar": "0x1"}}},
			result: assert.False, err: assert.NotNil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mes := &nfdv1alpha1.MatchExpressionSet{}
			if err := yaml.Unmarshal([]byte(tc.mes), mes); err != nil {
				t.Fatal("failed to parse data of test case")
			}

			out, err := api.MatchGetInstances(mes, tc.input)
			assert.Equal(t, tc.output, out)
			tc.err(t, err)

			res, err := api.MatchInstances(mes, tc.input)
			tc.result(t, res)
			tc.err(t, err)
		})
	}
}

func TestMatchKeyNames(t *testing.T) {
	type O = []api.MatchedElement
	type I = map[string]nfdv1alpha1.Nil

	type TC struct {
		name   string
		me     *nfdv1alpha1.MatchExpression
		input  I
		result bool
		output O
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{
			name:   "empty input",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchAny},
			input:  I{},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchAny",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchAny},
			input:  I{"key1": {}, "key2": {}},
			result: true,
			output: O{{"Name": "key1"}, {"Name": "key2"}},
			err:    assert.Nil,
		},
		{
			name:   "MatchExists",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchExists},
			input:  I{"key1": {}, "key2": {}},
			result: true,
			output: O{{"Name": "key1"}, {"Name": "key2"}},
			err:    assert.Nil,
		},
		{
			name:   "MatchDoesNotExist",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchDoesNotExist},
			input:  I{"key1": {}, "key2": {}},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchIn matches",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchIn, Value: nfdv1alpha1.MatchValue{"key1"}},
			input:  I{"key1": {}, "key2": {}},
			result: true,
			output: O{{"Name": "key1"}},
			err:    assert.Nil,
		},
		{
			name:   "MatchIn no match",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchIn, Value: nfdv1alpha1.MatchValue{"key3"}},
			input:  I{"key1": {}, "key2": {}},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchNotIn",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchNotIn, Value: nfdv1alpha1.MatchValue{"key1"}},
			input:  I{"key1": {}, "key2": {}},
			result: true,
			output: O{{"Name": "key2"}},
			err:    assert.Nil,
		},
		{
			name:   "error",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchExists, Value: nfdv1alpha1.MatchValue{"key1"}},
			input:  I{"key1": {}, "key2": {}},
			result: false,
			output: nil,
			err:    assert.NotNil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, ret, err := api.MatchKeyNames(tc.me, tc.input)
			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.output, ret)
			tc.err(t, err)
		})
	}
}

func TestMatchValueNames(t *testing.T) {
	type O = []api.MatchedElement
	type I = map[string]string

	type TC struct {
		name   string
		me     *nfdv1alpha1.MatchExpression
		input  I
		result bool
		output O
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{
			name:   "empty input",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchAny},
			input:  I{},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchExists",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchExists},
			input:  I{"key1": "val1", "key2": "val2"},
			result: true,
			output: O{{"Name": "key1", "Value": "val1"}, {"Name": "key2", "Value": "val2"}},
			err:    assert.Nil,
		},
		{
			name:   "MatchDoesNotExist",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchDoesNotExist},
			input:  I{"key1": "val1", "key2": "val2"},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchIn matches",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchIn, Value: nfdv1alpha1.MatchValue{"key1"}},
			input:  I{"key1": "val1", "key2": "val2"},
			result: true,
			output: O{{"Name": "key1", "Value": "val1"}},
			err:    assert.Nil,
		},
		{
			name:   "MatchIn no match",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchIn, Value: nfdv1alpha1.MatchValue{"key3"}},
			input:  I{"key1": "val1", "key2": "val2"},
			result: false,
			output: O{},
			err:    assert.Nil,
		},
		{
			name:   "MatchNotIn",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchNotIn, Value: nfdv1alpha1.MatchValue{"key1"}},
			input:  I{"key1": "val1", "key2": "val2"},
			result: true,
			output: O{{"Name": "key2", "Value": "val2"}},
			err:    assert.Nil,
		},
		{
			name:   "error",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchNotIn},
			input:  I{"key1": "val1", "key2": "val2"},
			result: false,
			output: nil,
			err:    assert.NotNil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, ret, err := api.MatchValueNames(tc.me, tc.input)
			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.output, ret)
			tc.err(t, err)
		})
	}
}

func TestMatchInstanceAttributeNames(t *testing.T) {
	type O = []api.MatchedElement
	type I = []nfdv1alpha1.InstanceFeature
	type A = map[string]string

	type TC struct {
		name   string
		me     *nfdv1alpha1.MatchExpression
		input  I
		output O
		err    ValueAssertionFunc
	}

	tcs := []TC{
		{
			name:   "empty input",
			me:     &nfdv1alpha1.MatchExpression{Op: nfdv1alpha1.MatchAny},
			input:  I{},
			output: O{},
			err:    assert.Nil,
		},
		{
			name: "no match",
			me: &nfdv1alpha1.MatchExpression{
				Op:    nfdv1alpha1.MatchIn,
				Value: nfdv1alpha1.MatchValue{"foo"},
			},
			input: I{
				{Attributes: A{"bar": "1"}},
				{Attributes: A{"baz": "2"}},
			},
			output: O{},
			err:    assert.Nil,
		},
		{
			name: "match",
			me: &nfdv1alpha1.MatchExpression{
				Op:    nfdv1alpha1.MatchIn,
				Value: nfdv1alpha1.MatchValue{"foo"},
			},
			input: I{
				{Attributes: A{"foo": "1"}},
				{Attributes: A{"bar": "2"}},
				{Attributes: A{"foo": "3", "baz": "4"}},
			},
			output: O{
				{"foo": "1"},
				{"foo": "3", "baz": "4"},
			},
			err: assert.Nil,
		},
		{
			name: "error",
			me: &nfdv1alpha1.MatchExpression{
				Op: nfdv1alpha1.MatchIn,
			},
			input: I{
				{Attributes: A{"foo": "1"}},
			},
			output: nil,
			err:    assert.NotNil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			matched, err := api.MatchInstanceAttributeNames(tc.me, tc.input)
			assert.Equal(t, tc.output, matched)
			tc.err(t, err)
		})
	}
}
