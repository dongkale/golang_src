package controllers

import "testing"

func TestNoticeFCM(t *testing.T) {
	type args struct {
		gbn      string
		mem_cd   string
		cont     string
		sn       int64
		brdgbncd string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NoticeFCM(tt.args.gbn, tt.args.mem_cd, tt.args.cont, tt.args.sn, tt.args.brdgbncd)
		})
	}
}
