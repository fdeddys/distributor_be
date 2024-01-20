package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// GetAllDataDetailPurchaseOrder ...
func GetAllDataDetailPurchaseOrder(purchaseOrderID int64) []dbmodels.PurchaseOrderDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var purchaseOrderDetails []dbmodels.PurchaseOrderDetail

	db.Preload("Product").Preload("UOM").Preload("PoUOM").Find(&purchaseOrderDetails, " po_id = ? and qty > 0 ", purchaseOrderID)

	return purchaseOrderDetails
}

func GetAllDataDetailPurchaseOrderByPoNo(purchaseOrderNo string) []dbmodels.PurchaseOrderDetail {

	db := GetDbCon()
	db.Debug().LogMode(true)

	var purchaseOrder dbmodels.PurchaseOrder

	var purchaseOrderDetails []dbmodels.PurchaseOrderDetail

	errPO := db.Find(&purchaseOrder, " po_no = ?  ", purchaseOrderNo).Error
	if errPO == gorm.ErrRecordNotFound {
		return purchaseOrderDetails
	}

	err := db.Order("id asc").Find(&purchaseOrderDetails, " po_id = ?  and qty > 0 ", purchaseOrder.ID).Error
	if err == gorm.ErrRecordNotFound {
		return purchaseOrderDetails
	}

	return purchaseOrderDetails
}

// GetPurchaseOrderDetailPage ...
func GetPurchaseOrderDetailPage(param dto.FilterPurchaseOrderDetail, offset, limit int) ([]dbmodels.PurchaseOrderDetail, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var purchaseOrderDetails []dbmodels.PurchaseOrderDetail
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&purchaseOrderDetails).Error
		if err != nil {
			return purchaseOrderDetails, 0, err
		}
		return purchaseOrderDetails, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysPurchaseOrderDetails(db, offset, limit, &purchaseOrderDetails, param.PurchaseOrderID, errQuery)
	go AsyncQueryCountsPurchaseOrderDetails(db, &total, param.PurchaseOrderID, offset, limit, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return purchaseOrderDetails, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return purchaseOrderDetails, 0, resErrCount
	}
	return purchaseOrderDetails, total, nil
}

// AsyncQueryCountsPurchaseOrderDetails ...
func AsyncQueryCountsPurchaseOrderDetails(db *gorm.DB, total *int, purchaseOrderID int64, offset int, limit int, resChan chan error) {

	var err error

	err = db.Model(&dbmodels.PurchaseOrderDetail{}).Where("po_id = ?", purchaseOrderID).Count(total).Error

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysPurchaseOrderDetails ...
func AsyncQuerysPurchaseOrderDetails(db *gorm.DB, offset int, limit int, purchaseOrderDetails *[]dbmodels.PurchaseOrderDetail, purchaseOrderID int64, resChan chan error) {

	var err error

	err = db.Offset(offset).Limit(limit).Order("id desc").Preload("Product.BigUom").Preload("Product.SmallUom").Preload("UOM").Preload("PoUOM").Find(&purchaseOrderDetails, "po_id = ? ", purchaseOrderID).Error
	if err != nil {
		fmt.Println("error --> ", err)
	}

	fmt.Println("order--> ", purchaseOrderDetails)

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

//SavePurchaseOrderDetail ...
func SavePurchaseOrderDetail(purchaseOrderDetail *dbmodels.PurchaseOrderDetail) (errCode string, errDesc string) {

	fmt.Println(" Update PurchaseOrder Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Save(&purchaseOrderDetail); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
		return
	}

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", purchaseOrderDetail.ID)
	return

}

// DeletePurchaseOrderDetailById ...
func DeletePurchaseOrderDetailById(id int64) (errCode string, errDesc string) {

	fmt.Println(" Delete PurchaseOrder Detail  ---- ", id)

	db := GetDbCon()
	db.Debug().LogMode(true)

	if r := db.Where("id = ? ", id).Delete(dbmodels.PurchaseOrderDetail{}); r.Error != nil {
		errCode = "99"
		errDesc = r.Error.Error()
		return
	}

	errCode = "00"
	errDesc = fmt.Sprintf("%v", id)
	return

}

func GetLastPricePurchaseOrderDetail(productId int64) (res dto.ResultLastPrice) {

	db := GetDbCon()
	db.Debug().LogMode(true)

	// db.Raw("select "+
	// 	" price::numeric::integer, disc1, p.hpp "+
	// 	" from receive_detail rd "+
	// 	" left join receive r on rd.receive_id = r.id and r.status in(20, 40, 50, 60)  "+
	// 	" inner join product p on rd.product_id = p.id   "+
	// 	" where  rd.product_id = ? "+
	// 	" order by rd.id desc  limit 1 ", productCode).Scan(&res)

	db.Raw("select  price::numeric::integer, disc1, disc2, p.hpp "+
		" from product p   "+
		" left  join receive_detail rd  on rd.product_id = p.id "+
		" left join receive r on rd.receive_id = r.id  and r.status in(20, 40, 50, 60) "+
		" where  p.id = ?  order by rd.id desc  limit 1  ", productId).Scan(&res)

	return

}
