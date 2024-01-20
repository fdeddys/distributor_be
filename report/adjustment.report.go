package report

import (
	"distribution-system-be/database"
	dbmodels "distribution-system-be/models/dbModels"
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/leekchan/accounting"
	"github.com/signintech/gopdf"
)

type DataAdjHdr struct {
	TransAt   string
	SourceDoc string
}

type DataAdjDetail struct {
	Item     string
	Quantity int64
	Unit     string
	Price    int64
	Total    int64
}

var (
	adjNumb    string
	dataAdjHdr DataAdjHdr
)

func GenerateSalesAdjustmentReport(adjId int64) {

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
	fmt.Println("data recv send to fillData Details : ", adjId)
	dataDetails := fillDataAdjDetail(adjId)

	fmt.Println("hasil fill")
	for i, recvDetail := range dataDetails {
		fmt.Println(i, "====", recvDetail)
	}
	fmt.Println("=============")
	// setFont(&pdf, 12)
	setHeaderAdjust(&pdf)
	pdf.Br(20)

	setAdjDetail(&pdf, dataDetails)
	setSummaryAdj(&pdf)
	setSignAdj(&pdf, "Admin", "Warehouse", "")
	// 595, H: 842
	// pdf.SetFont("open-sans", "", 14)

	// pdf.SetFont("open-sans", "", 10)
	// for i := 2; i <= 83; i++ {
	// 	pdf.SetX(1)
	// 	pdf.SetY(10 * float64(i))
	// 	pdf.Text(fmt.Sprintf("%v", i))
	// }
	pdf.WritePdf("adjustment.pdf")

}

func fillDataAdjDetail(adjId int64) []DataAdjDetail {

	adj, err := database.GetAdjustmentByAdjustmentID(adjId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(adj)

	// invoiceNumb = "IVyymm999999"
	adjNumb = adj.AdjustmentNo

	adjDetails := database.GetAllDataDetailAdjustment(adjId)

	fmt.Println("receive Details : ", adjDetails)

	go fillDataAdj(adj)
	// tdk blh kosong
	// per halaman max 25 item detail
	totalRec = len(adjDetails)
	res := make([]DataAdjDetail, totalRec+1)
	var data DataAdjDetail

	subTotal = 0
	tax = 0
	grandTotal = 0
	for i, adjDetail := range adjDetails {
		data.Item = adjDetail.Product.Name
		data.Quantity = int64(adjDetail.Qty)
		data.Unit = adjDetail.UOM.Name
		data.Price = int64(adjDetail.Product.Hpp)
		total := data.Price * data.Quantity
		data.Total = int64(adjDetail.Product.Hpp) * int64(adjDetail.Qty)
		subTotal += total
		res[i+1] = data
		fmt.Println("total sub total", subTotal)
	}
	totalRec = len(res)
	fmt.Println("Jumlah record [fill] =>", totalRec)

	tax = 0
	// if receive.Tax > 0 {
	// 	tax = int64(receive.Tax)
	// }
	grandTotal = subTotal + tax

	return res
}

func fillDataAdj(adj dbmodels.Adjustment) {
	// dataHdr.SupplierCode = receive.Supplier.Code
	// dataHdr.SupplierName = receive.Supplier.Name
	dataAdjHdr.TransAt = adj.AdjustmentDate.Format("02-01-2006")
	dataAdjHdr.SourceDoc = adj.AdjustmentNo
}

func setHeaderAdjust(pdf *gopdf.GoPdf) {

	showLogo(pdf)
	showCompany(pdf)
	space(pdf)
	showLine(pdf)
	showAdjNo(pdf)

}

func showAdjNo(pdf *gopdf.GoPdf) {

	pdf.SetY(30)
	pdf.SetX(450)
	setFontBold(pdf, 10)
	pdf.Text("ADJUSTMENT")

	space(pdf)
	setFont(pdf, 12)
	pdf.SetX(450)
	pdf.Text(adjNumb)
}

func setAdjDetail(pdf *gopdf.GoPdf, data []DataAdjDetail) {

	setPageNumb(pdf, curPage)
	pdf.SetX(20)
	pdf.SetY(spaceLen * 8)

	showSupplierAdj(pdf)

	space(pdf)
	showHeaderTableAdj(pdf)

	fmt.Println("Panjang array ", len(data), "] ")
	fmt.Println("Total rec => set detail => ", totalRec, "] ")
	fmt.Println("start iterate")
	// var dataDetail DataDetail
	if totalRec > 1 {
		for i := 1; i <= 25; i++ {
			fmt.Println("idx ke [", i, "]", data[number])
			space(pdf)
			showDataAdj(pdf, fmt.Sprintf("%v", number), data[number].Item, data[number].Unit, data[number].Quantity, data[number].Price, data[number].Total)
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
		setHeaderAdjust(pdf)
		setAdjDetail(pdf, data)
	}
}

func setSummaryAdj(pdf *gopdf.GoPdf) {

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

func showHeaderTableAdj(pdf *gopdf.GoPdf) {

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

func showDataAdj(pdf *gopdf.GoPdf, no, item, unit string, qty, price, total int64) {

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

func showSupplierAdj(pdf *gopdf.GoPdf) {
	// , code, name, transDate, ssNo string
	// space(pdf)
	setFont(pdf, 10)

	// pdf.SetX(spaceCustomerInfo1)
	// pdf.Text("Supplier Code")
	// pdf.SetX(spaceTitik1)
	// pdf.Text(":")
	// pdf.SetX(spaceValue1)
	// pdf.Text(dataHdr.SupplierCode)

	// // space(pdf)
	// pdf.SetX(spaceCustomerInfo)
	// pdf.Text("Supplier ")
	// pdf.SetX(spaceTitik)
	// pdf.Text(":")
	// pdf.SetX(spaceValue)
	// pdf.Text(dataHdr.SupplierName)

	// space(pdf)
	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Transaction at ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(dataAdjHdr.TransAt)

	space(pdf)
	pdf.SetX(spaceCustomerInfo)
	pdf.Text("Source Document ")
	pdf.SetX(spaceTitik)
	pdf.Text(":")
	pdf.SetX(spaceValue)
	pdf.Text(dataAdjHdr.SourceDoc)

}

func setSignAdj(pdf *gopdf.GoPdf, sign1, sign2, sign3 string) {

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
	}

}
