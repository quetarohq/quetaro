package quetaro

import (
	sq "github.com/Masterminds/squirrel"
)

func init() {
	sq.StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
