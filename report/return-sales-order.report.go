package report

import (
	"distribution-system-be/database"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/signintech/gopdf"
)

var (
	returnSONumber  string
	returnInvoiceNo string
)

func GenerateReturnSalesOrderReport(returnSoID int64) {

	title = "Sales Order Return"

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

	spaceCustomerInfo1 = tblCol1
	spaceTitik1 = spaceCustomerInfo1 + 150
	spaceValue1 = spaceCustomerInfo1 + 160

	spaceCustomerInfo = 300
	spaceTitik = spaceCustomerInfo + 150
	spaceValue = spaceCustomerInfo + 160

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
	fmt.Println("data  send to fillData Details : ", returnSoID)
	dataDetails := fillDataDetailReturnSO(returnSoID)

	fmt.Println("hasil fill")
	for i, ordDetail := range dataDetails {
		fmt.Println(i, "====", ordDetail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeader(&pdf, "rso")
	pdf.Br(20)

	setDetail(&pdf, dataDetails)
	setSummary(&pdf)
	setSign(&pdf, "Admin", "Salesman", "Customer")

	pdf.WritePdf("return-so.pdf")

}

func fillDataDetailReturnSO(returnSoID int64) []DataDetail {

	returnSo, err := database.GetSalesOrderReturnById(returnSoID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(returnSo)

	returnSONumber = returnSo.ReturnSalesOrderNo
	returnInvoiceNo = returnSo.InvoiceNo

	returnSoDetails := database.GetAllSalesOrderReturnDetail(returnSoID)

	fmt.Println("orderDetails : ", returnSoDetails)

	go fillDataCustomer(
		returnSo.Customer.Code,
		returnSo.Customer.Name,
		returnSo.ReturnSalesOrderDate.Format("02-01-2006"),
		returnSo.ReturnSalesOrderNo,
		"", "",
	)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(returnSoDetails)
	res := make([]DataDetail, totalRec+1)
	var data DataDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, ordDetail := range returnSoDetails {
		data.Item = ordDetail.Product.Name
		data.Quantity = int64(ordDetail.Qty)
		data.Unit = ordDetail.UOM.Name
		data.Price = int64(ordDetail.Price)
		total := data.Price * data.Quantity
		data.Total = int64(ordDetail.Price) * int64(ordDetail.Qty)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	// tax = subTotal / 10
	if returnSo.Tax > 0 {
		tax = int64(returnSo.Tax)
	}

	grandTotal = subTotal + tax

	return res
}
