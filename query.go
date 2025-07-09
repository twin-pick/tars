package tars

func NewQueryParams(usernames, genres []string) *QueryParams {
	return &QueryParams{
		usernames: usernames,
		genres:    genres,
	}
}
