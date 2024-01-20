package report

import (
	"distribution-system-be/database"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/signintech/gopdf"
)

var (
	returnReceiveNumber string
)

func GenerateReturnReceiveReport(returnSoID int64) {

	title = "Return Receive"

	fmt.Println("Proses RR report ", title)
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
	fmt.Println("data  send to fillData Details : ", returnSoID)
	dataDetails := fillDataDetailReturnReceive(returnSoID)

	fmt.Println("hasil fill")
	for i, ordDetail := range dataDetails {
		fmt.Println(i, "====", ordDetail)
	}
	fmt.Println("=============", title)
	// setFont(&pdf, 12)
	setHeader(&pdf, "rr")
	pdf.Br(20)

	setDetail(&pdf, dataDetails, "rr")
	setSummary(&pdf)
	setSign(&pdf, "Salesman", "", "Apoteker")

	pdf.WritePdf("return-receive.pdf")

}

func fillDataDetailReturnReceive(returnReceiveID int64) []DataDetail {

	returnReceive, err := database.GetReturnReceiveById(returnReceiveID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(returnReceive)

	returnReceiveNumber = returnReceive.ReturnReceiveNo
	returnInvoiceNo = returnReceive.InvoiceNo

	returnDetails := database.GetAllReceiveReturnDetail(returnReceiveID)

	fmt.Println("orderDetails : ", returnDetails)

	go fillDataCustomer(
		returnReceive.Supplier.Code,
		returnReceive.Supplier.Name,
		returnReceive.ReturnReceiveDate.Format("02-01-2006"),
		returnReceive.ReturnReceiveNo,
		returnReceive.Supplier.Alamat,
		returnReceive.Supplier.Kota,
	)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(returnDetails)
	res := make([]DataDetail, totalRec+1)
	var data DataDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, detail := range returnDetails {
		data.Item = detail.Product.Name
		data.Quantity = int64(detail.Qty)
		data.Unit = detail.UOM.Name
		data.Price = int64(detail.Price)
		total := data.Price * data.Quantity
		data.Total = int64(detail.Price) * int64(detail.Qty)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	tax = subTotal / 10
	grandTotal = subTotal + tax

	return res
}
