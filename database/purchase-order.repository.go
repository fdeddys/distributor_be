package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

//SavePurchaseOrder ...
func SavePurchaseOrder(purchaseOrder *dbmodels.PurchaseOrder) (errCode string, errDesc string, id int64, status int8) {

	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&purchaseOrder)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, purchaseOrder.ID, purchaseOrder.Status
}

func ApprovePurchaseOrder(purchaseOrder *dbmodels.PurchaseOrder) (errCode string, errDesc string) {

	fmt.Println(" Reject Purchase Order numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)
	totalPO := CountTotalPO(purchaseOrder.PurchaserNo)
	tax := float32(0)

	if purchaseOrder.IsTax {
		tax = getTaxValue()
		tax = totalPO * (tax / 100)
	}
	r := db.Model(&dbmodels.PurchaseOrder{}).Where("id =?", purchaseOrder.ID).Update(dbmodels.PurchaseOrder{
		Status:        constants.STATUS_APPROVE,
		Total:         totalPO,
		Tax:           tax,
		GrandTotal:    totalPO + tax,
		SupplierID:    purchaseOrder.SupplierID,
		Note:          purchaseOrder.Note,
		PurchaserDate: purchaseOrder.PurchaserDate,
	})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func getTaxValue() float32 {

	tax := float32(10)
	parameteTax, errCode, _, _ := GetParameterByNama(constants.PARAMETER_TAX_VALUE)
	if errCode != constants.ERR_CODE_00 {
		value, err := strconv.ParseFloat(parameteTax.Value, 32)
		if err != nil {
			tax = float32(value)
		}
	}

	return tax
}

func CountTotalPO(poNo string) (total float32) {

	poDetails := GetAllDataDetailPurchaseOrderByPoNo(poNo)
	total = 0
	for _, poDetail := range poDetails {
		subTotal := poDetail.Price * float32(poDetail.Qty)
		subTotal -= (subTotal * poDetail.Disc1 / 100)
		subTotal -= (subTotal * poDetail.Disc2 / 100)
		fmt.Println("total ", poDetail.ID, "   => ", total)
		total += subTotal
	}
	return
}

// GetPurchaseOrderPage ...
func GetPurchaseOrderPage(param dto.FilterPurchaseOrder, offset, limit, internalStatus int) ([]dbmodels.PurchaseOrder, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var purchaseOrders []dbmodels.PurchaseOrder
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&purchaseOrders).Error
		if err != nil {
			return purchaseOrders, 0, err
		}
		return purchaseOrders, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysPurchaseOrders(db, offset, limit, internalStatus, &purchaseOrders, param, errQuery)
	go AsyncQueryCountsPurchaseOrders(db, &total, internalStatus, &purchaseOrders, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return purchaseOrders, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return purchaseOrders, 0, resErrCount
	}
	return purchaseOrders, total, nil
}

// AsyncQueryCountsPurchaseOrders ...
func AsyncQueryCountsPurchaseOrders(db *gorm.DB, total *int, status int, purchaseOrders *[]dbmodels.PurchaseOrder, param dto.FilterPurchaseOrder, resChan chan error) {

	purchaseOrderNumber, byStatus, bySupplierID, suppName := getParamPurchaseOrder(param, status)


	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("1 Rec Number ", purchaseOrderNumber, "  status ", status, " fill status ", byStatus, " supp name : " , suppName)
		if (suppName == "%") {
			err = db.
			Model(&purchaseOrders).
			Joins("left join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Where(" ( (po.status = ?) or ( not ?) ) AND COALESCE(po_no, '') ilike ? AND po_date between ? and ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.StartDate, param.EndDate, param.SupplierId, bySupplierID).
			Count(&*total).
			Error
		} else {
			err = db.
			Model(&purchaseOrders).
			Joins("inner join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Where(" ( (po.status = ?) or ( not ?) ) AND COALESCE(po_no, '') ilike ? AND po_date between ? and ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.StartDate, param.EndDate, param.SupplierId, bySupplierID).
			Count(&*total).
			Error
		}
		
	} else {
		fmt.Println("2 Rec Number ", purchaseOrderNumber, "  status ", status, " fill status ", byStatus, " supp name : " , suppName)
		if (suppName == "%"){
			err = db.
			Model(&purchaseOrders).
			Joins("left join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Where(" ( (po.status = ?) or ( not ?) ) AND COALESCE(po_no,'') ilike ?  AND ( ( supplier_id = ? ) or ( not ?) ) ", status, byStatus, purchaseOrderNumber, param.SupplierId, bySupplierID).
			Count(&*total).
			Error
		} else {
			err = db.
			Model(&purchaseOrders).
			Joins("inner join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Where(" ( (po.status = ?) or ( not ?) ) AND COALESCE(po_no,'') ilike ?  AND ( ( supplier_id = ? ) or ( not ?) ) ", status, byStatus, purchaseOrderNumber, param.SupplierId, bySupplierID).
			Count(&*total).
			Error
		}
		
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysPurchaseOrders ...
func AsyncQuerysPurchaseOrders(db *gorm.DB, offset int, limit int, status int, purchaseOrders *[]dbmodels.PurchaseOrder, param dto.FilterPurchaseOrder, resChan chan error) {

	var err error

	purchaseOrderNumber, byStatus, bySupplierID, suppName := getParamPurchaseOrder(param, status)

	fmt.Println(" PurchaseOrder no ", purchaseOrderNumber, "  status ", status, " fill status ", byStatus, " suppname : " , suppName)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("3 isi dari filter [", param.StartDate, '-', param.EndDate, "] ")

		if (suppName == "%") {
			err = db.
			Joins("left join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Preload("Supplier").
			Order("id DESC").
			Offset(offset).
			Limit(limit).
			Find(&purchaseOrders, " ( ( po.status = ?) or ( not ?) ) AND COALESCE(po_no, '') ilike ? AND po_date between ? and ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.StartDate, param.EndDate, param.SupplierId, bySupplierID).
			Error
		} else{
			err = db.
			Joins("inner join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Preload("Supplier").
			Order("id DESC").
			Offset(offset).
			Limit(limit).
			Find(&purchaseOrders, " ( ( po.status = ?) or ( not ?) ) AND COALESCE(po_no, '') ilike ? AND po_date between ? and ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.StartDate, param.EndDate, param.SupplierId, bySupplierID).
			Error
		}
		
	} else {
		fmt.Println("4 isi dari kosong ")
		if (suppName == "%"){
			err = db.
			Joins("left join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Offset(offset).
			Limit(limit).
			Preload("Supplier").
			Order("id DESC").Find(&purchaseOrders, " ( ( po.status = ?) or ( not ?) ) AND COALESCE(po_no,'') ilike ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.SupplierId, bySupplierID).
			Error
		} else {
			err = db.
			Joins("inner join supplier on supplier.id = po.supplier_id and supplier.name ilike ? ", suppName).
			Offset(offset).
			Limit(limit).
			Preload("Supplier").
			Order("id DESC").Find(&purchaseOrders, " ( ( po.status = ?) or ( not ?) ) AND COALESCE(po_no,'') ilike ?  AND ( ( supplier_id = ? ) or ( not ?) )  ", status, byStatus, purchaseOrderNumber, param.SupplierId, bySupplierID).
			Error
		}
		
		if err != nil {
			fmt.Println("purchaseOrder --> ", err)
		}
		fmt.Println("purchaseOrder--> ", purchaseOrders)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func getParamPurchaseOrder(param dto.FilterPurchaseOrder, status int) (purchaseOrderNumber string, byStatus, bySupplierID bool, supplierName string) {

	purchaseOrderNumber = param.PurchaseOrderNumber
	if purchaseOrderNumber == "" {
		purchaseOrderNumber = "%"
	} else {
		purchaseOrderNumber = "%" + param.PurchaseOrderNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	bySupplierID = true
	if param.SupplierId == 0 {
		bySupplierID = false
	}

	supplierName = param.SupplierName
	// fmt.Println("param ==> ", param.SupplierName)
	if supplierName == "" {
		supplierName = "%"
	} else {
		supplierName = "%" + param.SupplierName + "%"
	}

	// fmt.Println("param ==> ", supplierName)
	// byDate = true
	// if param.StartDate == "" || param.EndDate == "" {
	// 	byDate = false
	// }

	return
}

// GetPurchaseOrderByPurchaseOrderID ...
func GetPurchaseOrderByPurchaseOrderID(purchaseOrderID int64) (dbmodels.PurchaseOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	purchaseOrder := dbmodels.PurchaseOrder{}

	err := db.Preload("Supplier").Where(" id = ?  ", purchaseOrderID).First(&purchaseOrder).Error

	return purchaseOrder, err

}

//RejectPurchaseOrder ...
func RejectPurchaseOrder(purchaseOrder dbmodels.PurchaseOrder) (errCode string, errDesc string) {

	fmt.Println(" Reject PurchaseOrder numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.PurchaseOrder{}).Where("id =?", purchaseOrder.ID).Update(dbmodels.PurchaseOrder{Status: 30})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

//RejectPurchaseOrder ...
func CancelSubmitPurchaseOrder(purchaseOrder dbmodels.PurchaseOrder) (errCode string, errDesc string) {

	fmt.Println(" Reject PurchaseOrder numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.PurchaseOrder{}).Where("id =?", purchaseOrder.ID).Update(dbmodels.PurchaseOrder{Status: 10})
	if r.Error != nil {
		fmt.Println("err reject ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func UpdatePoPaid(poNo string) (errCode string, errDesc string) {

	fmt.Println(" Update PurchaseOrder numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.PurchaseOrder{}).Where("po_no =?", poNo).Update(dbmodels.PurchaseOrder{Status: 40})
	if r.Error != nil {
		fmt.Println("err PO Paid ", r.Error)
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

func UpdateTotal(poID int64, totalPO, tax, grandTotal float32) {

	fmt.Println(" Update Total ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(&dbmodels.PurchaseOrder{}).Where("id =?", poID).Update(dbmodels.PurchaseOrder{
		Total:      totalPO,
		Tax:        tax,
		GrandTotal: grandTotal,
	},
	)

	if r.Error != nil {
		fmt.Println("err PO Paid ", r.Error)
	}

}

// GetPurchaseOrderByPurchaseOrderID ...
func GetPurchaseOrderByPurchaseOrderDetailID(purchaseOrderDetailID int64) (dbmodels.PurchaseOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	purchaseOrder := dbmodels.PurchaseOrder{}

	purchaseDetail := dbmodels.PurchaseOrderDetail{}
	db.Where("id = ? ", purchaseOrderDetailID).First(&purchaseDetail)

	err := db.Preload("Supplier").Where(" id = ?  ", purchaseDetail.PurchaseOrderID).First(&purchaseOrder).Error

	return purchaseOrder, err

}

// GetPurchaseOrderByPurchaseOrderID ...
func GetPurchaseOrderByPurchaseOrderNo(pono string) (dbmodels.PurchaseOrder, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	purchaseOrder := dbmodels.PurchaseOrder{}

	err := db.Preload("Supplier").Where(" po_no = ?  ", pono).First(&purchaseOrder).Error

	return purchaseOrder, err

}

//UpdateReceiveDetail ...
// purchaseOrderDetail.ID, purchaseOrderDetail.Qty, purchaseOrderDetail.Price, purchaseOrderDetail.UomID
func UpdatePODetail(poDetail dbmodels.PurchaseOrderDetail) (errCode string, errDesc string) {

	fmt.Println(" Update Receive Detail  ------------------------------------------ ")

	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Model(dbmodels.PurchaseOrderDetail{}).Where("id = ?", poDetail.ID).Updates(
		dbmodels.PurchaseOrderDetail{
			PoUomID:      poDetail.PoUomID,
			PoQty:        poDetail.PoQty,
			PoPrice:      poDetail.PoPrice,
			PoUOMQty:     poDetail.PoUOMQty,
			Qty:          poDetail.Qty,
			Price:        poDetail.Price,
			UomID:        poDetail.UomID,
			LastUpdate:   util.GetCurrDate(),
			LastUpdateBy: dto.CurrUser,
		})
	if r.Error != nil {
		errCode = constants.ERR_CODE_30
		errDesc = r.Error.Error()
		return
	}

	errCode = constants.ERR_CODE_00
	errDesc = fmt.Sprintf("%v", r.RowsAffected)
	return

}
