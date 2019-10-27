package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToArgv(t *testing.T) {
	tests := []struct {
		name string
		give string
		want []string
	}{
		{
			name: "command without quotes",
			give: "_test a b c",
			want: []string{"_test", "a", "b", "c"},
		},
		{
			name: "command with quotes",
			give: "_test \"my long\" \"command with quotes\"",
			want: []string{"_test", "my long", "command with quotes"},
		},
		{
			name: "different type of quotes",
			give: "_test “my long“ 'command with quotes'",
			want: []string{"_test", "my long", "command with quotes"},
		},
		{
			name: "mixed quotes and no quotes",
			give: "_test \"my long\" abc \"command with quotes\"",
			want: []string{"_test", "my long", "abc", "command with quotes"},
		},
		{
			name: "unfinished quotes",
			give: "_test \"my long\" abc \"command that's not finished",
			want: []string{"_test", "my long", "abc", "command that's not finished"},
		},
		{
			name: "escape character at the end",
			give: "_test \"my long\" abc \"command with quotes\"\\",
			want: []string{"_test", "my long", "abc", "command with quotes"},
		},
		{
			name: "quote at the beginning",
			give: "\"_test\" \"my long\" abc \"command with quotes\"",
			want: []string{"_test", "my long", "abc", "command with quotes"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := ToArgv(tt.give)
			assert.Equal(t, tt.want, got)
		})
	}
}
