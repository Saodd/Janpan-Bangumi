package libs

func QuickSortBlog(li []*Blog) {
	quickSortBlog(li, 0, len(li)-1)
}

func quickSortBlog(li []*Blog, lo, hi int) {
	if hi-lo < 5 {
		quickSortBlogSelect(li, lo, hi)
		return
	}
	mid := quickSortBlogPartition(li, lo, hi)
	quickSortBlog(li, lo, mid-1)
	quickSortBlog(li, mid+1, hi)
}

func quickSortBlogPartition(li []*Blog, lo, hi int) (mid int) {
	l, r := lo, hi
	midValue := li[lo].WatchDate
	for l < r {
		for l <= hi {
			if li[l].WatchDate < midValue { // 比较处
				break
			}
			l++
		}
		for r >= lo {
			if li[r].WatchDate >= midValue { // 比较处
				break
			}
			r--
		}
		if l < r {
			li[l], li[r] = li[r], li[l]
		} else {
			break
		}
	}
	li[lo], li[r] = li[r], li[lo]
	return r
}

func quickSortBlogSelect(li []*Blog, lo, hi int) {
	var min int
	for ; lo < hi; lo++ {
		min = lo
		for i := lo + 1; i <= hi; i++ {
			if li[i].WatchDate > li[min].WatchDate { // 比较处
				min = i
			}
		}
		if lo != min {
			li[lo], li[min] = li[min], li[lo]
		}
	}
}
