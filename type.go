package entchi

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
)

type extension struct {
	entc.DefaultExtension
	hooks []gen.Hook
	data  data
}

type Driver = string

const (
	SQLite     Driver = "sqlite3"
	MySQL      Driver = "mysql"
	PostgreSQL Driver = "postgres"
)

type option = func(*extension)

type data struct {
	*gen.Graph
	DBConfig      *DBConfig
	TSConfig      *TSConfig
	EntchiConfig  *EntchiConfig
	CurrentSchema *load.Schema
	IDType        string
}

type DBConfig struct {
	Path   string
	Driver string
	Dsn    string
}

type TSConfig struct {
	TypesPath string
	ApiPath   string
	Runtime   string
	IDType    string
}

type EntchiConfig struct {
	HandlersPath string
	RoutesPath   string
}

type comparable interface {
	~string | ~int | ~float32 | ~uint
}

var gots = map[string]string{
	"time.Time": "string",
	"bool":      "boolean",
	"int":       "number",
	"uint":      "number",
	"float":     "number",
	"enum":      "string",
	"any":       "any",
	"other":     "any",
	"json":      "any",
}

type skipAnnotation struct {
	name  string
	Skips []uint
}