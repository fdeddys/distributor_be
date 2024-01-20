package report

import (
	"distribution-system-be/database"
	"fmt"
	"log"
	"strings"

	"github.com/astaxie/beego"
	"github.com/leekchan/accounting"
	"github.com/signintech/gopdf"
)

type InvHdrInfo struct {
	CustCode  string
	CustName  string
	TransAt   string
	SourceDoc string
	Address   string
	City      string
}

type DataDetail struct {
	Item     string
	Quantity int64
	Unit     string
	Price    int64
	Total    int64
	disc     int64
}

var (
	// length New Line
	// spaceLen float64

	// page margin
	// pageMargin float64

	// customer region
	// spaceCustomerInfo float64
	// spaceTitik        float64
	// spaceValue        float64

	// spaceSummaryInfo  float64
	// spaceTitikSumamry float64
	// spaceValueSummary float64

	// table
	// tblCol1 float64
	// tblCol2 float64
	// tblCol3 float64
	// tblCol4 float64
	// tblCol5 float64
	// tblCol6 float64

	// curPage     int
	// number      int
	dataDetails []DataDetail
	// totalRec    int
	invoiceNumb string
	invoiceNo   string

	// count by system
	// subTotal   int64
	// tax        int64
	// grandTotal int64

	invInfo InvHdrInfo
	title   string

	spaceCustomerInfo1 float64
	spaceTitik1        float64
	spaceValue1        float64

	spaceSupplierInfo1 float64
)

