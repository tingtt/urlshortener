package registry

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"urlshortener/utils/tree"
)

var _ Registry = &registry{}

type registry struct {
	data map[string]string

	persistentDataFilePath string
}

// Save implements Registry.
func (r *registry) Save(path string, redirectTarget string) error {
	r.data[path] = redirectTarget
	return nil
}

// Find implements Registry.
func (r *registry) Find(path string) (redirectTarget string, err error) {
	redirectTarget, exists := r.data[path]
	if !exists {
		return "", ErrRedirectTargetNotFound
	}
	return redirectTarget, nil
}

// FindAll implements Registry.
func (r *registry) FindAll() (map[string]string, error) {
	return r.data, nil
}

// Remove implements Registry.
func (r *registry) Remove(path string) error {
	delete(r.data, path)
	return nil
}

// SavePersistently implements Registry.
func (r *registry) SavePersistently() error {
	err := os.MkdirAll(path.Dir(r.persistentDataFilePath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(r.persistentDataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	type redirect struct {
		from string
		to   string
	}
	var root *tree.Node[redirect]
	cnt := 0
	for from, to := range r.data {
		cnt++
		root = tree.Insert(root, redirect{from, to}, func(new, curr redirect) (isLeft bool) {
			return new.from < curr.from
		})
	}
	var redirects = make([]redirect, 0, cnt)
	tree.InOrderTraversal(root, &redirects)

	for _, redirect := range redirects {
		_, err = file.WriteString(fmt.Sprintf("%s,%s\n", redirect.from, redirect.to))
		if err != nil {
			return err
		}
	}
	return nil
}

// loadToCache implements Registry.
func (r *registry) loadToCache() error {
	file, err := os.Open(r.persistentDataFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	if r.data == nil {
		r.data = map[string]string{}
	}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := strings.Split(scanner.Text(), ",")
		if len(line) < 2 {
			return fmt.Errorf("invaid data (%s:%d)", r.persistentDataFilePath, i)
		}
		r.data[line[0]] = line[1]
	}
	return scanner.Err()
}
