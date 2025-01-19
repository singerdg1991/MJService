package drivers

type Mongo struct {
}

func (m *Mongo) MigrateUp() error {
	return nil
}
