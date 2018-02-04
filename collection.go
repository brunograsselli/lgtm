package lgtm

func FilterPullRequest(vs []PullRequest, f func(PullRequest) bool) []PullRequest {
	vsf := make([]PullRequest, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func IndexString(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func IncludeString(vs []string, t string) bool {
	return IndexString(vs, t) >= 0
}

func MapUserString(vs []User, f func(User) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
