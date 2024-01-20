package database

import (
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetAllDataDetailReceive ...
func GetAllDataDetailReceive(receiveID int64) []dbmodels.ReceiveDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var receiveDetails []dbmodels.ReceiveDetail

	db.Preload("Product").Preload("UOM").Find(&receiveDetails, " receive_id = ? and qty > 0 ", receiveID)

	return receiveDetails
}

// GetReceiveDetailPage ...
func GetReceiveDetailPage(param dto.FilterReceiveDetail, offset, limit int) ([]dbmodels.ReceiveDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var receiveDetails []dbmodels.ReceiveDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&receiveDetails).Error
		if err != nil {
			return receiveDetails, 0, err
		}
		return receiveDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReceiveDetails(db, offset, limit, &receiveDetails, param.ReceiveID, errQuery)
	go AsyncQueryCountsReceiveDetails(db, &total, param.ReceiveID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return receiveDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return receiveDetails, 0, resErrCount
	}
	return receiveDetails, total, nil
}

// AsyncQueryCountsReceiveDetails ...
func AsyncQueryCountsReceiveDetails(db *gorm.DB, total *int, receiveID int64, offset int, limit int, resChan chan error) {

	var err error

	err = db.Model(&dbmodels.ReceiveDetail{}).Offset(offset).Where("receive_id = ?", receiveID).Count(total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysReceiveDetails ...
func AsyncQuerysReceiveDetails(db *gorm.DB, offset int, limit int, receiveDetails *[]dbmodels.ReceiveDetail, receiveID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Order("product.name ASC")
		}).
		Preload("Product.BigUom").
		Preload("Product.SmallUom").
		Preload("UOM").
		Order("id desc").
		Find(&receiveDetails, "receive_id = ? ", receiveID).
		Error

	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("order--> ", receiveDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SaveReceiveDetail ...
func SaveReceiveDetail(receiveDetail *dbmodels.ReceiveDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Receive Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Save(&receiveDetail); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
		return
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", receiveDetail.ID)
	return

}

//UpdateReceiveDetail ...
func UpdateReceiveDetail(idDetail, qty int64, price, disc1, disc2 float32, batchNo, ed string) (errCode string, errDesc string) {

	fmt.Println(" Update Receive Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(dbmodels.ReceiveDetail{}).Where("id = ?", idDetail).Updates(
		dbmodels.ReceiveDetail{
			Qty:          qty,
			Price:        price,
			Disc1:        disc1,
			Disc2:        disc2,
			LastUpdate:   util.GetCurrDate(),
			LastUpdateBy: dto.CurrUser,
			BatchNo:      batchNo,
			Ed:           ed,
		})
	if r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
		return
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", r.RowsAffected)
	return

}

// DeleteReceiveDetailById ...
func DeleteReceiveDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete Receive Detail  ---- ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Where("id = ? ", id).Delete(dbmodels.ReceiveDetail{}); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
		return
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", id)
	return

}

// GetReceiveOrderDetailPage ...
func GetReceiveOrderDetailBatchExpiredPage(param dto.FilterBatchExpired, offset, limit int) ([]dbmodels.ReceiveDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var ReceiveOrderDetails []dbmodels.ReceiveDetail
	var total int

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysReceiveOrderDetailsBatchExpired(db, offset, limit, &ReceiveOrderDetails, param, errQuery)
	go AsyncQueryCountsReceiveOrderDetailsBatchExpired(db, &total, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return ReceiveOrderDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return ReceiveOrderDetails, 0, resErrCount
	}
	return ReceiveOrderDetails, total, nil
}

// AsyncQuerysReceiveOrders ...
func AsyncQuerysReceiveOrderDetailsBatchExpired(db *gorm.DB, offset int, limit int, ReceiveOrderDetails *[]dbmodels.ReceiveDetail, param dto.FilterBatchExpired, resChan chan error) {

	var err error
	byDate := true
	if param.ExpiredStart == "" || param.ExpiredEnd == "" {
		byDate = false
	}
	batch := "%" + param.Batch + "%"
	productName := "%" + param.ProductName + "%"

	// err = db.Offset(offset).Limit(limit).Preload("Product.SmallUom").Preload("UOM").Find(&ReceiveOrderDetails, " ( batch_no ilike ? ) and ( (TO_DATE(ed,'YYYY-MM-DD') between TO_DATE(?, 'YYYY-MM-DD') and TO_DATE(?, 'YYYY-MM-DD') ) or ( not ?))  ", batch, param.ExpiredStart, param.ExpiredEnd, byDate).Error

	err = db.Offset(offset).
		Limit(limit).
		// Preload("Product", "Name ilike ? ", productName).
		Preload("Product.SmallUom").
		Preload("UOM").
		Joins("Join product on product.id = receive_detail.product_id ").
		Where("product.name ilike ? ", productName).
		Find(&ReceiveOrderDetails, " ( batch_no ilike ? ) and ( ( ed between ? and ? ) or ( not ?))  ", batch, param.ExpiredStart, param.ExpiredEnd, byDate).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQueryCountsReceiveOrders ...
func AsyncQueryCountsReceiveOrderDetailsBatchExpired(db *gorm.DB, total *int, param dto.FilterBatchExpired, resChan chan error) {

	var ReceiveOrderDetails []dbmodels.ReceiveDetail

	byDate := true
	if param.ExpiredStart == "" || param.ExpiredEnd == "" {
		byDate = false
	}
	batch := "%" + param.Batch + "%"
	productName := "%" + param.ProductName + "%"
	var err error
	// err = db.Model(&ReceiveOrderDetails).Where(" ( batch_no ilike ? ) and ( (TO_DATE(ed,'YYYY-MM-DD') between TO_DATE(?,'YYYY-MM-DD' ) and TO_DATE(?, 'YYYY-MM-DD') ) or ( not ?))  ", batch, param.ExpiredStart, param.ExpiredEnd, byDate).Count(&*total).Error

	err = db.Model(&ReceiveOrderDetails).
		Joins("Join product on product.id = receive_detail.product_id ").Where("product.name ilike ? ", productName).
		Where(" ( batch_no ilike ? ) and ( (ed  between ? and ? ) or ( not ?))  ", batch, param.ExpiredStart, param.ExpiredEnd, byDate).Count(&*total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

type Result struct {
	Harga int64
}

// GetReceiveDetailPage ...
func GetDataPriceProduct(productId int64) (res models.ResponseReceiveCheckPrice) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	// recDetail := dbmodels.ReceiveDetail{}

	fmt.Println(" product Receive Detail  ---- ", productId)

	// success
	// var receiveDetails []dbmodels.ReceiveDetail
	// db.Preload("Product", "id = ?", productId).Find(&receiveDetails)

	// var result Result
	// db.Raw(` select price
	// 	 from "public"."receive_detail"
	// 	 where  product_id = 7361 `).Scan(&result)

	// var result string
	// row := db.Raw(`select *  from receive_detail r
	// inner join product p on p.id = r.product_id
	// where p."name"  = 'ABBOTIC 125 MG/  5ML GRANUL 30ML SYR'
	//  `).Row()
	// row.Scan(&recDetail)
	// " (  SELECT max(id) FROM receive_detail where product_id = ? ) ", productId).Scan(&result)

	// var result Result
	// db.Raw("select  price::numeric::integer harga " +
	// 	" from  receive_detail r  " +
	// 	" inner join product p on p.id = r.product_id  " +
	// 	" where p.name  = 'ABBOTIC 125 MG/  5ML GRANUL 30ML SYR'   ").Scan(&result)

	// err := db.Find(&recDetail, "product_id = ?", productId).Error
	// err := db.Where("product_id = ?", productId).Find(&recDetail).Error

	// if err != nil {
	// 	fmt.Println("err " + err.Error())
	// }
	// row := db.Table("receive_detail").Where("product_id = ?", 7361).Select("price").Row()
	// row.Scan(&price)

	// var receiveDetails []dbmodels.ReceiveDetail
	// db.Find(&receiveDetails).Preload("Product", "id = ?", productId)

	// db := GetDbCon()
	// db.Debug().LogMode(true)

	// db.Raw("select "+
	// 	" price::numeric::integer, disc1, p.hpp "+
	// 	" from receive_detail rd "+
	// 	" left join receive r on rd.receive_id = r.id and r.status in(20, 40, 50, 60)  "+
	// 	" inner join product p on rd.product_id = p.id   "+
	// 	" where  rd.product_id = ? "+
	// 	" order by rd.id desc  limit 1 ", productCode).Scan(&res)

	resp := dto.ResultLastPrice2{}
	db.Raw("select  rd.price , p.hpp  "+
		" from product p   "+
		" left  join receive_detail rd  on rd.product_id = p.id "+
		" left join receive r on rd.receive_id = r.id  and r.status in(20, 40, 50, 60) "+
		" where  p.id = ?  order by rd.id desc  limit 1  ", productId).Scan(&resp)

	fmt.Println("isi +", (resp))
	// , " hpp ", resp[0].Hpp)
	fmt.Println("err db", db.GetErrors())
	res.Price = int64(resp.Price)
	// int64(resp.Price)
	return
}
