package tempengine

import "strings"

func getTblName(db_tblName string) string {
	return strings.Split(db_tblName, "-")[1]
}
