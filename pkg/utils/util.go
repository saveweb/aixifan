package utils

func GetUA() string {
	return "aixifan/" + GetVersion().Version
}

// return f"AcFun-{dougaId}_p{partNum}" if partNum else f"AcFun-{dougaId}_p1"
func ToIdentifier(dougaId string, partNum ...string) string {
	if len(partNum) > 0 {
		return "AcFun-" + dougaId + "_p" + partNum[0]
	}
	return "AcFun-" + dougaId + "_p1"
}
