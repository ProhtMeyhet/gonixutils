package cp

import(
	pathtools "path"
)

func PathsAddPrefix(prefix string, to ...string) (result []string) {
	result = make([]string, len(to))
	for _, path := range to {
		result = append(result, PathAddPrefix(prefix, path))
	}; return
}

func PathAddPrefix(prefix, to string) string {
	return pathtools.Join(prefix, pathtools.Base(to))
}
