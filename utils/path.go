package utils

// FormatPath will standardize the url path.
func FormatPath(path string) string {
	/* Condition validation: Turn empty string into "/" */
	if len(path) == 0 {
		return "/"
	}

	// Try var buffer bytes.Buffer
	var buf []byte
	n := len(path)
	r := 1
	w := 1

	if path[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := (n > 2 && path[n-1] == '/')
	for r < n {
		switch {
		case path[r] == '/': // Empty path element, trailing slash is added after the end
			r++

		case path[r] == '.' && r+1 == n:
			trailing = true
			r++

		case path[r] == '.' && path[r+1] == '/': // . element
			r++

		case path[r] == '.' && path[r+1] == '.' && (r+2 == n || path[r+2] == '/'): // .. element: remove to last /
			r += 2

			if w > 1 { // can backtrack
				w--

				if buf == nil {
					for w > 1 && path[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, path, w, '/')
				w++
			}

			// copy element
			for r < n && path[r] != '/' {
				bufApp(&buf, path, w, path[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, path, w, '/')
		w++
	}

	if buf == nil {
		return path[:w]
	}
	return string(buf[:w])
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
