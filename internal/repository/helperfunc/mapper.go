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

func Int(num int32) pgtype.Int4 {
	return pgtype.Int4 {
		Int32: num,
		Valid: true,
	}
}


