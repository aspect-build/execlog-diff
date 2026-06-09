package sawmill

func stringListFind(list []string, item string, i int) int {
	// i is the expected index (from the other spawn).
	// Check the same location first to try and avoid looping through the list.
	if i >= 0 && len(list) > i && list[i] == item {
		return i
	}

	for i, listItem := range list {
		if item == listItem {
			return i
		}
	}

	return -1
}
