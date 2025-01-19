package ports

import "github.com/hoitek/Maja-Service/cmd/migrator/drivers"

type Migrator interface {
	*drivers.Postgres | *drivers.Mongo | any
}
