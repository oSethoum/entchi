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

		if e.data.EntchiConfig != nil {
			s = parseTemplate("Entchi/routes", e.data)
			writeFile(path.Join(e.data.EntchiConfig.RoutesPath, "routes.go"), s)

			s = parseTemplate("Entchi/util", e.data)
			writeFile(path.Join(e.data.EntchiConfig.HandlersPath, "util.go"), s)

			for _, schema := range g.Schemas {
				if skip_schema_query(schema) && skip_schema_create(schema) &&
					skip_schema_update(schema) && skip_schema_delete(schema) {
					continue
				}
				e.data.CurrentSchema = schema
				s := parseTemplate("Entchi/handler", e.data)
				writeFile(path.Join(e.data.EntchiConfig.HandlersPath, snake(plural(schema.Name))+".go"), s)
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