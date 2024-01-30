// filters.go
package main

func FilterArticles(articles *[]Article, history *History) {
	temp := (*articles)[:0]

	for _, article := range *articles {
		if !contains(history.Ids, article.Guid) {
			temp = append(temp, article)
		}
	}

	*articles = temp
}

func FindHistory(histories []*History, tag string) *History {
	for _, history := range histories {
		if history.Tag == tag {
			return history
		}
	}

	return &History{}
}

func contains(ids []string, id string) bool {
	if ids == nil || len(ids) <= 0 {
		return false
	}

	for _, i := range ids {
		if id == i {
			return true
		}
	}

	return false
}
