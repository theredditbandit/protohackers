package main

import (
	"testing"
)

func Test_isPrime(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "even", args: args{n: 4}, want: false},
		{name: "even2", args: args{n: 10}, want: false},
		{name: "even3", args: args{n: 18}, want: false},
		{name: "odd", args: args{n: 5}, want: true},
		{name: "odd2", args: args{n: 7}, want: true},
		{name: "odd3", args: args{n: 9}, want: false},
		{name: "odd4", args: args{n: 37}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrime(tt.args.n); got != tt.want {
				t.Errorf("isPrime() = %v, want %v", got, tt.want)
			}
		})
	}
}
