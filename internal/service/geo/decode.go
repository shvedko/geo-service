package geo

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func DecodeContext(ctx context.Context, v any) error {
	o := reflect.ValueOf(v)
	if o.Kind() != reflect.Pointer || o.IsNil() {
		return fmt.Errorf("not a pointer or nil pointer")
	}

	e := o.Elem()

	switch e.Kind() {
	case reflect.Struct:
		t := e.Type()

		for i := 0; i < e.NumField(); i++ {
			f := e.Field(i)
			m := t.Field(i)

			k := strings.SplitN(m.Tag.Get("json"), ",", 2)[0]
			if k == "-" {
				continue
			}
			if k == "" {
				k = m.Name
			}

			s := chi.URLParamFromCtx(ctx, k)
			if s == "" {
				continue
			}

			if f.CanSet() {
				switch f.Kind() {
				case reflect.Int32:
					x, err := strconv.ParseInt(s, 10, 32)
					if err != nil {
						return fmt.Errorf("%s is invalid: %w", k, err)
					}
					f.SetInt(x)
				case reflect.Float64:
					x, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("%s is invalid: %w", k, err)
					}
					f.SetFloat(x)
				case reflect.String:
					f.SetString(s)
				default:
					return fmt.Errorf("type not implemented yet: %s", f.Type())
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("not a pointer to struct")
	}
}

func Decode(r *http.Request, v any) error {
	return DecodeContext(r.Context(), v)
}