func GenerateSalesOrderReport(orderId int64, reportType string) {

	switch reportType {
	case "so":
		title = "Sales Order"
	case "invoice":
		title = "Invoice"
	}
	spaceLen = beego.AppConfig.DefaultFloat("report.space-len", 15)
	pageMargin = beego.AppConfig.DefaultFloat("report.page-margin", 12)

	curPage = 1

	tblCol1 = 25
	tblCol2 = 80
	tblCol3 = 300
	tblCol4 = 370
	tblCol5 = 430
	tblCol6 = 500

	// customer baris kiri
	spaceCustomerInfo1 = 50
	spaceTitik1 = spaceCustomerInfo1 + 150
	spaceValue1 = spaceCustomerInfo1 + 160

	// customer baris kanan
	spaceCustomerInfo = 300
	spaceTitik = spaceCustomerInfo + 150
	spaceValue = spaceCustomerInfo + 160

	spaceSummaryInfo = spaceCustomerInfo
	spaceTitikSumamry = spaceTitik
	spaceValueSummary = spaceValue

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
	fmt.Println("data order send to fillData Details : ", orderId)
	dataDetails := fillDataDetail(orderId)

	fmt.Println("hasil fill")
	for i, ordDetail := range dataDetails {
		fmt.Println(i, "====", ordDetail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeader(&pdf, reportType)
	pdf.Br(20)

	setDetail(&pdf, dataDetails, reportType)
	setSummary(&pdf)
	setSign(&pdf, "Admin", "", "Cashier")
	// 595, H: 842
	// pdf.SetFont("open-sans", "", 14)

	// pdf.SetFont("open-sans", "", 10)
	// for i := 2; i <= 83; i++ {
	// 	pdf.SetX(1)
	// 	pdf.SetY(10 * float64(i))
	// 	pdf.Text(fmt.Sprintf("%v", i))
	// }
	pdf.WritePdf("invoice.pdf")

}

func fillDataDetail(orderID int64) []DataDetail {

	order, err := database.GetSalesOrderByOrderId(orderID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order)

	// invoiceNumb = "IVyymm999999"
	invoiceNumb = order.SalesOrderNo
	invoiceNo = order.InvoiceNo

	orderDetails := database.GetAllDataDetail(order.ID)

	fmt.Println("orderDetails : ", orderDetails)

	go fillDataCustomer(
		order.Customer.Code,
		order.Customer.Name,
		order.OrderDate.Format("02-01-2006"),
		order.SalesOrderNo,
		"", "",
	)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(orderDetails)
	res := make([]DataDetail, totalRec+1)
	var data DataDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, ordDetail := range orderDetails {
		data.Item = ordDetail.Product.Name
		data.Quantity = int64(ordDetail.QtyOrder)
		data.Unit = ordDetail.UOM.Name
		data.Price = int64(ordDetail.Price)
		data.disc = int64(ordDetail.Disc1)
		total := (data.Price * data.Quantity) - (data.Price * data.Quantity * data.disc / 100)
		data.Total = total
		// int64(ordDetail.Price) * int64(ordDetail.QtyOrder)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	// tax = subTotal / 10
	if order.Tax > 0 {
		tax = int64(order.Tax * float32(subTotal) / 100)
	}

	grandTotal = subTotal + tax

	return res
}

// func fillDataCustomer(order dbmodels.SalesOrder) {
// 	invInfo.CustCode = order.Customer.Code
// 	invInfo.CustName = order.Customer.Name
// 	invInfo.TransAt = order.OrderDate.Format("02-01-2006")
// 	invInfo.SourceDoc = order.SalesOrderNo
// }

func fillDataCustomer(custCode, custName, transDate, orderNo, address, city string) {
	invInfo.CustCode = custCode
	invInfo.CustName = custName
	invInfo.TransAt = transDate
	invInfo.SourceDoc = orderNo
	invInfo.Address = address
	invInfo.City = city
}

func setHeader(pdf *gopdf.GoPdf, poType string) {

	fmt.Println("set header", title)
	showLogo(pdf)
	showCompany(pdf)
	space(pdf)
	showLine(pdf)
	showInvNo(pdf)

}

func showInvNo(pdf *gopdf.GoPdf) {

	pdf.SetY(30)
	pdf.SetX(450)
	setFontBold(pdf, 10)
	fmt.Println("set inv no", title)
	pdf.Text(title)

	space(pdf)
	setFont(pdf, 12)
	pdf.SetX(450)
	// pdf.Text(invoiceNumb)
	pdf.Text("")

	space(pdf)
	setFont(pdf, 12)
	pdf.SetX(450)
	pdf.Text("")
	// pdf.Text(invoiceNo)

}

func showCompany(pdf *gopdf.GoPdf) {

	line1 := beego.AppConfig.DefaultString("report.line1", "")
	line2 := beego.AppConfig.DefaultString("report.line2", "")
	line3 := beego.AppConfig.DefaultString("report.line3", "")
	line4 := beego.AppConfig.DefaultString("report.line4", "")
	line5 := beego.AppConfig.DefaultString("report.line5", "")
	line6 := beego.AppConfig.DefaultString("report.line6", "")

	pdf.Br(15)

	setFontBold(pdf, 12)
	pdf.SetX(200)
	pdf.SetTextColor(0, 0, 255)
	pdf.Text(line1)

	pdf.SetX(248)
	pdf.SetTextColor(255, 0, 0)
	pdf.Text(line2)

	pdf.SetTextColor(0, 0, 0)
	// space(pdf)
	// pdf.SetX(200)
	// pdf.Text(line2)

	setFont(pdf, 10)
	space(pdf)
	pdf.SetX(200)
	pdf.Text(line3)

	space(pdf)
	pdf.SetX(200)
	pdf.Text(line4)

	space(pdf)
	pdf.SetX(200)
	pdf.Text(line5)

	space(pdf)
	pdf.SetX(200)
	pdf.Text(line6)
}

func showLogo(pdf *gopdf.GoPdf) {

	imgSize := spaceLen * 5
	posX := 60.00
	posY := spaceLen

	pdf.Image("imgs/logo_bg_putih.png", posX, posY, &gopdf.Rect{W: imgSize, H: imgSize})
}

// param [0] = tipe report
// ex: SO, MT
func setDetail(pdf *gopdf.GoPdf, data []DataDetail, param ...string) {
	isPo := false
	isInv := false
	setPageNumb(pdf, curPage)
	pdf.SetX(20)
	pdf.SetY(spaceLen * 8)

	fmt.Println("set detail", param[0])
	if len(param) > 0 {
		switch strings.TrimSpace(param[0]) {
		case "mt":
			showWarehouse(pdf)
		case "so":
			showCustomer(pdf)
		case "po":
			showSupplier(pdf, "PO Number")
			isPo = true
		case "invoice":
			showCustomerInvoice(pdf)
			isInv = true
			// isPo = true
		case "rr":
			fmt.Println("set supplier rr")
			showSupplier(pdf, "Number")
		default:
			showCustomer(pdf)
		}
	} else {
		showCustomer(pdf)
	}

	space(pdf)
	if isPo {
		showHeaderTablePO(pdf)
	} else {
		if isInv {
			showHeaderTableInv(pdf)
		} else {
			showHeaderTable(pdf)
		}
	}

	fmt.Println("Panjang array ", len(data), "] ")
	fmt.Println("Total rec => set detail => ", totalRec, "] ")
	fmt.Println("start iterate")
	// var dataDetail DataDetail
	if totalRec > 1 {
		for i := 1; i <= 25; i++ {
			fmt.Println("idx ke [", i, "]", data[number])
			space(pdf)
			if isPo {
				showDataPO(pdf, fmt.Sprintf("%v", number), data[number].Item, data[number].Unit, data[number].Quantity)
			} else if isInv {
				showDataInv(pdf, fmt.Sprintf("%v", number), data[number].Item, data[number].Unit, data[number].Quantity, data[number].Price, data[number].Total, data[number].disc)
			} else {
				showData(pdf, fmt.Sprintf("%v", number), data[number].Item, data[number].Unit, data[number].Quantity, data[number].Price, data[number].Total)

			}
			number++
			if number >= totalRec {
				break
			}
		}
	}
	// }

	space(pdf)
	showLine(pdf)

	// jika data masih ada utk next page
	// 1. add page
	// 2. set header
	// 3. rekursif
	if totalRec > number {
		setSign(pdf, "", "", "Apoteker")
		fmt.Println("NEW page")
		curPage++
		pdf.AddPage()
		setHeader(pdf, param[0])
		setDetail(pdf, data, param[0])
	}
}

func setSummary(pdf *gopdf.GoPdf) {

	rectangle := gopdf.Rect{}
	rectangle.UnitsToPoints(gopdf.Unit_PT)

	ac := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	setFont(pdf, 10)

	space(pdf)
	// pdf.SetY(spaceLen * 42)

	pdf.SetX(spaceSummaryInfo)
	// pdf.Text("Subtotal")
	pdf.CellWithOption(&rectangle, "Subtotal ", gopdf.CellOption{Align: gopdf.Left, Border: 0, Float: gopdf.Left})
	pdf.SetX(spaceTitikSumamry)
	// pdf.Text(":")
	pdf.CellWithOption(&rectangle, ": ", gopdf.CellOption{Align: gopdf.Center, Border: 0, Float: gopdf.Center})
	// pdf.SetX(spaceValueSummary)
	// pdf.Text(fmt.Sprintf("%v", subTotal))
	// pdf.Text(ac.FormatMoney(subTotal))
	fmt.Println("isi space summ ", spaceValueSummary)
	pdf.SetX(spaceValueSummary + 100)
	pdf.CellWithOption(&rectangle, ac.FormatMoney(subTotal), gopdf.CellOption{Align: gopdf.Right, Border: 0, Float: gopdf.Top})

	space(pdf)
	pdf.SetX(spaceSummaryInfo)
	// pdf.Text("Tax ")
	pdf.CellWithOption(&rectangle, "Tax", gopdf.CellOption{Align: gopdf.Left, Border: 0, Float: gopdf.Left})
	pdf.SetX(spaceTitikSumamry)
	// pdf.Text(":")
	pdf.CellWithOption(&rectangle, ": ", gopdf.CellOption{Align: gopdf.Center, Border: 0, Float: gopdf.Center})
	// pdf.SetX(spaceValueSummary)
	// pdf.Text(fmt.Sprintf("%v", tax))
	// pdf.Text(ac.FormatMoney(tax))
	pdf.SetX(spaceValueSummary + 100)
	pdf.CellWithOption(&rectangle, ac.FormatMoney(tax), gopdf.CellOption{Align: gopdf.Right, Border: 0, Float: gopdf.Top})

	space(pdf)
	pdf.SetX(spaceSummaryInfo)
	// pdf.Text("GrandTotal ")
	pdf.CellWithOption(&rectangle, "GrandTotal", gopdf.CellOption{Align: gopdf.Left, Border: 0, Float: gopdf.Left})

	pdf.SetX(spaceTitikSumamry)
	// pdf.Text(":")
	pdf.CellWithOption(&rectangle, ": ", gopdf.CellOption{Align: gopdf.Center, Border: 0, Float: gopdf.Center})
	// pdf.SetX(spaceValueSummary)
	// // pdf.Text(fmt.Sprintf("%v", grandTotal))
	// pdf.Text(ac.FormatMoney(grandTotal))
	pdf.SetX(spaceValueSummary + 100)
	pdf.CellWithOption(&rectangle, ac.FormatMoney(grandTotal), gopdf.CellOption{Align: gopdf.Right, Border: 0, Float: gopdf.Top})

}

func setSummaryX(pdf *gopdf.GoPdf) {

	setFont(pdf, 10)

	space(pdf)
	// pdf.SetY(spaceLen * 42)
	pdf.SetX(spaceSummaryInfo)
	pdf.Text("Subtotal")
	pdf.SetX(spaceTitikSumamry)
	pdf.Text(":")
	pdf.SetX(spaceValueSummary)
	// fmt.Println("isi subtotal utk summart ", subTotal)
	pdf.Text(fmt.Sprintf("%v", subTotal))

	space(pdf)
	pdf.SetX(spaceSummaryInfo)
	pdf.Text("Tax ")
	pdf.SetX(spaceTitikSumamry)
	pdf.Text(":")
	pdf.SetX(spaceValueSummary)
	pdf.Text(fmt.Sprintf("%v", tax))

	space(pdf)
	pdf.SetX(spaceSummaryInfo)
	pdf.Text("GrandTotal ")
	pdf.SetX(spaceTitikSumamry)
	pdf.Text(":")
	pdf.SetX(spaceValueSummary)
	pdf.Text(fmt.Sprintf("%v", grandTotal))

}

func showHeaderTable(pdf *gopdf.GoPdf) {

	showLine(pdf)
	space(pdf)
	setFontBold(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text("#")

	pdf.SetX(tblCol2)
	pdf.Text("Item")

	pdf.SetX(tblCol3)
	pdf.Text("Quantity")

	pdf.SetX(tblCol4)
	pdf.Text("Unit")

	pdf.SetX(tblCol5)
	pdf.Text("Price")

	pdf.SetX(tblCol6)
	pdf.Text("Total")

	space(pdf)
	showLine(pdf)
}

func showDataInv(pdf *gopdf.GoPdf, no, item, unit string, qty, price, total int64, disc int64) {

	ac := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	setFont(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text(no)

	pdf.SetX(tblCol2)
	pdf.Text(item)

	pdf.SetX(tblCol3)
	pdf.Text(fmt.Sprintf("%v", qty))

	pdf.SetX(tblCol4)
	pdf.Text(unit)

	pdf.SetX(tblCol5)
	// pdf.Text(fmt.Sprintf("%v", price))
	pdf.Text(ac.FormatMoney(price))

	if disc > 0 {
		pdf.SetX(tblCol5 + 45)
		// pdf.Text(fmt.Sprintf("%v", price))
		pdf.Text(fmt.Sprintf("%v", disc))
		// pdf.Text(ac.FormatMoney(disc))
	}

	pdf.SetX(tblCol6 + 35)
	// pdf.Text(fmt.Sprintf("%v", total))
	pdf.Text(ac.FormatMoney(total))
}

func showData(pdf *gopdf.GoPdf, no, item, unit string, qty, price, total int64) {

	ac := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	setFont(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text(no)

	pdf.SetX(tblCol2)
	pdf.Text(item)

	pdf.SetX(tblCol3)
	pdf.Text(fmt.Sprintf("%v", qty))

	pdf.SetX(tblCol4)
	pdf.Text(unit)

	pdf.SetX(tblCol5)
	// pdf.Text(fmt.Sprintf("%v", price))
	pdf.Text(ac.FormatMoney(price))

	pdf.SetX(tblCol6)
	// pdf.Text(fmt.Sprintf("%v", total))
	pdf.Text(ac.FormatMoney(total))
}

func showHeaderTablePO(pdf *gopdf.GoPdf) {

	showLine(pdf)
	space(pdf)
	setFontBold(pdf, 10)
	pdf.SetX(tblCol1 + 20)
	pdf.Text("#")

	pdf.SetX(tblCol2 + 10)
	pdf.Text("Item")

	pdf.SetX(tblCol3 + 100)
	pdf.Text("Quantity")

	pdf.SetX(tblCol4 + 120)
	pdf.Text("Unit")

	space(pdf)
	showLine(pdf)
}
func showDataPO(pdf *gopdf.GoPdf, no, item, unit string, qty int64) {

	setFont(pdf, 10)
	pdf.SetX(tblCol1 + 20)
	pdf.Text(no)

	pdf.SetX(tblCol2 + 10)
	pdf.Text(item)

	pdf.SetX(tblCol3 + 110)
	pdf.Text(fmt.Sprintf("%v", qty))

	pdf.SetX(tblCol4 + 120)
	pdf.Text(unit)

}

func showHeaderTableInv(pdf *gopdf.GoPdf) {

	showLine(pdf)
	space(pdf)
	setFontBold(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text("#")

	pdf.SetX(tblCol2)
	pdf.Text("Item")

	pdf.SetX(tblCol3)
	pdf.Text("Quantity")

	pdf.SetX(tblCol4)
	pdf.Text("Unit")

	pdf.SetX(tblCol5)
	pdf.Text("Price")

	pdf.SetX(tblCol5 + 35)
	pdf.Text("Disc[%]")

	pdf.SetX(tblCol6 + 35)
	pdf.Text("Total")

	space(pdf)
	showLine(pdf)

}

func showDataX(pdf *gopdf.GoPdf, no, item, unit string, qty, price, total int64) {

	setFont(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text(no)

	pdf.SetX(tblCol2)
	pdf.Text(item)

	pdf.SetX(tblCol3)
	pdf.Text(fmt.Sprintf("%v", qty))

	pdf.SetX(tblCol4)
	pdf.Text(unit)

	pdf.SetX(tblCol5)
	pdf.Text(fmt.Sprintf("%v", price))

	pdf.SetX(tblCol6)
	pdf.Text(fmt.Sprintf("%v", total))
}

func showCustomer(pdf *gopdf.GoPdf) {
	// , code, name, transDate, ssNo string
	// space(pdf)
	setFont(pdf, 10)

	pdf.SetX(25)
	pdf.Text("Supplier ")
	pdf.SetX(100)
	pdf.Text(":")
	pdf.SetX(110)
	pdf.Text(invInfo.CustName)

	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Order Date ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.TransAt)

	space(pdf)
	pdf.SetX(spaceCustomerInfo1)
	pdf.Text("")
	pdf.SetX(spaceTitik1)
	pdf.Text("")
	pdf.SetX(spaceValue1)
	pdf.Text("")

	pdf.SetX(spaceCustomerInfo)
	pdf.Text("PO Number ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.SourceDoc)

}

func showCustomerInvoice(pdf *gopdf.GoPdf) {
	// , code, name, transDate, ssNo string
	// space(pdf)
	setFont(pdf, 10)

	pdf.SetX(25)
	// pdf.Text("Supplier ")
	// pdf.SetX(100)
	// pdf.Text(":")
	// pdf.SetX(110)
	pdf.Text(invInfo.CustName)

	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Order Date ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.TransAt)

	space(pdf)
	pdf.SetX(spaceCustomerInfo1)
	pdf.Text("")
	pdf.SetX(spaceTitik1)
	pdf.Text("")
	pdf.SetX(spaceValue1)
	pdf.Text("")

	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Invoice Number ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invoiceNo)

}

func showWarehouse(pdf *gopdf.GoPdf) {
	// , code, name, transDate, ssNo string
	space(pdf)
	setFont(pdf, 10)

	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Warehouse Source")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.CustCode)

	space(pdf)
	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Warehouse Dest ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.CustName)

	space(pdf)
	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Transaction at ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.TransAt)

	space(pdf)
	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Source Document ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.SourceDoc)

}

func setSign(pdf *gopdf.GoPdf, sign1, sign2, sign3 string) {

	// pdf.SetY(spaceLen * 48)

	xSign1 := tblCol1
	xSign2 := tblCol1 + 200
	xSign3 := tblCol1 + 380
	maxLengLine := 100

	xLengSign1 := xSign1 + float64(maxLengLine)
	xLengSign2 := xSign2 + float64(maxLengLine)
	xLengSign3 := xSign3 + float64(maxLengLine)

	space(pdf)
	space(pdf)
	space(pdf)
	space(pdf)

	if sign1 != "" {
		pdf.SetX(xSign1)
		pdf.Text(sign1)
	}

	if sign2 != "" {
		pdf.SetX(xSign2)
		pdf.Text(sign2)
	}

	if sign3 != "" {
		pdf.SetX(xSign3)
		pdf.Text(sign3)
	}

	space(pdf)
	space(pdf)
	space(pdf)
	space(pdf)

	if sign1 != "" {
		pdf.SetX(xSign1)
		pdf.Line(xSign1, pdf.GetY(), xLengSign1, pdf.GetY())
	}

	if sign2 != "" {
		pdf.SetX(xSign2)
		pdf.Line(xSign2, pdf.GetY(), xLengSign2, pdf.GetY())
	}

	if sign3 != "" {
		pdf.SetX(xSign3)
		pdf.Line(xSign3, pdf.GetY(), xLengSign3, pdf.GetY())

		if sign3 == "Apoteker" {
			space(pdf)
			pdf.SetX(xSign3)
			pdf.Text(apoteker)
			space(pdf)
			pdf.SetX(xSign3)
			pdf.Text(sipa)
		}
	}

}
