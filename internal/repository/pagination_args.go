package repository

func appendPaginationArgs(args []any, limit, offset int) []any {
	out := make([]any, len(args), len(args)+2)
	copy(out, args)
	return append(out, limit, offset)
}
