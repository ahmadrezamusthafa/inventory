package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rezamusthafa/inventory/api/services/inputs"
)

func TestGenerateSKU(t *testing.T) {
	type args struct {
		product inputs.Product
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal",
			args: args{
				product: inputs.Product{
					Name:  "Baju kebaya",
					Size:  "L",
					Color: "Yellow",
				},
			},
			want: "-LL-YEL",
		},
		{
			name: "Normal",
			args: args{
				product: inputs.Product{
					Name:  "Baju kebaya",
					Size:  "XXL",
					Color: "blue",
				},
			},
			want: "-XX-BLU",
		},
		{
			name: "Normal",
			args: args{
				product: inputs.Product{
					Name:  "Baju kebaya",
					Size:  "M",
					Color: "blue",
				},
			},
			want: "-MM-BLU",
		},
		{
			name: "Failed",
			args: args{
				product: inputs.Product{
					Name:  "Baju kebaya",
					Size:  "",
					Color: "blue",
				},
			},
			want: "",
		},
		{
			name: "Normal",
			args: args{
				product: inputs.Product{
					Name:  "Baju kebaya",
					Size:  "M",
					Color: "B",
				},
			},
			want: "-MM-BBB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSKU(tt.args.product)
			if !strings.Contains(got, tt.want) {
				t.Errorf("GenerateSKU() = %v, want %v", got, tt.want)
			}

			fmt.Println(got)
		})
	}
}

func TestIsValidSKU(t *testing.T) {
	type args struct {
		sku string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Normal",
			args: args{
				sku: "SSI-D14320263-MM-BLU",
			},
			want: true,
		},
		{
			name: "Normal",
			args: args{
				sku: "SSI-D41483963-XX-BLU",
			},
			want: true,
		},
		{
			name: "Normal",
			args: args{
				sku: "SSI-D31331013-LL-YEL",
			},
			want: true,
		},
		{
			name: "Failed",
			args: args{
				sku: "23S-D31331013-SLL-YEL",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSKU(tt.args.sku); got != tt.want {
				t.Errorf("IsValidSKU() = %v, want %v", got, tt.want)
			}
		})
	}
}
