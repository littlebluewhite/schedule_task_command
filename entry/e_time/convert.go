package e_time

func S2Status(s *string) Status {
	if s == nil {
		return Prepared
	}
	switch *s {
	case "Success":
		return Success
	case "Failure":
		return Failure
	default:
		return Prepared
	}
}
