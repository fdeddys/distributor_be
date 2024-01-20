package database

import (
	"distribution-system-be/constants"
	dbmodels "distribution-system-be/models/dbModels"
	"distribution-system-be/models/dto"
	"distribution-system-be/utils/util"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

//SaveStockOpname ...
func SaveStockOpname(stockOpname *dbmodels.StockOpname) (errCode string, errDesc string, id int64, status int8) {

	fmt.Println(" Update  StockOpname numb ------------------------------------------ ")
	db := GetDbCon()
	db.Debug().LogMode(true)

	r := db.Save(&stockOpname)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		id = 0
		status = 0
		fmt.Println("Error update ", errDesc)
		return
	}

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG, stockOpname.ID, stockOpname.Status
}

// SaveStockOpnameApprove ...
func SaveStockOpnameApprove(stockOpname *dbmodels.StockOpname) (errCode string, errDesc string) {

	fmt.Println(" Approve StockOpname ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// iterate semua product
	// - cek product stock di warehouse
	// - if found di table opname detail -> isi qty system
	// -    - update stock di warehouse sesuai -> qty opname opname detail
	// -    - insert history stock sesuai opname detail -> sesuai qty opname
	// - if not found di table opname detail->  create new baris di opname detail, isi qty on system
	// -    - update stock di warenouse -> 0
	// -    - insert history stock  = 0
	products := ProductList()
	total := float32(0)
	curStock := int64(0)
	for _, product := range products {

		stockFound := true
		stock, errCode, _ := GetStockByProductAndWarehouse(product.ID, stockOpname.WarehouseID)
		if errCode != constants.ERR_CODE_00 {
			stock.LastUpdate = time.Now()
			stock.LastUpdateBy = dto.CurrUser
			stock.ProductID = product.ID
			stock.Qty = 0
			stock.WarehouseID = stockOpname.WarehouseID
			tx.Save(&stock)
			stockFound = false
		} else {
			curStock = stock.Qty
		}

		var stockOpnameDetail dbmodels.StockOpnameDetail

		r := tx.Where(" stock_opname_id = ? and product_id = ? ", stockOpname.ID, product.ID).Find(&stockOpnameDetail).Error

		// opname tidak ditemukan
		if r != nil {
			stockOpnameDetail.Hpp = product.Hpp
			stockOpnameDetail.LastUpdateBy = dto.CurrUser
			stockOpnameDetail.LastUpdate = time.Now()
			stockOpnameDetail.ProductID = product.ID
			stockOpnameDetail.Qty = 0
			stockOpnameDetail.QtyOnSystem = stock.Qty
			stockOpnameDetail.StockOpnameID = stockOpname.ID
			stockOpnameDetail.UomID = product.SmallUomID
			tx.Save(&stockOpnameDetail)

			if stockFound {
				stock.Qty = 0
				stock.WarehouseID = stockOpname.WarehouseID
				stock.LastUpdate = time.Now()
				stock.LastUpdateBy = dto.CurrUser
				tx.Save(&stock)
				total += float32(stockOpnameDetail.Qty-curStock) * product.Hpp
				fmt.Println("1.name = ", product.Name, ":total=", total, ":so qty", stockOpnameDetail.Qty, "- product qty:", curStock, ":hpp-", product.Hpp)

			} else {
				total += float32(stockOpnameDetail.Qty) * product.Hpp
				fmt.Println("2.name = ", product.Name, ":total=", total, ":so qty", stockOpnameDetail.Qty, "- product qty:", curStock, ":hpp-", product.Hpp)

			}

			var history dbmodels.HistoryStock
			history.Code = product.Code
			history.Name = product.Name
			history.Debet = 0
			history.Kredit = 0
			history.Saldo = 0
			history.Description = "STOCK OPNAME"
			history.TransDate = stockOpname.StockOpnameDate
			history.ReffNo = stockOpname.StockOpnameNo
			history.Price = 0
			history.Hpp = product.Hpp
			history.LastUpdateBy = dto.CurrUser
			history.LastUpdate = util.GetCurrDate()
			history.Disc1 = 0
			history.Total = 0
			history.WarehouseID = stockOpname.WarehouseID
			tx.Save(&history)

			continue
		}
		// opame ditemukan
		stockOpnameDetail.LastUpdateBy = dto.CurrUser
		stockOpnameDetail.LastUpdate = time.Now()
		stockOpnameDetail.QtyOnSystem = stock.Qty
		tx.Save(&stockOpnameDetail)

		// if stockFound {
		stock.Qty = stockOpnameDetail.Qty
		stock.WarehouseID = stockOpname.WarehouseID
		tx.Save(&stock)
		// }

		var history dbmodels.HistoryStock
		history.Code = product.Code
		history.Name = product.Name
		history.Debet = stockOpnameDetail.Qty
		history.Kredit = 0
		history.Saldo = stockOpnameDetail.Qty
		history.Description = "STOCK OPNAME"
		history.TransDate = stockOpname.StockOpnameDate
		history.ReffNo = stockOpname.StockOpnameNo
		history.Price = 0
		history.Hpp = stockOpnameDetail.Hpp
		history.LastUpdateBy = dto.CurrUser
		history.LastUpdate = util.GetCurrDate()
		history.Disc1 = 0
		history.Total = 0
		history.WarehouseID = stockOpname.WarehouseID
		tx.Save(&history)

		total += float32(stockOpnameDetail.Qty-curStock) * (stockOpnameDetail.Hpp)
		fmt.Println("3.name = ", product.Name, ":total=", total, ":so qty", stockOpnameDetail.Qty, "- product qty:", curStock, ":hpp-", product.Hpp)

	}

	db.Debug().LogMode(false)

	fmt.Println("Sebelom save total = ", total)
	stockOpname.LastUpdateBy = dto.CurrUser
	stockOpname.LastUpdate = util.GetCurrDate()
	stockOpname.Status = 20
	stockOpname.Total = total
	r := db.Save(&stockOpname)
	if r.Error != nil {
		errCode = constants.ERR_CODE_80
		errDesc = r.Error.Error()
		fmt.Println("Error update stock mutation ", errDesc)
		return
	}

	tx.Commit()
	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}

// GetStockOpnameByStockOpnameNo ...
func GetStockOpnameByStockOpnameNo(stockOpnameNo string) (dbmodels.StockOpname, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	order := dbmodels.StockOpname{}

	err := db.Preload("Warehouse").Where(" stock_opname_no = ?  ", stockOpnameNo).First(&order).Error

	return order, err

}

// GetStockOpnameByStockOpnameId ...
func GetStockOpnameById(stockOpnameID int64) (dbmodels.StockOpname, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)
	stockOpname := dbmodels.StockOpname{}

	err := db.Preload("Warehouse").Where(" id = ?  ", stockOpnameID).First(&stockOpname).Error
	return stockOpname, err

}

