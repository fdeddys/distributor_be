package database

import (
	"distribution-system-be/models/dto"
)

func ReportPaymentCashByDate(dateStart, dateEnd string) []dto.ReportPaymentCash {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var datas []dto.ReportPaymentCash

	db.Raw("select l.name payment_type_name, p.payment_no , to_char(p.payment_date, 'DD-MON-YYYY') payment_date1, "+
		" so.sales_order_no , to_char(so.order_date , 'DD/Mon/YYYY') order_date, pd.total::integer total_order, so.tax total_ppn , p.total_payment::integer as total_payment , "+
		" to_char(p.last_update, 'DD/Mon/YYYY') last_update, p.last_update_by "+
		" from payment p "+
		" inner join payment_detail pd on pd.payment_id = p.id and p.is_cash = true and p.status = 20 "+
		" left join lookup l on l.id = pd.payment_type_id  "+
		" inner join payment_order po on po.payment_id  = p.id "+
		" inner join sales_order so on po.sales_order_id = so.id  "+
		" where payment_date between ?  and ? "+
		" order by payment_date desc, payment_no asc ", dateStart, dateEnd).Scan(&datas)

	// db.Raw("select l.name payment_type_name, p.payment_no , to_char(p.payment_date, 'DD-MON-YYYY') payment_date1, "+
	// 	" so.sales_order_no , to_char(so.order_date , 'DD/Mon/YYYY') order_date, p.total_order::integer, pd.total::integer as total_payment , "+
	// 	" to_char(p.last_update, 'DD/Mon/YYYY') last_update, p.last_update_by "+
	// 	" from payment p "+
	// 	" inner join payment_detail pd on pd.payment_id = p.id and p.is_cash = true and p.status = 20 "+
	// 	" left join lookup l on l.id = pd.payment_type_id  "+
	// 	" inner join payment_order po on po.payment_id  = p.id "+
	// 	" inner join sales_order so on po.sales_order_id = so.id  "+
	// 	" where payment_date between ?  and ? "+
	// 	" order by payment_date desc, payment_no asc ", dateStart, dateEnd).Scan(&datas)

	return datas

}

func ReportSalesByDate(dateStart, dateEnd string) []dto.ReportSales {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var datas []dto.ReportSales

	db.Raw("select "+
		" to_char(so.order_date , 'DD/Mon/YYYY') as order_date , so.sales_order_no , "+
		" ( case  "+
		"      when so.status = 0 or so.status = 1 or so.status = 10   then 'Outstanding' "+
		"      when  so.status = 20 then 'Submit' "+
		"      when  so.status = 30 then 'Cancel' "+
		"      when  so.status = 40 then 'Receiving' "+
		"      when  so.status = 50 then 'Paid' "+
		"      when  so.status = 60 then 'Reject Payment' "+
		" else so.status::text  end ) as status, "+
		" p.plu, p.name as product_name , sod.qty_order , "+
		" l.name as uom , floor(sod.price) as price , sod.disc1  "+
		" from sales_order so  "+
		" inner join sales_order_detail sod on sod.sales_order_id = so.id "+
		" left join product p on sod.product_id = p.id "+
		" left join lookup l on l.id = sod.uom  "+
		" where so.order_date between ?  and ?  "+
		" and so.status in (20, 50)"+
		" order by so.order_date, so.sales_order_no ", dateStart, dateEnd).Scan(&datas)

	return datas

}

func ReportPaymentSupplierByDate(dateStart, dateEnd string) []dto.ReportPaymentSupplier {
	db := GetDbCon()
	db.Debug().LogMode(true)

	var datas []dto.ReportPaymentSupplier

	// " 	to_char(r.receive_date + interval '1' day, 'DD/Mon/YYYY') as receive_date ,  "+
	db.Raw(
		" select r.receive_no, "+
			" 	to_char(r.receive_date   , 'DD/Mon/YYYY')  as receive_tgl ,  "+
			" 	ps.payment_no ,  "+
			" 	to_char(ps.payment_date , 'DD/Mon/YYYY') as payment_date , s.name as supplier , l.name as payment_type ,ps.payment_reff ,  "+
			" 	( case  when ps.status = 0 or ps.status = 1 or ps.status = 10   then 'Outstanding'  "+
			" 		when  ps.status = 20 then 'Submit' "+
			" 		when  ps.status = 30 then 'Cancel' "+
			" 		when  ps.status = 40 then 'Receiving' "+
			" 		when  ps.status = 50 then 'Paid'  "+
			" 		when  ps.status = 60 then 'Reject Payment'  "+
			" 	else ps.status::text  end ) as status, "+
			"  	r.grand_total  "+
			" from payment_supplier ps  "+
			" inner join payment_supplier_detail psd on psd.payment_supplier_id = ps.id  "+
			" inner join receive r on psd.receiving_id =r.id  "+
			" left join supplier s on ps.supplier_id  = s.id  "+
			" left join lookup l on ps.payment_type_id = l.id  "+
			" where ps.payment_date between ?  and ? "+
			" order by ps.payment_date desc, ps.payment_no asc 	 ", dateStart, dateEnd).Scan(&datas)

	// for _, data := range datas {
	// 	fmt.Println(data.ReceiveNo, "- ", data.ReceiveTgl)
	// }
	return datas

}
