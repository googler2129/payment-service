package static

type Model struct {
	ID string `gorm:"column:id;primaryKey:false"`
}

func (m *Model) SetID(id string) {
	m.ID = id
}
