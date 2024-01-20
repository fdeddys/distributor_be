package dbmodels

// Sequence ...
type Sequence struct {
	ID      int64  `gorm:"column:id"`
	Subject string `gorm:"column:subj"`
	Year    string `gorm:"column:year"`
	Month   string `gorm:"column:month" `
	Seq     int32  `gorm:"column:last_seq"`
}

// TableName ...
func (t *Sequence) TableName() string {
	return "public.tb_sequence"
}
