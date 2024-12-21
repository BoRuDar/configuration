package configuration

import (
	"reflect"
	"testing"
)

func Test_fetchTagKey(t *testing.T) {
	t.Parallel()

	registredTags := map[string]struct{}{
		"json":    {},
		"xml":     {},
		"flag":    {},
		"default": {},
	}

	tests := []struct {
		name string
		in   reflect.StructTag
		want map[string]struct{}
	}{
		{
			name: "empty tags",
			in:   reflect.StructTag(""),
			want: map[string]struct{}{},
		},
		{
			name: "empty tag value",
			in:   reflect.StructTag(`json:""`),
			want: map[string]struct{}{
				"json": {},
			},
		},
		{
			name: "non-empty tag value",
			in:   reflect.StructTag(`default:"one;two"`),
			want: map[string]struct{}{
				"default": {},
			},
		},
		{
			name: "multiple tags",
			in:   reflect.StructTag(`json:"id" xml:"ID"`),
			want: map[string]struct{}{
				"json": {},
				"xml":  {},
			},
		},
		{
			name: "malformed tag",
			in:   reflect.StructTag(`json`),
			want: map[string]struct{}{},
		},
		{
			name: "tag with spaces in value",
			in:   reflect.StructTag(`flag:"name_flag||Some description" json:"id" xml:"ID"`),
			want: map[string]struct{}{
				"flag": {},
				"json": {},
				"xml":  {},
			},
		},
		{
			name: "special characters in tag value",
			in:   reflect.StructTag(`default:"one;two-three_1,/ and this!=2*7&^3:,.	$"`),
			want: map[string]struct{}{
				"default": {},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert(t, tt.want, fetchTagKey(tt.in, registredTags))
		})
	}
}
