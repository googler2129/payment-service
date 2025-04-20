package scd

// SCDModel is the base model for all SCD tables
type SCDModel struct {
	ID       string `gorm:"column:id;primaryKey:false"`
	Version  int    `gorm:"column:version;primaryKey:false"`
	UID      string `gorm:"column:uid;primaryKey"`
	IsLatest bool   `gorm:"column:is_latest;not null" json:"is_latest"`
}

func (m SCDModel) GetID() string {
	return m.ID
}

func (m SCDModel) GetVersion() int {
	return m.Version
}

func (m SCDModel) GetUID() string {
	return m.UID
}

func (m SCDModel) GetIsLatest() bool {
	return m.IsLatest
}

func (m *SCDModel) SetVersion(version int) {
	m.Version = version
}

func (m *SCDModel) SetUID(uid string) {
	m.UID = uid
}

func (m *SCDModel) SetID(id string) {
	m.ID = id
}

func (m *SCDModel) SetIsLatest(isLatest bool) {
	m.IsLatest = isLatest
}
