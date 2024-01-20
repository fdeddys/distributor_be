package report

import (
	"distribution-system-be/database"
	dbmodels "distribution-system-be/models/dbModels"
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/leekchan/accounting"
	"github.com/signintech/gopdf"
)

type DataRecvHdr struct {
	SupplierCode string
	SupplierName string
	TransAt      string
	SourceDoc    string
}

type DataRecvDetail struct {
	Item     string
	Quantity int64
	Unit     string
	Price    int64
	Disc1    int64
	Disc2    int64
	Total    int64
	UomQty   int64
}

var (
	receiveNumb string
	dataHdr     DataRecvHdr
)

func GenerateReceiveReport(receiveId int64) {

	var (
	// length New Line
	// spaceLen float64

	// page margin
	// pageMargin float64

	// spaceSupplierInfo float64
	// spaceTitik float64
	// spaceValue float64

	// spaceSummaryInfo  float64
	// spaceTitikSumamry float64
	// spaceValueSummary float64

	// // table
	// tblCol1 float64
	// tblCol2 float64
	// tblCol3 float64
	// tblCol4 float64
	// tblCol5 float64
	// tblCol6 float64

	// curPage int
	// number  int
	// // dataDetails []DataRecvDetail
	// totalRec int

	// count by system
	// subTotal   int64
	// tax        int64
	// grandTotal int64
	)

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
	fmt.Println("data recv send to fillData Details : ", receiveId)
	dataDetails := fillDataRecvDetail(receiveId)

	fmt.Println("hasil fill")
	for i, recvDetail := range dataDetails {
		fmt.Println(i, "====", recvDetail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeaderReceive(&pdf)
	pdf.Br(20)

	setReceiveDetail(&pdf, dataDetails)
	setSummaryReceive(&pdf)
	setSignReceive(&pdf, "Admin", "Warehouse", "Supplier")
	// 595, H: 842
	// pdf.SetFont("open-sans", "", 14)

	// pdf.SetFont("open-sans", "", 10)
	// for i := 2; i <= 83; i++ {
	// 	pdf.SetX(1)
	// 	pdf.SetY(10 * float64(i))
	// 	pdf.Text(fmt.Sprintf("%v", i))
	// }
	pdf.WritePdf("receive.pdf")

}

func fillDataRecvDetail(receiveId int64) []DataRecvDetail {

	receive, err := database.GetReceiveByReceiveID(receiveId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(receive)

	// invoiceNumb = "IVyymm999999"
	receiveNumb = receive.ReceiveNo

	receiveDetails := database.GetAllDataDetailReceive(receiveId)

	fmt.Println("receive Details : ", receiveDetails)

	fillDataSupplier(receive)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(receiveDetails)
	res := make([]DataRecvDetail, totalRec+1)
	var data DataRecvDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, receiveDetail := range receiveDetails {

		if len(receiveDetail.Product.Name) > 40 {
			data.Item = receiveDetail.Product.Name[0:40]
		} else {
			data.Item = receiveDetail.Product.Name
		}
		data.Unit = receiveDetail.UOM.Name
		if receiveDetail.UomID == receiveDetail.Product.BigUomID {
			data.Quantity = int64(receiveDetail.Qty) / int64(receiveDetail.Product.QtyUom)
			data.UomQty = int64(receiveDetail.Product.QtyUom)
		} else {
			data.Quantity = int64(receiveDetail.Qty)
			data.UomQty = 1
		}
		data.Price = int64(receiveDetail.Price)
		data.Disc1 = int64(receiveDetail.Disc1)
		data.Disc2 = int64(receiveDetail.Disc2)
		total := data.Price * data.Quantity * data.UomQty

		data.Total = int64(receiveDetail.Price) * int64(receiveDetail.Qty)
		disc1 := data.Total * int64(receiveDetail.Disc1) / 100
		data.Total -= disc1
		total -= disc1

		disc2 := data.Total * int64(receiveDetail.Disc2) / 100
		data.Total -= disc2
		total -= disc2

		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	tax = 0
	if receive.Tax > 0 {
		tax = subTotal * int64(receive.Tax) / 100
	}
	grandTotal = subTotal + tax

	return res
}

func fillDataSupplier(receive dbmodels.Receive) {
	fmt.Println("isi header : supplier ===>", receive)
	dataHdr.SupplierCode = receive.Supplier.Code
	dataHdr.SupplierName = receive.Supplier.Name
	// dataHdr.TransAt = receive.ReceiveDate.Format("02-01-2006")
	dataHdr.TransAt = receive.ReceiveDate.Add(time.Hour * time.Duration(7+5)).Format("02-01-2006")
	// EST -5
	// total +7 +5

	// var date, _ = time.Parse(time.RFC822, receive.ReceiveDate.String())
	// fmt.Println("date baru", receive.ReceiveDate.Add(time.Hour*time.Duration(7+5)).Format("02-01-2006"))
	// dataHdr.TransAt = receive.ReceiveDate.String()
	// print location and local time
	// location := receive.ReceiveDate.LoadLocation("Asia/Jakarta")
	// fmt.Println("Location : ", location, " Time : ", receive.ReceiveDate.In(location))

	dataHdr.SourceDoc = receive.PoNo

	fillDataCustomer(dataHdr.SupplierCode, dataHdr.SupplierName, dataHdr.TransAt, dataHdr.SourceDoc, receive.Supplier.Alamat, receive.Supplier.Kota)
}

func setHeaderReceive(pdf *gopdf.GoPdf) {

	showLogo(pdf)
	showCompany(pdf)
	space(pdf)
	showLine(pdf)
	showReceiveNo(pdf)

}

// func showLogoRecv(pdf *gopdf.GoPdf) {

// 	imgSize := spaceLen * 5
// 	posX := 20.0
// 	posY := spaceLen

// 	pdf.Image("imgs/logo4.jpg", posX, posY, &gopdf.Rect{W: imgSize + 68, H: imgSize})
// }

func showReceiveNo(pdf *gopdf.GoPdf) {

	pdf.SetY(30)
	pdf.SetX(450)
	setFontBold(pdf, 10)
	pdf.Text("RECEIVE")

	space(pdf)
	setFont(pdf, 12)
	pdf.SetX(450)
	pdf.Text(receiveNumb)
}

// func showCompany(pdf *gopdf.GoPdf) {

// 	line1 := beego.AppConfig.DefaultString("report.line1", "PT. Reksa Transaksi Sukses Makmur")
// 	line2 := beego.AppConfig.DefaultString("report.line2", "Plaza Mutiara Lt 21 Suite 2105")
// 	line3 := beego.AppConfig.DefaultString("report.line3", "Jl. DR. Ide Anak Agung Gde Agung")
// 	line4 := beego.AppConfig.DefaultString("report.line4", "Kav")
// 	line5 := beego.AppConfig.DefaultString("report.line5", "Setiabudi")
// 	line6 := beego.AppConfig.DefaultString("report.line6", "Postal code")

// 	pdf.Br(15)

// 	setFontBold(pdf, 10)
// 	pdf.SetX(200)
// 	pdf.Text(line1)

// 	space(pdf)
// 	setFont(pdf, 10)
// 	pdf.SetX(200)
// 	pdf.Text(line2)

// 	space(pdf)
// 	pdf.SetX(200)
// 	pdf.Text(line3)

// 	space(pdf)
// 	pdf.SetX(200)
// 	pdf.Text(line4)

// 	space(pdf)
// 	pdf.SetX(200)
// 	pdf.Text(line5)

// 	space(pdf)
// 	pdf.SetX(200)
// 	pdf.Text(line6)
// }

// func showLogo(pdf *gopdf.GoPdf) {

// 	imgSize := spaceLen * 5
// 	posX := 20.0
// 	posY := spaceLen

// 	pdf.Image("imgs/logo3.png", posX, posY, &gopdf.Rect{W: imgSize + 68, H: imgSize})
// }

func setReceiveDetail(pdf *gopdf.GoPdf, data []DataRecvDetail) {

	setPageNumb(pdf, curPage)
	pdf.SetX(20)
	pdf.SetY(spaceLen * 8)

	showSupplier(pdf, "Number")

	space(pdf)
	showHeaderTableReceive(pdf)

	fmt.Println("Panjang array ", len(data), "] ")
	fmt.Println("Total rec => set detail => ", totalRec, "] ")
	fmt.Println("start iterate")
	// var dataDetail DataDetail
	if totalRec > 1 {
		for i := 1; i <= 25; i++ {
			fmt.Println("idx ke [", i, "]", data[number])
			space(pdf)
			showDataReceive(pdf, fmt.Sprintf("%v", number), data[number].Item, data[number].Unit, data[number].Quantity, data[number].Price, data[number].Disc1, data[number].Disc2, data[number].Total)
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
		fmt.Println("NEW page")
		curPage++
		pdf.AddPage()
		setHeaderReceive(pdf)
		setReceiveDetail(pdf, data)
	}
}

func setSummaryReceive(pdf *gopdf.GoPdf) {

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

func showHeaderTableReceive(pdf *gopdf.GoPdf) {

	showLine(pdf)
	space(pdf)
	setFontBold(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text("#")

	pdf.SetX(tblCol2 - 25)
	pdf.Text("Item")

	pdf.SetX(tblCol3 - 10)
	pdf.Text("Quantity")

	pdf.SetX(tblCol4 - 10)
	pdf.Text("Unit")

	pdf.SetX(tblCol5 - 40)
	pdf.Text("Price")

	pdf.SetX(tblCol5 + 10)
	pdf.Text("Disc1")

	pdf.SetX(tblCol5 + 45)
	pdf.Text("Disc2")

	pdf.SetX(tblCol6 + 20)
	pdf.Text("Total")

	space(pdf)
	showLine(pdf)
}

func showDataReceive(pdf *gopdf.GoPdf, no, item, unit string, qty, price, disc1, disc2, total int64) {

	ac := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	setFont(pdf, 10)
	pdf.SetX(tblCol1)
	pdf.Text(no)

	pdf.SetX(tblCol2 - 25)
	pdf.Text(item)

	pdf.SetX(tblCol3)
	pdf.Text(fmt.Sprintf("%v", qty))

	pdf.SetX(tblCol4 - 10)
	pdf.Text(unit)

	pdf.SetX(tblCol5 - 40)
	// pdf.Text(fmt.Sprintf("%v", price))
	pdf.Text(ac.FormatMoney(price))

	pdf.SetX(tblCol5 + 15)
	pdf.Text(ac.FormatMoney(disc1))

	pdf.SetX(tblCol5 + 50)
	pdf.Text(ac.FormatMoney(disc2))

	pdf.SetX(tblCol6 + 25)
	// pdf.Text(fmt.Sprintf("%v", total))
	pdf.Text(ac.FormatMoney(total))
}

func showSupplier(pdf *gopdf.GoPdf, TitleNumber string) {
	// , code, name, transDate, ssNo string
	// space(pdf)
	// setFont(pdf, 10)
	setFontBold(pdf, 12)
	pdf.SetX(25)
	// pdf.Text("Supplier ")
	// pdf.SetX(100)
	// pdf.Text(":")
	// pdf.SetX(110)
	fmt.Println("supplier RR", invInfo.CustName)
	pdf.Text(invInfo.CustName)
	setFont(pdf, 10)

	// space(pdf)
	pdf.SetX(spaceCustomerInfo + 70)
	pdf.Text("Date ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.TransAt)

	space(pdf)
	pdf.SetX(25)
	// pdf.Text("Address ")
	// pdf.SetX(100)
	// pdf.Text(":")
	// pdf.SetX(110)
	pdf.Text(invInfo.Address)

	// space(pdf)
	pdf.SetX(spaceCustomerInfo + 70)

	pdf.Text(TitleNumber)
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(invInfo.SourceDoc)

	space(pdf)
	pdf.SetX(25)
	// pdf.Text("City ")
	// pdf.SetX(100)
	// pdf.Text(":")
	// pdf.SetX(110)
	pdf.Text(invInfo.City)

	// space(pdf)
	pdf.SetX(spaceCustomerInfo + 70)
	pdf.Text("SIA")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(sia)
}

func setSignReceive(pdf *gopdf.GoPdf, sign1, sign2, sign3 string) {

	// pdf.SetY(spaceLen * 48)

	xSign1 := tblCol1
	xSign2 := tblCol1 + 200
	xSign3 := tblCol1 + 400
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
		space(pdf)
		pdf.SetX(xSign3)
		pdf.Text(sipa)
	}

}
