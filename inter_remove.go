package jint

import "strconv"

func DeleteKey(json []byte, path ... string) ([]byte, error){
	if len(path) == 0 {
		return json, NULL_PATH_ERROR()
	}
	offset := 0
	currentPath := path[0]
	chars := []byte{34, 44, 58, 91, 93, 123, 125}
	isJsonChar := make([]bool, 256)
	for _,v := range chars {
		isJsonChar[v] = true
	}
	for space(json[offset]) {
		offset++
	}
	keyStart := -1
	braceType := json[offset]
	for k := 0 ; k < len(path) ; k ++ {
		if braceType == 91 {
			arrayIndex, err := strconv.Atoi(currentPath)
			if err != nil {
				return json, INDEX_EXPECTED_ERROR()
			}
			if arrayIndex == 0 {
				offset++
				for i := offset; i < len(json) ; i ++ {
					curr := json[i]
					if curr == 91 || curr == 123 {
						braceType = curr
						if k != len(path) - 1{
							currentPath = path[k + 1]
						}else{
							return json, KEY_EXPECTED_ERROR()
						}
						offset = i
						break
					}
					if !space(curr){
						break
					}
				}
			}else{
				level := 0
				inQuote := false
				found := false
				indexCount := 0
				isJsonChar[58] = false
				for i := offset ; i < len(json) ; i ++ {
					curr := json[i]
					if !isJsonChar[curr]{
						continue
					}
					if curr == 34 {
						if json[i - 1] == 92 {
							continue
						}
						inQuote = !inQuote
						continue
					}
					if inQuote {
						continue
					}else{
						if curr == 91 || curr == 123{
							if found {
								offset = i
								level++
								braceType = curr
								currentPath = path[k + 1]
								found = false
								break
							}
							level++
							continue
						}
						if curr == 93 || curr == 125 {
							level--
							if level < 1 {
								return json, INDEX_OUT_OF_RANGE_ERROR()
							}
							continue
						}
						if !found {
							if level == 1 {
								if curr == 44 {
									indexCount++
									if indexCount == arrayIndex {
										offset = i + 1
										if k == len(path) - 1{
											return json, KEY_EXPECTED_ERROR()
										}
										found = true
										continue
									}
									continue
								}
								continue
							}
							continue
						}
						continue
					}
				}
				isJsonChar[58] = true
			}
		}else{
			inQuote := false
			found := false
			start := 0
			end := 0
			level := k
			isJsonChar[44] = false
			for i := offset ; i < len(json) ; i ++ {
				curr := json[i]
				if !isJsonChar[curr]{
					continue
				}
				if curr == 34 {
					inQuote = !inQuote
					if found {
						continue
					}
					if level != k + 1 {
						continue
					}
					if inQuote {
						start = i + 1
						continue
					}
					end = i
					continue
				}
				if inQuote {
					continue
				}else{
					if curr == 91 {
						if found {
							braceType = curr
							currentPath = path[k + 1]
							break
						}
						level++
						continue
					}
					if curr == 123 {
						if found {
							k++
							level++
							currentPath = path[k]
							found = false
							continue
						}
						level++
						continue
					}
					if curr == 93 || curr == 125 {
						level--
						continue
					}
					if level == k + 1 {
						if curr == 58 {
							if compare(json, start, end, currentPath) {
								offset = i + 1
								found = true
								if k == len(path) - 1{
									isJsonChar[44] = true
									keyStart = start - 1
									break
								}else{
									continue
								}
							}
							isJsonChar[44] = true
							isJsonChar[58] = false
							for j := i ;  j < len(json) ; j ++ {
								curr := json[j]
								if !isJsonChar[curr]{
									continue
								}
								if curr == 34 {
									if json[j - 1] == 92 {
										continue
									}
									inQuote = !inQuote
									continue
								}
								if inQuote {
									continue
								}else{
									if curr == 91 || curr == 123 {
										level++
										continue
									}
									if curr == 93 || curr == 125 {
										level--
										continue
									}
									if curr == 44 {
										if level == k + 1 {
											i = j
											break
										}
										continue
									}
									continue
								}

							}
							isJsonChar[44] = false
							isJsonChar[58] = true
							continue
						}
						continue
					}
				}
			}
			isJsonChar[44] = true
			if !found {
				return json, KEY_NOT_FOUND_ERROR()
			}
		}
	}
	if offset == 0 {
		return json, BAD_JSON_ERROR()
	}
	if keyStart == -1 {
		return json, BAD_JSON_ERROR()
	}
	for space(json[offset]) {
		offset++
	}
	done := false
	end := -1
	if json[offset] == 123 ||  json[offset] == 91{
		level := 0
		inQuote := false
		for i := offset ; i < len(json) ; i ++ {
			curr := json[i]
			if !isJsonChar[curr]{
				continue
			}
			if curr == 34 {
				if json[i - 1] == 92 {
					continue
				}
				inQuote = !inQuote
				continue
			}
			if inQuote {
				continue
			}else{
				if curr == 91 || curr == 123 {
					level++
				}
				if curr == 93 || curr == 125 {
					level--
					if level == 0 {
						done = true
						end = i
						break
					}
					continue
				}
				continue
			}
			continue
		}
	}else{
		if json[offset] == 34 {
			inQuote := false
			for i := offset ;  i < len(json) ; i ++ {
				curr := json[i]
				if curr == 34 {
					if json[i - 1] == 92 {
						continue
					}
					if inQuote {
						done = true
						end = i
						break
					}
					inQuote = !inQuote
					continue
				}
			}
		}else{
			for i := offset ;  i < len(json) ; i ++ {
				curr := json[i]
				if space(curr) || curr == 44 || curr == 93 || curr == 125 {
					if offset == i {
						return json, EMPTY_ARRAY_ERROR()
					}
					done = true
					end = i
					break
				}
			}
		}
	}
	if done {
		for i := end + 1;  i < len(json) ; i ++ {
			curr := json[i]
			if !space(curr) {
				if curr == 44 {
					return replace(json, []byte{}, keyStart, i + 1), nil
				}
				if curr == 93 || curr == 125 {
					end = i
					break
				}
			}
		}
		for i := 0;  i < len(json) - keyStart ; i ++ {
			curr := json[keyStart - i - 1]
			if !space(curr) {
				if curr == 44 {
					return replace(json, []byte{}, keyStart - i - 1, end), nil
				}
				if curr == 91 || curr == 123 {
					return replace(json, []byte{}, keyStart - i, end), nil
				}
				break
			}
		}
	}
	return json, BAD_JSON_ERROR()
}

