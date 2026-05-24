package helperfunc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func Text(str string) pgtype.Text {
	return pgtype.Text{
		String: str,
		Valid:  str != "",
	}
}

func Int(num int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: num,
		Valid: true,
	}
}
