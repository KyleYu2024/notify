package pluginsdk

import (
	"context"
	"fmt"
	"reflect"
)

type Adapter struct {
	target any
}

func WrapPlugin(target any) (Plugin, error) {
	val := reflect.ValueOf(target)
	typ := val.Type()

	// 检查必须的方法是否存在
	requiredMethods := []struct {
		name   string
		numOut int
	}{
		{"ID", 1},
		{"Name", 1},
		{"Version", 1},
		{"Desc", 1},
		{"DefaultSettings", 1},
		{"Process", 2}, // (*Output, error)
	}
	for _, m := range requiredMethods {
		method, ok := typ.MethodByName(m.name)
		if !ok {
			return nil, fmt.Errorf("target missing method: %s", m.name)
		}
		if method.Type.NumOut() != m.numOut {
			return nil, fmt.Errorf("method %s has wrong return count", m.name)
		}
	}

	return &Adapter{target: target}, nil
}

func (a *Adapter) callString(name string) string {
	m := reflect.ValueOf(a.target).MethodByName(name)
	if !m.IsValid() {
		return ""
	}
	out := m.Call(nil)
	if len(out) == 0 {
		return ""
	}
	return out[0].String()
}

func (a *Adapter) ID() string      { return a.callString("ID") }
func (a *Adapter) Name() string    { return a.callString("Name") }
func (a *Adapter) Version() string { return a.callString("Version") }
func (a *Adapter) Desc() string    { return a.callString("Desc") }

func (a *Adapter) DefaultSettings() map[string]any {
	m := reflect.ValueOf(a.target).MethodByName("DefaultSettings")
	if !m.IsValid() {
		return nil
	}
	out := m.Call(nil)
	if len(out) == 0 || out[0].IsNil() {
		return nil
	}
	return out[0].Interface().(map[string]any)
}

func (a *Adapter) Process(ctx context.Context, input map[string]any, settings map[string]any) (*Output, error) {
	m := reflect.ValueOf(a.target).MethodByName("Process")
	if !m.IsValid() {
		return nil, fmt.Errorf("no Process method")
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(input),
		reflect.ValueOf(settings),
	}
	out := m.Call(args)

	var res *Output
	var err error
	if !out[0].IsNil() {
		res, err = convertOutput(out[0].Interface())
		if err != nil {
			return nil, err
		}
	}
	if !out[1].IsNil() {
		err = out[1].Interface().(error)
		if err != nil {
			return nil, err
		}
	}
	return res, err
}

// convertOutput 用反射将任意来源的 Output 转成本包的 Output
func convertOutput(v any) (*Output, error) {
	if v == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	res := &Output{}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		val := rv.Field(i)

		switch field.Name {
		case "Title":
			res.Title = val.String()
		case "Content":
			res.Content = val.String()
		case "Image":
			res.Image = val.String()
		case "URL":
			res.URL = val.String()
		case "Targets":
			if !val.IsNil() {
				res.Targets = make([]string, val.Len())
				for j := 0; j < val.Len(); j++ {
					res.Targets[j] = val.Index(j).String()
				}
			}
		case "Meta":
			if !val.IsNil() {
				meta, err := convertMeta(val.Interface())
				if err != nil {
					return nil, err
				}
				res.Meta = meta
			}
		}
	}

	return res, nil
}

// convertMeta 用反射转换 MetaData
func convertMeta(v any) (*MetaData, error) {
	if v == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	meta := &MetaData{}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		val := rv.Field(i)

		switch field.Name {
		case "Req":
			if !val.IsNil() {
				m := make(map[string]any)
				for _, key := range val.MapKeys() {
					m[key.String()] = val.MapIndex(key).Interface()
				}
				meta.Req = m
			}
		case "PluginID":
			meta.PluginID = val.String()
		case "ProcessedAt":
			meta.ProcessedAt = val.String()
		case "Extra":
			if !val.IsNil() {
				m := make(map[string]any)
				for _, key := range val.MapKeys() {
					m[key.String()] = val.MapIndex(key).Interface()
				}
				meta.Extra = m
			}
		}
	}

	return meta, nil
}
