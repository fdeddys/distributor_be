package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

func ApprovePaymentDirect(paymentDirect *dto.PaymentDirectModel) (errCode string, errDesc string) {
	fmt.Println(" Approve Payment Direct ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order, _ := GetSalesOrderByOrderId(paymentDirect.SalesOrderID)
	// update stock -> jika penjualan kredit , klo tunai update saat payment
	// update history stock
	// hitung ulang
	var total float32
	total = 0
	salesOrderDetails := GetAllDataDetail(paymentDirect.SalesOrderID)
	for idx, orderDetail := range salesOrderDetails {
		fmt.Println("idx -> ", idx)

		// update stock pindah saat sales order
		// updateStockInsertHistory(order, orderDetail)

		total = total + (orderDetail.Price * float32(orderDetail.QtyOrder))
	}
	if order.Tax > 0 {
		total += total * order.Tax / 100
	}
	tx.Debug().LogMode(true)
	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	payment, _ := GetPaymentById(paymentDirect.PaymentID)
	payment.Status = constants.STATUS_APPROVE
	payment.LastUpdate = time.Now()
	payment.LastUpdateBy = dto.CurrUser
	payment.TotalPayment = total
	r := tx.Save(&payment)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}
	order.InvoiceNo = payment.PaymentNo
	o := tx.Save((&order))
	if o.Error != nil {
		fmt.Println("Error save inv no in order => ", o.Error.Error())
	}
	// fmt.Println("Order [database]=> order id", order.OrderNo)

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetOrderPage ...
func GetPaymentDirectPage(param dto.FilterPaymentDirect, offset, limit int) ([]dto.PaymentDirectModel, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var paymentDirects []dto.PaymentDirectModel
	var total int
	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&paymentDirects).Error
		if err != nil {
			return paymentDirects, 0, err
		}
		return paymentDirects, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan []dto.PaymentDirectModel)
	errCount := make(chan error)

	go asyncQuerysPaymentDirect(db, offset, limit, &paymentDirects, param, errQuery)
	go asyncQueryCountsPaymentDirect(db, &total, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	fmt.Println("total =>", *&total)
	if resErrQuery != nil {
		log.Println("errr-->", resErrCount)
		return paymentDirects, 0, resErrCount
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return paymentDirects, 0, resErrCount
	}
	return paymentDirects, total, nil
}

func getParamDirectPayment(param dto.FilterPaymentDirect) (salesOrderNo, paymentNo string, isSearchPaymentStatus, isSearchPaymentNo bool) {

	paymentNo = "%"
	salesOrderNo = "%"
	isSearchPaymentStatus = true
	isSearchPaymentNo = false
	// paymentNo = param.PaymentNo
	if param.PaymentNo != "" {
		paymentNo = "%" + param.PaymentNo + "%"
		isSearchPaymentStatus = true
	}
	if param.SalesOrderNo != "" {
		salesOrderNo = "%" + param.SalesOrderNo + "%"
	}
	if param.PaymentStatus == 0 {
		isSearchPaymentStatus = false
	}

	return
}

// AsyncQuerysOrders ...
func asyncQuerysPaymentDirect(db *gorm.DB, offset int, limit int, paymentDirects *[]dto.PaymentDirectModel, param dto.FilterPaymentDirect, resChan chan []dto.PaymentDirectModel) {

	salesOrderNo, paymentNo, isSearchPaymentStatus, isSearchPaymentNo := getParamDirectPayment(param)
	fmt.Println("Search sales order no = ", salesOrderNo, " payment no ", paymentNo, " search pay status", isSearchPaymentStatus)

	dateStart := param.StartDate
	dateEnd := param.EndDate
	// dateStart := param.StartDate + " 00:00:00.000"
	// dateEnd := param.EndDate + " 23:59:59.999"

	db.Raw("select p.status as payment_status , p.payment_no , so.sales_order_no , so.order_date , "+
		" so.status as so_status , so.grand_total  ,so.is_cash, p.id as payment_id, so.id as sales_order_id "+
		" 	from sales_order so "+
		" 	left join payment_order po on po.sales_order_id = so.id  "+
		" 	left join payment p on po.payment_id = p.id  "+
		" where so.status in ('20','40', '50', '60') "+
		" 	and so.is_cash "+
		"   and so.order_date >= ? and so.order_date <= ? "+
		" 	and so.sales_order_no like ?  "+
		" 	and ( p.payment_no like ? or ( not ? )) "+
		" 	and ( p.status = ?        or ( not ? )) "+
		" order by  so.id desc "+
		" limit ?  offset  ?", dateStart, dateEnd,
		salesOrderNo, paymentNo, isSearchPaymentNo, param.PaymentStatus, isSearchPaymentStatus, limit, offset).Scan(&paymentDirects)

	// var result []dto.PaymentDirectModel
	// for rows.Next() {

	// 	var paymentStatus int64
	// 	var paymentNo string
	// 	var salesOrderNo string
	// 	var orderDate time.Time
	// 	var soStatus int64
	// 	var grandTotal float32
	// 	rows.Scan(&paymentStatus, &paymentNo, &salesOrderNo, &orderDate, &soStatus, &grandTotal)
	// 	var paymentDirect dto.PaymentDirectModel

	// 	fmt.Println("isi->", paymentStatus, "paymentno", paymentNo)
	// 	paymentDirect.PaymentStatus = int8(paymentStatus)
	// 	result = append(result, paymentDirect)
	// 	// paymentDirects = append(paymentDirects, paymentDirect)
	// }

	// paymentDirects= append(paymentDirects, result)
	resChan <- nil
}

