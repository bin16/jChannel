package main

type messageSorter struct {
	messages []message
}

func (ms messageSorter) Len() int {
	return len(ms.messages)
}

func (ms messageSorter) Less(a, b int) bool {
	if ms.messages[a].Time.Before(ms.messages[b].Time) {
		return true
	}

	return false
}

func (ms messageSorter) Swap(a, b int) {
	x := ms.messages[a]
	ms.messages[a] = ms.messages[b]
	ms.messages[b] = x
}