// GetStockOpnamePage ...
func GetStockOpnamePage(param dto.FilterStockOpname, offset, limit, internalStatus int) ([]dbmodels.StockOpname, int, error) {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var stockOpnames []dbmodels.StockOpname
	var total int

	var err error

	if offset == 0 && limit == 0 {
		err = db.Find(&stockOpnames).Error
		if err != nil {
			return stockOpnames, 0, err
		}
		return stockOpnames, 0, nil
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	errQuery := make(chan error)
	errCount := make(chan error)

	go AsyncQuerysStockOpnames(db, offset, limit, internalStatus, &stockOpnames, param, errQuery)
	go AsyncQueryCountsStockOpnames(db, &total, internalStatus, &stockOpnames, param, errCount)

	resErrQuery := <-errQuery
	resErrCount := <-errCount

	wg.Done()

	if resErrQuery != nil {
		return stockOpnames, 0, resErrQuery
	}

	if resErrCount != nil {
		log.Println("errr-->", resErrCount)
		return stockOpnames, 0, resErrCount
	}
	return stockOpnames, total, nil
}

func getParamStockOpname(param dto.FilterStockOpname, status int) (stockOpnameNumber string, byStatus bool) {

	stockOpnameNumber = param.StockOpnameNumber
	if stockOpnameNumber == "" {
		stockOpnameNumber = "%"
	} else {
		stockOpnameNumber = "%" + param.StockOpnameNumber + "%"
	}

	byStatus = true
	if status == -1 {
		byStatus = false
	}

	return
}

// AsyncQueryCountsStockOpnames ...
func AsyncQueryCountsStockOpnames(db *gorm.DB, total *int, status int, stockOpname *[]dbmodels.StockOpname, param dto.FilterStockOpname, resChan chan error) {

	stockOpnameNumber, byStatus := getParamStockOpname(param, status)
	fmt.Println("ISI  stock_opname Number ", stockOpnameNumber, "  status ", status, " fill status ", byStatus)

	var err error
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		err = db.Model(&stockOpname).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(stock_opname_no, '') ilike ? AND stock_opname_date between ? and ?  ", status, byStatus, stockOpnameNumber, param.StartDate, param.EndDate).Count(&*total).Error
	} else {
		err = db.Model(&stockOpname).Where(" ( (status = ?) or ( not ?) ) AND COALESCE(stock_opname_no,'') ilike ? ", status, byStatus, stockOpnameNumber).Count(&*total).Error
	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

// AsyncQuerysStockOpnames ...
func AsyncQuerysStockOpnames(db *gorm.DB, offset int, limit int, status int, stockOpnames *[]dbmodels.StockOpname, param dto.FilterStockOpname, resChan chan error) {

	var err error
	stockOpnameNumber, byStatus := getParamStockOpname(param, status)
	fmt.Println("ISI  stock_opname Number ", stockOpnameNumber, "  status ", status, " fill status ", byStatus)

	fmt.Println("isi dari filter [", param, "] ")
	if strings.TrimSpace(param.StartDate) != "" && strings.TrimSpace(param.EndDate) != "" {
		fmt.Println("isi dari filter [", param.StartDate, '-', param.EndDate, "] ")
		err = db.Preload("Warehouse").Order("stock_opname_date DESC").Offset(offset).Limit(limit).Find(&stockOpnames, " ( ( status = ?) or ( not ?) ) AND COALESCE(stock_opname_no, '') ilike ? AND stock_opname_date between ? and ?   ", status, byStatus, stockOpnameNumber, param.StartDate, param.EndDate).Error
	} else {
		fmt.Println("isi dari kosong ")
		err = db.Offset(offset).Limit(limit).Preload("Warehouse").Order("stock_opname_date DESC").Find(&stockOpnames, " ( ( status = ?) or ( not ?) ) AND COALESCE(stock_opname_no,'') ilike ?  ", status, byStatus, stockOpnameNumber).Error
		if err != nil {
			fmt.Println("error --> ", err)
		}
		fmt.Println("order--> ", stockOpnames)

	}

	if err != nil {
		resChan <- err
	}
	resChan <- nil
}

func getTotalQtyByProduct(db *gorm.DB, opnameID, productID int64) int64 {

	var total int64 = 0
	db.Raw("SELECT SUM(qty) FROM stock_opname_detail WHERE product_id = ? and stock_opname_id = ? ", productID, opnameID).Scan(&total)

	return total

}

func FindDataStockByWarehouseID(warehousID int64) []dto.TemplateReportStockOpname {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var datas []dto.TemplateReportStockOpname

	db.Raw("select p.id as product_id, p.name as product_name , s.qty, l.name as uom_name , p.small_uom_id as uom_id "+
		" from product p "+
		" inner join stock s on s.product_id = p.id and s.warehouse_id = ? "+
		" left join lookup l on p.small_uom_id  = l.id  "+
		" where p.status =1 order by p.name ", warehousID).Scan(&datas)

	return datas

}

func RecalculateTotal() (errCode string, errDesc string) {

	fmt.Println(" Recalculate StockOpname ------------------------------------------ ")
	db := GetDbCon()
	tx := db.Begin()
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		tx.Rollback()
	// 	}
	// }()

	var stockOpnames []dbmodels.StockOpname
	err := db.Where("total = 0   ").Find(&stockOpnames).Error
	if err != nil {
		return
	}

	products := ProductList()
	for _, stockOpname := range stockOpnames {
		fmt.Println("stockopname ", stockOpname.StockOpnameNo)

		total := float32(0)
		curStock := int64(0)
		for _, product := range products {

			stockFound := true
			stock, errCode, _ := GetStockByProductAndWarehouse(product.ID, stockOpname.WarehouseID)
			if errCode != constants.ERR_CODE_00 {
				stockFound = false
			} else {
				curStock = stock.Qty
			}

			var stockOpnameDetail dbmodels.StockOpnameDetail

			r := tx.Where(" stock_opname_id = ? and product_id = ? ", stockOpname.ID, product.ID).Find(&stockOpnameDetail).Error

			// opname tidak ditemukan
			if r != nil {

				if stockFound {
					total += float32(stockOpnameDetail.Qty-curStock) * product.Hpp
				} else {
					total += float32(stockOpnameDetail.Qty) * product.Hpp
				}
				fmt.Println("Update detail =", total)
				continue
			}
			if stockOpnameDetail.Qty > 0 || curStock > 0 {
				fmt.Println("stockopname = ", stockOpnameDetail.Qty, " : ", stock.Qty)
			}

			total += float32(stockOpnameDetail.Qty-curStock) * (stockOpnameDetail.Hpp)
			fmt.Println("Update detail =", total)

		}

		db.Debug().LogMode(true)

		stockOpname.Total = total
		r := db.Save(&stockOpname)
		if r.Error != nil {
			errCode = constants.ERR_CODE_80
			errDesc = r.Error.Error()
			fmt.Println("Error update stock  ", errDesc)
			return
		}
	}
	tx.Commit()

	return constants.ERR_CODE_00, constants.ERR_CODE_00_MSG
}
