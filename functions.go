package entchi

import (
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"entgo.io/ent/schema/field"
)

var (
	snake      = gen.Funcs["snake"].(func(string) string)
	plural     = gen.Funcs["plural"].(func(string) string)
	buggyCamel = gen.Funcs["camel"].(func(string) string)
	camel      = func(s string) string { return buggyCamel(snake(s)) }
)

func init() {
	gen.Funcs["tag"] = tag
	gen.Funcs["imports"] = imports
	gen.Funcs["null_field_create"] = null_field_create
	gen.Funcs["null_field_update"] = null_field_update
	gen.Funcs["extract_type"] = extract_type
	gen.Funcs["edge_field"] = edge_field
	gen.Funcs["is_comparable"] = is_comparable
	gen.Funcs["enum_or_edge_filed"] = enum_or_edge_filed
	gen.Funcs["get_name"] = get_name
	gen.Funcs["get_type"] = get_type
	gen.Funcs["is_slice"] = is_slice
	gen.Funcs["id_type"] = id_type
	gen.Funcs["go_ts"] = go_ts
	gen.Funcs["order_fields"] = order_fields
	gen.Funcs["select_fields"] = select_fields
	gen.Funcs["dir"] = path.Dir
	gen.Funcs["skip_field_create"] = skip_field_create
	gen.Funcs["skip_field_update"] = skip_field_update
	gen.Funcs["skip_field_query"] = skip_field_query
	gen.Funcs["skip_field_type"] = skip_field_type
	gen.Funcs["skip_schema_query"] = skip_schema_query
	gen.Funcs["skip_schema_create"] = skip_schema_create
	gen.Funcs["skip_schema_update"] = skip_schema_update
	gen.Funcs["skip_schema_delete"] = skip_schema_delete
	gen.Funcs["skip_edge_query"] = skip_edge_query
	gen.Funcs["skip_handler_find"] = skip_handler_find
	gen.Funcs["skip_handler_find_many"] = skip_handler_find_many
	gen.Funcs["skip_handler_create"] = skip_handler_create
	gen.Funcs["skip_handler_create_many"] = skip_handler_create_many
	gen.Funcs["skip_handler_update"] = skip_handler_update
	gen.Funcs["skip_handler_update_many"] = skip_handler_update_many
	gen.Funcs["skip_handler_delete"] = skip_handler_delete
	gen.Funcs["skip_handler_delete_many"] = skip_handler_delete_many
}

func tag(f *load.Field) string {
	if f.Tag == "" {
		name := camel(f.Name)
		if strings.HasSuffix(name, "ID") {
			name = strings.TrimSuffix(name, "ID")
			name += "Id"
		}
		return fmt.Sprintf("json:\"%s,omitempty\"", name)
	}
	return f.Tag
}

func imports(g *gen.Graph, isInput ...bool) []string {
	imps := []string{}
	for _, s := range g.Schemas {
		for _, f := range s.Fields {
			if len(f.Enums) > 0 && len(isInput) > 0 && isInput[0] {
				imps = append(imps, path.Join(g.Package, strings.Split(f.Info.Ident, ".")[0]))
			}
			if f.Info != nil && len(f.Info.PkgPath) != 0 {
				if !in(f.Info.PkgPath, imps) {
					imps = append(imps, f.Info.PkgPath)
				}
			}
		}
	}
	return imps
}

func null_field_create(f *load.Field) bool {
	return f.Optional || f.Default
}

func null_field_update(field *load.Field) bool {
	return !strings.HasPrefix(extract_type(field), "[]")
}

func extract_type(field *load.Field) string {
	if field.Info.Ident != "" {
		return field.Info.Ident
	}
	return field.Info.Type.String()
}

func edge_field(e *load.Edge) bool {
	return e.Field != ""
}

func is_comparable(f *load.Field) bool {
	return has_prefixes(extract_type(f), []string{
		"string",
		"int",
		"uint",
		"float",
		"time.Time",
	})
}

func enum_or_edge_filed(s *load.Schema, f *load.Field) bool {
	for _, e := range s.Edges {
		if e.Field == f.Name {
			return extract_type(f) == "enum"
		}
	}
	return false
}

func get_name(f *load.Field) string {
	n := camel(f.Name)
	if strings.HasSuffix(n, "ID") {
		n = strings.TrimSuffix(n, "ID") + "Id"
	}
	return n
}

func get_type(t *field.TypeInfo) string {
	return go_ts(t.Type.String())
}

func go_ts(s string) string {
	slice := false
	if strings.HasPrefix(s, "[]") {
		slice = true
		s = strings.TrimPrefix(s, "[]")
	}
	for k, v := range gots {
		if strings.HasPrefix(s, k) {
			if slice {
				return v + "[]"
			}
			return v
		}
	}
	if slice {
		return s + "[]"
	}
	return s
}

func is_slice(f *load.Field) bool {
	return strings.HasPrefix(get_type(f.Info), "[]")
}

func id_type(s *load.Schema) string {
	for _, f := range s.Fields {
		if strings.ToLower(f.Name) == "id" {
			return get_type(f.Info)
		}
	}
	return "number"
}

func order_fields(s *load.Schema) string {
	fields := []string{}
	for _, f := range s.Fields {
		if orderable(f) {
			fields = append(fields, get_name(f))
		}
	}
	return "\"" + strings.Join(fields, "\" | \"") + "\""
}

func select_fields(s *load.Schema) string {
	fields := []string{}
	for _, f := range s.Fields {
		fields = append(fields, get_name(f))
	}
	return "\"" + strings.Join(fields, "\" | \"") + "\""
}

func orderable(f *load.Field) bool {
	return has_prefixes(extract_type(f), []string{
		"string",
		"int",
		"uint",
		"float",
		"time.Time",
		"bool",
	})
}

func skip_field_create(f *load.Field) bool {
	return (f.Default && f.Name == "id") || shouldSkip(f.Annotations, SkipFieldCreate)
}

func skip_field_update(f *load.Field) bool {
	return f.Immutable || (f.Default && f.Name == "id") || shouldSkip(f.Annotations, SkipFieldUpdate)
}

func skip_field_query(f *load.Field) bool {
	return shouldSkip(f.Annotations, SkipFieldQuery)
}

func skip_field_type(f *load.Field) bool {
	return shouldSkip(f.Annotations, SkipFieldType)
}

func skip_schema_query(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipSchemaQuery)
}

func skip_schema_create(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipSchemaCreate)
}

func skip_schema_update(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipSchemaUpdate)
}

func skip_schema_delete(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipSchemaDelete)
}

func skip_edge_query(e *load.Edge) bool {
	return shouldSkip(e.Annotations, SkipEdgeQuery)
}

func skip_handler_find(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerFind)
}

func skip_handler_find_many(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerFindMany)
}

func skip_handler_create(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerCreate)
}

func skip_handler_create_many(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerCreateMany)
}

func skip_handler_update(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerUpdate)
}

func skip_handler_update_many(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerUpdateMany)
}

func skip_handler_delete(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerDelete)
}

func skip_handler_delete_many(s *load.Schema) bool {
	return shouldSkip(s.Annotations, SkipHandlerDeleteMany)
}

func shouldSkip(annotations map[string]any, skip uint) bool {
	a := &skipAnnotation{}
	a.decode(annotations[skipAnootationName])
	if in(skip, a.Skips) || in(SkipAll, a.Skips) {
		return true
	}
	return false
}
