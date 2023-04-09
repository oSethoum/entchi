package entchi

import (
	"path"

	"entgo.io/ent/entc/gen"
)

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		e.data.Graph = g

		s := parseTemplate("ent/input", e.data)
		writeFile("ent/input.go", s)

		s = parseTemplate("ent/query", e.data)
		writeFile("ent/query.go", s)

		s = parseTemplate("ent/errors", e.data)
		writeFile("ent/errors.go", s)

		if e.data.ChiConfig != nil {
			s = parseTemplate("chi/routes", e.data)
			writeFile(path.Join(e.data.ChiConfig.RoutesPath, "routes.go"), s)

			s = parseTemplate("chi/util", e.data)
			writeFile(path.Join(e.data.ChiConfig.HandlersPath, "util.go"), s)

			for _, schema := range g.Schemas {
				if skip_schema_query(schema) && skip_schema_create(schema) &&
					skip_schema_update(schema) && skip_schema_delete(schema) {
					continue
				}
				e.data.CurrentSchema = schema
				s := parseTemplate("chi/handler", e.data)
				writeFile(path.Join(e.data.ChiConfig.HandlersPath, snake(plural(schema.Name))+".go"), s)
			}
		}

		if e.data.DBConfig != nil {
			s := parseTemplate("ent/db", e.data)
			writeFile(path.Join(e.data.DBConfig.Path, "db.go"), s)
		}

		if e.data.TSConfig != nil {
			s := parseTemplate("ts/api", e.data)
			writeFile(path.Join("ts/", "api.ts"), s)
			s = parseTemplate("ts/types", e.data)
			writeFile(path.Join("ts/", "types.ts"), s)
		}

		return next.Generate(g)
	})
}
