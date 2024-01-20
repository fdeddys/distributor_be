package report

import (
	"distribution-system-be/database"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/signintech/gopdf"
)

var (
	purchaseOrderNumber string
	apoteker            string
	sia                 string
	sipa                string
)

func getParamValue() {

	// parameter, errCode, _, _ := database.GetParameterByNama(constants.PARAMETER_APOTEKER_NAME)
	// if errCode != constants.ERR_CODE_00 {
	// 	fmt.Println("Parameter not set")
	// 	return
	// }
	apoteker = "???"
	sia = "??"
	sipa = "??"
	parameters, err := database.GetParameter()
	if err != nil {
		return
	}

	for _, parameter := range parameters {
		if parameter.Name == "SIA" {
			sia = parameter.Value
		}
		if parameter.Name == "SIPA" {
			sipa = parameter.Value
			fmt.Println("sipa =", sipa)
		}
		if parameter.Name == "apoteker" {
			apoteker = parameter.Value
		}

	}

}

func init() {
	getParamValue()
}

func GeneratePurchaseOrderReportByPoNo(pono string) error {
	po, err := database.GetPurchaseOrderByPurchaseOrderNo(pono)
	if err != nil {

		return err
	}
	GeneratePurchaseOrderReport(po.ID)
	return nil
}

func GeneratePurchaseOrderReport(purchaseOrderID int64) {

	title = "Purchase Order"

	spaceLen = beego.AppConfig.DefaultFloat("report.space-len", 15)
	pageMargin = beego.AppConfig.DefaultFloat("report.page-margin", 12)

	curPage = 1

	spaceCustomerInfo = 300
	spaceTitik = spaceCustomerInfo + 150
	spaceValue = spaceCustomerInfo + 160

	spaceSummaryInfo = spaceCustomerInfo
	spaceTitikSumamry = spaceTitik
	spaceValueSummary = spaceValue

	tblCol1 = 25
	tblCol2 = 80
	tblCol3 = 300
	tblCol4 = 370
	tblCol5 = 430
	tblCol6 = 500

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.SetMargins(pageMargin, pageMargin, pageMargin, pageMargin)
	pdf.AddPage()

	if err := pdf.AddTTFFont("open-sans", "font/OpenSans-Regular.ttf"); err != nil {
		log.Print(err.Error())
		return
	}

	if err := pdf.AddTTFFont("open-sans-bold", "font/OpenSans-Bold.ttf"); err != nil {
		log.Print(err.Error())
		return
	}

	// untuk nomor urut barang
	number = 1

	// get Data mockup utk display ke grid
	fmt.Println("data  send to fillData Details : ", purchaseOrderID)
	dataDetails := fillDataDetailPurchaseOrder(purchaseOrderID)

	fmt.Println("hasil fill")
	for i, ordDetail := range dataDetails {
		fmt.Println(i, "====", ordDetail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeader(&pdf, "po")
	pdf.Br(20)

	setDetail(&pdf, dataDetails, "po")
	// setSummary(&pdf)
	setSign(&pdf, "", "", "Apoteker")

	pdf.WritePdf("purchase-order.pdf")

}

func fillDataDetailPurchaseOrder(purchaseOrderID int64) []DataDetail {

	purchaseOrder, err := database.GetPurchaseOrderByPurchaseOrderID(purchaseOrderID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Header : ", purchaseOrder)

	purchaseOrderNumber = purchaseOrder.PurchaserNo
	// purchaseOrderNo = purchaseOrder.PurchaserNo

	purchaseOrderDetails := database.GetAllDataDetailPurchaseOrder(purchaseOrderID)

	fmt.Println("Details : ", purchaseOrderDetails)

	fillDataCustomer(
		purchaseOrder.Supplier.Code,
		purchaseOrder.Supplier.Name,
		purchaseOrder.PurchaserDate.Format("02-01-2006"),
		purchaseOrder.PurchaserNo,
		purchaseOrder.Supplier.Alamat,
		purchaseOrder.Supplier.Kota,
	)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(purchaseOrderDetails)
	res := make([]DataDetail, totalRec+1)
	var data DataDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, detail := range purchaseOrderDetails {
		data.Item = detail.Product.Name
		data.Quantity = int64(detail.PoQty)
		// data.Unit = detail.UOM.Name
		data.Unit = detail.PoUOM.Name
		// data.Price = int64(detail.Price)
		data.Price = int64(detail.PoPrice)
		total := data.Price * data.Quantity
		data.Total = int64(detail.Price) * int64(detail.Qty)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	tax = 0
	if purchaseOrder.Tax > 0 {
		tax = subTotal * int64(purchaseOrder.Tax) / 100
	}
	grandTotal = subTotal + tax

	return res
}
