package database

import "distribution-system-be/models/dto"

// GetQtyOrd ...
func GetQtyOrd(dateStart, dateEnd string) (dto.DBoardQtyOrderDto, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var dashqo dto.DBoardQtyOrderDto
	var dashq []dto.DBoardQtyDto
	var err error

	err = db.Raw(`
	SELECT COUNT(order_no) AS order_count, internal_status
	FROM public.order
	WHERE order_date between ? and ?
	GROUP BY internal_status
	`, dateStart, dateEnd).Scan(&dashq).Error

	for _, q := range dashq {
		switch q.InternalStatus {
		case 0:
			dashqo.NewOrder = q.OrderCount
		case 1:
			dashqo.Order = q.OrderCount
		case 2:
			dashqo.Payment = q.OrderCount
		case 3:
			dashqo.Complete = q.OrderCount
		case 4:
			dashqo.Reject = q.OrderCount
		}
	}

	return dashqo, err
}
