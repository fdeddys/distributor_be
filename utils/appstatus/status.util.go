package appstatus

func GetStatus(id int8) string {
	statusName := "Unknown"
	switch id {
	case 1:
	case 10:
		statusName = "Outstanding"
		break
	case 20:
		statusName = "Submit"
		break
	case 30:
		statusName = "Cancel"
		break
	case 40:
		statusName = "Receiving"
		break
	case 50:
		statusName = "Paid"
		break
	case 60:
		statusName = "Reject Payment"
		break
	}
	return statusName
}