func DeleteValue(json []byte, path ... string) ([]byte, error){
	if len(path) == 0 {
		return json, NULL_PATH_ERROR()
	}
	offset := 0
	currentPath := path[0]
	chars := []byte{34, 44, 58, 91, 93, 123, 125}
	isJsonChar := make([]bool, 256)
	for _,v := range chars {
		isJsonChar[v] = true
	}
	for space(json[offset]) {
		offset++
	}
	braceType := json[offset]
	for k := 0 ; k < len(path) ; k ++ {
		if braceType == 91 {
			arrayIndex, err := strconv.Atoi(currentPath)
			if err != nil {
				return json, INDEX_EXPECTED_ERROR()
			}
			if arrayIndex == 0 {
				offset++
				for i := offset; i < len(json) ; i ++ {
					curr := json[i]
					if curr == 123 || curr == 91{
						braceType = curr
						if k != len(path) - 1{
							currentPath = path[k + 1]
						}
						offset = i
						break
					}
					if !space(curr){
						break
					}
				}
			}else{
				level := 0
				inQuote := false
				found := false
				indexCount := 0
				isJsonChar[58] = false
				for i := offset ; i < len(json) ; i ++ {
					curr := json[i]
					if !isJsonChar[curr]{
						continue
					}
					if curr == 34 {
						if json[i - 1] == 92 {
							continue
						}
						inQuote = !inQuote
						continue
					}
					if inQuote {
						continue
					}else{
						if curr == 91 || curr == 123{
							if found {
								offset = i
								level++
								braceType = curr
								currentPath = path[k + 1]
								found = false
								break
							}
							level++
							continue
						}
						if curr == 93 || curr == 125 {
							level--
							if level < 1 {
								return json, INDEX_OUT_OF_RANGE_ERROR()
							}
							continue
						}
						if !found {
							if level == 1 {
								if curr == 44 {
									indexCount++
									if indexCount == arrayIndex {
										offset = i + 1
										if k == len(path) - 1{
											break
										}
										found = true
										continue
									}
									continue
								}
								continue
							}
							continue
						}
						continue
					}
				}
				isJsonChar[58] = true
			}
		}else{
			inQuote := false
			found := false
			start := 0
			end := 0
			level := k
			isJsonChar[44] = false
			for i := offset ; i < len(json) ; i ++ {
				curr := json[i]
				if !isJsonChar[curr]{
					continue
				}
				if curr == 34 {
					inQuote = !inQuote
					if found {
						continue
					}
					if level != k + 1 {
						continue
					}
					if inQuote {
						start = i + 1
						continue
					}
					end = i
					continue
				}
				if inQuote {
					continue
				}else{
					if curr == 91 {
						if found {
							braceType = curr
							currentPath = path[k + 1]
							break
						}
						level++
						continue
					}
					if curr == 123 {
						if found {
							k++
							level++
							currentPath = path[k]
							found = false
							continue
						}
						level++
						continue
					}
					if curr == 93 || curr == 125 {
						level--
						continue
					}
					if level == k + 1 {
						if curr == 58 {
							if compare(json, start, end, currentPath) {
								offset = i + 1
								found = true
								if k == len(path) - 1{
									return json, INDEX_EXPECTED_ERROR()
								}else{
									continue
								}
							}
							isJsonChar[44] = true
							isJsonChar[58] = false
							for j := i ;  j < len(json) ; j ++ {
								curr := json[j]
								if !isJsonChar[curr]{
									continue
								}
								if curr == 34 {
									if json[j - 1] == 92 {
										continue
									}
									inQuote = !inQuote
									continue
								}
								if inQuote {
									continue
								}else{
									if curr == 91 || curr == 123 {
										level++
										continue
									}
									if curr == 93 || curr == 125 {
										level--
										continue
									}
									if curr == 44 {
										if level == k + 1 {
											i = j
											break
										}
										continue
									}
									continue
								}

							}
							isJsonChar[44] = false
							isJsonChar[58] = true
							continue
						}
						continue
					}
				}
			}
			isJsonChar[44] = true
			if !found {
				return json, KEY_NOT_FOUND_ERROR()
			}
		}
	}
	if offset == 0 {
		return json, BAD_JSON_ERROR()
	}
	level := 0
	inQuote := false
	num, err := strconv.Atoi(currentPath)
	if err != nil {
		return json, BAD_JSON_ERROR()
	}
	for i := offset; i < len(json) ; i ++ {
		curr := json[i]
		if !isJsonChar[curr]{
			continue
		}
		if curr == 34 {
			// escape character control
			if json[i - 1] == 92 {
				continue
			}
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}else{
			if curr == 91 || curr == 123 {
				level++
			}
			if curr == 93 || curr == 125 {
				if level == 0 {
					if offset == i {
						return json, EMPTY_ARRAY_ERROR()
					}
					if num == 0 {
						json = replace(json, []byte{}, offset, i)
						return json, nil
					}else{
						json = replace(json, []byte{}, offset - 1, i)
						return json, nil
					}
				}
				level--
				continue
			}
			if curr == 44 {
				if level == 0 {
					if offset == i {
						return json, EMPTY_ARRAY_ERROR()
					}
					json = replace(json, []byte{}, offset, i + 1)
					return json, nil
				}
				continue
			}
			continue
		}
		continue
	}
	return json, BAD_JSON_ERROR()
}