package string

import (
	"github.com/stretchr/testify/assert"
	constants "github.com/vinicius.csantos/nagini-api/internal/util/constant"
	"testing"
)

func TestFormatGenericMessagesString(t *testing.T) {
	type args struct {
		strGeneric string
		strTarget  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should format generic message with target",
			args: args{strGeneric: constants.GenericCreateSuccessMessage, strTarget: "Test"},
			want: "Test created with successful",
		},

		{
			name: "Should format generic message without target",
			args: args{strGeneric: constants.GenericInternalServerErrorMessage, strTarget: ""},
			want: "It was not possible to perform this action",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			genericMessage := FormatGenericMessagesString(test.args.strGeneric, test.args.strTarget)

			assert.Equalf(t, test.want, genericMessage, test.name)

		})
	}
}

func TestRemoveSpecialCharacters(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should remove all special characters",
			args: args{str: "(T.e.s.t-S.p.e.c.i.a.l / C.h.a.r.a.c.t.e.r.s)"},
			want: "TestSpecialCharacters",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			removedCharacters := RemoveSpecialCharacters(test.args.str)

			assert.Equalf(t, test.want, removedCharacters, test.name)

		})
	}
}
