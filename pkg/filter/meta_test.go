package filter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
)

func TestMetaFiltering(t *testing.T) {
	t.Run("SuccessName", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			Name:    "success_object",
			UID:     runtime.UID("myuid"),
			Created: runtime.Time{},
			Labels: map[string]string{
				"first":  "f_value",
				"second": "s_value",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "success_object",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailName", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			Name: "fail_object",
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "success_object",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessUID", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			UID: runtime.UID("myuid"),
		}

		f := metaFilter{
			identifier:    "{{.UID}}",
			expectedValue: "myuid",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailUID", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			UID: "failuid",
		}

		f := metaFilter{
			identifier:    "{{.UID}}",
			expectedValue: "myuid",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessCreated", func(t *testing.T) {
		nowtime := runtime.Timestamp()
		oMeta := &runtime.ObjectMeta{
			Created: nowtime,
		}

		f := metaFilter{
			identifier:    "{{.Created}}",
			expectedValue: nowtime.String(),
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailCreated", func(t *testing.T) {
		nowtime := runtime.Timestamp()
		oMeta := &runtime.ObjectMeta{
			Created: nowtime,
		}

		othertime := nowtime.Add(time.Duration(5))
		f := metaFilter{
			identifier:    "{{.Created}}",
			expectedValue: othertime.String(),
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessLabels", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			Labels: map[string]string{
				"foo": "bar",
			},
		}

		f := metaFilter{
			identifier:    "{{.Labels.foo}}",
			expectedValue: "bar",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailLabels", func(t *testing.T) {
		oMeta := &runtime.ObjectMeta{
			Labels: map[string]string{
				"foo": "bar2",
			},
		}

		f := metaFilter{
			identifier:    "{{.Labels.foo}}",
			expectedValue: "bar",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
}

func TestExtractKeyValueFiltering(t *testing.T) {
	tests := []struct {
		name string
		str  string
		key  string
		val  string
		err  error
	}{
		{
			name: "Success1",
			str:  "{{.Name}}=ta-rg_et",
			key:  "{{.Name}}",
			val:  "ta-rg_et",
			err:  nil,
		},
		{
			name: "Success2",
			str:  "{{.Name}}=8",
			key:  "{{.Name}}",
			val:  "8",
			err:  nil,
		},
		{
			name: "FailEqualBadPlace",
			str:  "{{.Name=}}target",
			key:  "",
			val:  "",
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailEqualBadPlace2",
			str:  "={{.Name}}target",
			key:  "",
			val:  "",
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailEqualBadPlace3",
			str:  "{{.Name}}tar=get",
			key:  "",
			val:  "",
			err:  fmt.Errorf("expected error"),
		},
	}
	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			key, val, err := extractKeyValueFiltering(utest.str)
			if utest.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
			assert.Equal(t, utest.key, key)
			assert.Equal(t, utest.val, val)
		})
	}
}

func TestExtractMultipleKeyValueFiltering(t *testing.T) {
	tests := []struct {
		name string
		str  string
		res  []map[string]string
		err  error
	}{
		{
			name: "Success",
			str:  "{{.Name}}=target1,{{.Age}}=38",
			res: []map[string]string{
				map[string]string{
					"key":   "{{.Name}}",
					"value": "target1",
				},
				map[string]string{
					"key":   "{{.Age}}",
					"value": "38",
				},
			},
			err: nil,
		},
		{
			name: "FailWithoutSeparator",
			str:  "{{.Name}}=target1{{.Age}}=38",
			res:  nil,
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailBadFormat",
			str:  "{{.Name}}=target1{{.Age}}38",
			res:  nil,
			err:  fmt.Errorf("expected error"),
		},
	}
	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			res, err := extractMultipleKeyValueFiltering(utest.str)
			if err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, utest.res, res)
		})
	}
}

func TestMultipleMetaFilter(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		object   *runtime.ObjectMeta
		expected bool
		err      error
	}{
		{
			name: "SuccessOneFilter",
			str:  "{{.Name}}=hello",
			object: &runtime.ObjectMeta{
				Name: "hello",
				UID:  "123",
			},
			expected: true,
			err:      nil,
		},
		{
			name: "SuccessTwoFilter",
			str:  "{{.Name}}=hello,{{.UID}}=123",
			object: &runtime.ObjectMeta{
				Name: "hello",
				UID:  "123",
			},
			expected: true,
			err:      nil,
		},
		{
			name: "SuccessOneValueDiffer",
			str:  "{{.Name}}=hello,{{.UID}}=1234",
			object: &runtime.ObjectMeta{
				Name: "hello",
				UID:  "123",
			},
			expected: false,
			err:      nil,
		},
		{
			name: "FailBadFormat",
			str:  "{{.Name}}=hello,{{.Unexisting}}=1234",
			object: &runtime.ObjectMeta{
				Name: "hello",
				UID:  "123",
			},
			expected: false,
			err:      fmt.Errorf("expected error"),
		},
	}

	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			mmf, err := GenerateMultipleMetadataFiltering(utest.str)
			expected, err := mmf.AreExpected(utest.object)
			if utest.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, utest.expected, expected)
		})
	}
}
