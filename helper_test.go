package configuration

import (
	"reflect"
	"testing"
)

func Test_fetchTagKey(t *testing.T) {
	t.Parallel()

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
			in:   reflect.StructTag(`json:"id"`),
			want: map[string]struct{}{
				"json": {},
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
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert(t, tt.want, fetchTagKey(tt.in))
		})
	}
}
