package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	args, err := ParseArguments(os.Args[1:])

	if err != nil {
		log.Fatalf("could not parse arguments: %v", err)
		os.Exit(-1)
	}

	histories, err := LoadHistories(args.Histories)

	if err != nil {
		log.Fatalf("could not load histories: %v", err)
		os.Exit(-1)
	}

	blogs, err := LoadBlogs(args.Blogs)

	if err != nil {
		log.Fatalf("could not load blogs: %v", err)
		os.Exit(-1)
	}

	var articles map[string][]Article = make(map[string][]Article)

	for _, blog := range blogs {
		article, subErr := GetArticles(blog)

		if subErr != nil {
			err = errors.Join(err, fmt.Errorf("could not fetch articles for blog %v: %v", blog.Tag, subErr))
			continue
		}

		history := FindHistory(histories, blog.Tag)
		FilterArticles(&article, history)
		articles[blog.Tag] = article
	}

	if err != nil {
		log.Fatalf("could not load all articles %v", err)
	}

	client := NewNotionClient(args.Url, args.Secret, args.Version, args.Database)

	var multiErrs error

	for tag, blogArticles := range articles {
		history := FindHistory(histories, tag)

		for _, blogArticle := range blogArticles {
			blogErr := client.SaveArticle(&blogArticle)

			if blogErr != nil {
				multiErrs = errors.Join(multiErrs, blogErr)
				continue
			}

			history.Ids = append(history.Ids, blogArticle.Guid)
		}
	}

	if multiErrs != nil {
		log.Fatalf("could not save all articles %v", multiErrs)
	}

	err = SaveHistories(args.Histories, histories)

	if err != nil {
		log.Fatalf("couldn't save histories: %v", err)
	}
}
