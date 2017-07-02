package jira

const (
	MAX_RESULTS = 1000
)

type JqlResult struct {
	MaxResults int
	Issues     []Issue
}

type Issue struct {
}

func GetJqlResult(jql string) {
	startAt := 0
	result := getPart(startAt, MAX_RESULTS)
	for len(result.Issues) == startAt+MAX_RESULTS {
		startAt += MAX_RESULTS
		result.Issues = append(result.Issues, getPart(startAt, MAX_RESULTS).Issues...)
	}
}

func getPart(startAt, maxResults int) JqlResult {
	return JqlResult{}
	//fmt.Sprintf() // TODO get "http://#{@config[:jira][:credentials][:username]}:#{@config[:jira][:credentials][:password]}@#{@config[:jira][:credentials][:hostname]}/rest/api/2/search?#{jql}&maxResults=#{max_results}&startAt=#{start_at}"
}
