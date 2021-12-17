package cpfCNPJ

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCpfCNPJ(t *testing.T) {
	type args struct {
		cpfCNPJ string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "[CPF] Should validate CpfCnpj with successful",
			args: args{cpfCNPJ: "79873939083"},
			want: true,
		},

		{
			name: "[CPF] Should first CpfCnpj digit is false",
			args: args{cpfCNPJ: "79873939093"},
			want: false,
		},

		{
			name: "[CPF] Should second CpfCnpj digit is false",
			args: args{cpfCNPJ: "79873939080"},
			want: false,
		},

		{
			name: "[CPF] Should first CpfCnpj digit is zero",
			args: args{cpfCNPJ: "72626141003"},
			want: true,
		},

		{
			name: "[CPF] Should second CpfCnpj digit is zero",
			args: args{cpfCNPJ: "36033663870"},
			want: true,
		},

		{
			name: "[CNPJ] Should validate CpfCnpj with successful",
			args: args{cpfCNPJ: "56611458000131"},
			want: true,
		},

		{
			name: "[CNPJ] Should first CpfCnpj digit is false",
			args: args{cpfCNPJ: "56611458000181"},
			want: false,
		},

		{
			name: "[CNPJ] Should second CpfCnpj digit is false",
			args: args{cpfCNPJ: "56611458000138"},
			want: false,
		},

		{
			name: "[CNPJ] Should first CpfCnpj digit is zero",
			args: args{cpfCNPJ: "46557577000108"},
			want: true,
		},

		{
			name: "[CNPJ] Should second CpfCnpj digit is zero",
			args: args{cpfCNPJ: "48814815000130"},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testValidate := ValidateCpfCNPJ(test.args.cpfCNPJ)
			assert.Equalf(t, test.want, testValidate, test.name)
		})
	}
}
