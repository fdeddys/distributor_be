package report

import (
	"fmt"
	"log"

	"github.com/signintech/gopdf"
)

var (
	// length New Line
	spaceLen float64

	// page margin
	pageMargin float64

	// customer region
	spaceCustomerInfo float64
	spaceTitik        float64
	spaceValue        float64

	spaceSummaryInfo  float64
	spaceTitikSumamry float64
	spaceValueSummary float64

	// table
	tblCol1 float64
	tblCol2 float64
	tblCol3 float64
	tblCol4 float64
	tblCol5 float64
	tblCol6 float64

	curPage  int
	number   int
	totalRec int

	subTotal   int64
	tax        int64
	grandTotal int64
)

func showLine(pdf *gopdf.GoPdf) {
	pdf.SetX(200)
	pdf.Line(pdf.MarginLeft(), pdf.GetY(), 575.0, pdf.GetY())
}

func setFont(pdf *gopdf.GoPdf, size int) {
	if err := pdf.SetFont("open-sans", "", size); err != nil {
		log.Print(err.Error())
		return
	}
	// pdf.SetFont("open-sans", "", size)
}

func setFontBold(pdf *gopdf.GoPdf, size int) {
	if err := pdf.SetFont("open-sans-bold", "", size); err != nil {
		log.Print(err.Error())
		return
	}
	// pdf.SetFont("open-sans", "", size)
}

func space(pdf *gopdf.GoPdf) {
	pdf.Br(spaceLen)
}

func setPageNumb(pdf *gopdf.GoPdf, curPage int) {

	setFont(pdf, 10)
	pdf.SetX(595 - pageMargin - 40)
	pdf.SetY(842 - (pageMargin * 2))
	pdf.Text(fmt.Sprintf("Page %v", curPage))

}
