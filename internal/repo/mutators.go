package repo

import "github.com/jackc/pgx/v5/pgtype"

func NullString(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}
