package post

import (
	"errors"
	"sort"
	"sync"
)

var (
	mu      sync.Mutex
	idSeq   int
	storage = make(map[int]Post)
)

type Post struct {
	ID       int
	Title    string
	Body     string
	Comments []Comment
}

func Insert(p *Post) error {
	mu.Lock()
	defer mu.Unlock()
	idSeq++
	p.ID = idSeq
	storage[p.ID] = *p
	return nil
}

func All() ([]Post, error) {
	mu.Lock()
	defer mu.Unlock()

	posts := make([]Post, 0, len(storage))
	for _, v := range storage {
		posts = append(posts, v)
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].ID < posts[j].ID })

	return posts, nil
}

func FindByID(id int) (*Post, error) {
	mu.Lock()
	defer mu.Unlock()
	p, ok := storage[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return &p, nil
}

func Save(p *Post) error {
	mu.Lock()
	defer mu.Unlock()

	_, ok := storage[p.ID]
	if !ok {
		return errors.New("not found")
	}

	storage[p.ID] = *p
	return nil
}

func AddComment(p *Post, c *Comment) error {
	mu.Lock()
	defer mu.Unlock()

	_, ok := storage[p.ID]
	if !ok {
		return errors.New("not found")
	}

	p.Comments = append(p.Comments, *c)
	sp := storage[p.ID]
	sp.Comments = append(sp.Comments, *c)
	storage[p.ID] = sp
	return nil
}