// asyncQueryCountsPaymentDirect ...
func asyncQueryCountsPaymentDirect(db *gorm.DB, total *int, param dto.FilterPaymentDirect, resChan chan error) {

	salesOrderNo, paymentNo, isSearchPaymentStatus, isSearchPaymentNo := getParamDirectPayment(param)
	fmt.Println("Search sales order no = ", salesOrderNo, " payment no ", paymentNo, " search pay status", isSearchPaymentStatus)

	// dateStart := param.StartDate + " 00:00:00.000"
	// dateEnd := param.EndDate + " 23:59:59.999"
	dateStart := param.StartDate
	dateEnd := param.EndDate

	row := db.Raw("select count(*) "+
		" 	from sales_order so "+
		" 	left join payment_order po on po.sales_order_id = so.id  "+
		" 	left join payment p on po.payment_id = p.id  "+
		" where so.status in ('20','40', '50') "+
		" 	and so.is_cash "+
		" 	and Date(so.order_date) >= ? and Date(so.order_date)<= ? "+
		" 	and so.sales_order_no like ?  "+
		" 	and ( p.payment_no like ? or ( not ? )) "+
		" 	and ( p.status = ?         or ( not ? )) ", dateStart, dateEnd,
		salesOrderNo, paymentNo, isSearchPaymentNo, param.PaymentStatus, isSearchPaymentStatus).Row()
	row.Scan(total)

	fmt.Println("total rec =", *total)
	resChan <- nil
}

func RejectPaymentDirect(paymentDirect *dto.PaymentDirectModel) (errCode string, errDesc string) {
	fmt.Println(" REjeCT Payment Direct ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order, _ := GetSalesOrderByOrderId(paymentDirect.SalesOrderID)
	// update stock -> jika penjualan kredit , klo tunai update saat payment
	// update history stock
	// hitung ulang
	var total float32
	total = 0
	salesOrderDetails := GetAllDataDetail(paymentDirect.SalesOrderID)
	for idx, orderDetail := range salesOrderDetails {
		fmt.Println("idx -> ", idx)

		updateStockInsertHistoryReject(order, orderDetail)

		total = total + (orderDetail.Price * float32(orderDetail.QtyOrder))
	}

	tx.Debug().LogMode(true)
	// r := db.Model(&newOrder).Where("id = ?", order.ID).Update(dbmodels.SalesOrder{OrderNo: order.OrderNo, StatusCode: "001", WarehouseCode: order.WarehouseCode, InternalStatus: 1, OrderDate: order.OrderDate})

	payment, _ := GetPaymentById(paymentDirect.PaymentID)
	payment.Status = constants.STATUS_REJECT
	payment.LastUpdate = time.Now()
	payment.LastUpdateBy = dto.CurrUser
	payment.TotalPayment = total
	r := tx.Save(&payment)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}
	order.Status = constants.STATUS_REJECT_PAYMENT
	o := tx.Save(&order)
	if o.Error != nil {
		fmt.Println("Error save inv no in order => ", o.Error.Error())
	}
	// fmt.Println("Order [database]=> order id", order.OrderNo)

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func updateStockInsertHistoryReject(order dbmodels.SalesOrder, orderDetail dbmodels.SalesOrderDetail) {
	fmt.Println("Sales order yang di reject =>", order)

	product, _, _ := FindProductByID(orderDetail.ProductID)

	checkStock, _, _ := GetStockByProductAndWarehouse(product.ID, order.WarehouseID)
	curQty := checkStock.Qty

	updateQty := curQty + orderDetail.QtyOrder

	fmt.Println("cur qty =", curQty, " update =", updateQty)
	var historyStock dbmodels.HistoryStock
	historyStock.Code = product.Code
	if order.IsCash {
		historyStock.Description = "Direct Sales REJECT"
	}
	historyStock.Hpp = product.Hpp
	historyStock.WarehouseID = order.WarehouseID
	historyStock.Name = product.Name
	historyStock.Price = orderDetail.Price
	historyStock.ReffNo = order.SalesOrderNo
	historyStock.TransDate = order.OrderDate
	historyStock.Debet = orderDetail.QtyOrder
	historyStock.Kredit = 0
	historyStock.Saldo = updateQty
	historyStock.Disc1 = orderDetail.Disc1 * orderDetail.Price / 100
	historyStock.LastUpdate = time.Now()
	historyStock.LastUpdateBy = dto.CurrUser

	UpdateStockProductByID(orderDetail.ProductID, updateQty, order.WarehouseID)
	SaveHistory(historyStock)
}
