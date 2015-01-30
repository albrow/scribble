package util

import (
	"bufio"
	"io/ioutil"
)

const frontMatterDelim = "+++\n"

// SplitFrontMatter reads from r and splits its contents into two pieces:
// frontmatter and other content. If there is no frontmatter,
// then first return value will be an empty string.
func SplitFrontMatter(r *bufio.Reader) (frontMatter string, content string, err error) {
	if contains, err := ContainsFrontMatter(r); err != nil {
		return "", "", err
	} else if contains {
		// split the file into two pieces according to where we
		// find the closing delimiter
		frontMatter := ""
		content := ""
		scanner := bufio.NewScanner(r)
		// skip first line because it's just the delimiter
		scanner.Scan()
		// whether or not we have reached content portion yet
		reachedContent := false
		for scanner.Scan() {
			if line := scanner.Text() + "\n"; line == frontMatterDelim {
				// we have reached the closing delimiter, everything
				// else in the file is content.
				reachedContent = true
			} else {
				if reachedContent {
					content += line
				} else {
					frontMatter += line
				}
			}
		}
		return frontMatter, content, nil
	} else {
		// there is no front matter
		if content, err := ioutil.ReadAll(r); err != nil {
			return "", "", err
		} else {
			return "", string(content), nil
		}
	}
}

// ContainsFrontMatter returns true iff the contents of r include
// frontmatter. It checks the first couple of bytes to see if they
// equal the front matter delimiter without changing the position of
// r.
func ContainsFrontMatter(r *bufio.Reader) (bool, error) {
	firstBytes, err := r.Peek(4)
	if err != nil {
		return false, err
	}
	return string(firstBytes) == frontMatterDelim, nil
}
