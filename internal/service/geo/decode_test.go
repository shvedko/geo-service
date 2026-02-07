package geo

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/shvedko/geo-service/internal/repository"
)

func TestDecodeContext(t *testing.T) {
	type args struct {
		ctx context.Context
		v   any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				ctx: context.WithValue(context.TODO(), chi.RouteCtxKey, func() any {
					ctx := chi.NewRouteContext()

					ctx.URLParams.Add("id", "123")
					ctx.URLParams.Add("name", "A")
					ctx.URLParams.Add("lon", "11.22")
					ctx.URLParams.Add("lat", "33.44")

					return ctx
				}()),
				v: &repository.GetPointRow{},
			},
			want: &repository.GetPointRow{
				ID:   123,
				Name: "A",
				Lon:  11.22,
				Lat:  33.44,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeContext(tt.args.ctx, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DecodeContext() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, tt.args.v) {
				t.Errorf("DecodeContext() %v != %v", tt.want, tt.args.v)
			}
		})
	}
}
