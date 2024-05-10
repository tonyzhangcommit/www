package algorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

/*
	实现常见排序算法
*/

// 生成随机整数切片
func GenSliceByNum(num int) (s []int) {
	s = make([]int, 0, num)
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	if num > 0 {
		for i := 0; i < num; i++ {
			s = append(s, rng.Intn(60))
		}
	}
	return
}

// 冒泡排序--原地排序
func SortByBubbles() []int {
	s := GenSliceByNum(12)
	lenSlice := len(s)
	fmt.Println("origin slice is:", s)
	for i := 0; i < lenSlice-1; i++ {
		for j := i + 1; j < lenSlice; j++ {
			if s[i] > s[j] {
				temp := s[i]
				s[i] = s[j]
				s[j] = temp
			}
		}
	}
	fmt.Println("sorted slice is:", s)
	return s
}

// 选择排序--原地排序
func SortBySelect() []int {
	s := GenSliceByNum(12)
	lenSlice := len(s)
	fmt.Println("origin slice is:", s)
	for i := 0; i < lenSlice-1; i++ {
		minIndex := i
		for j := i + 1; j < lenSlice; j++ {
			if s[j] < s[minIndex] {
				minIndex = j
			}
		}
		if minIndex != i {
			temp := s[i]
			s[i] = s[minIndex]
			s[minIndex] = temp
		}
	}
	fmt.Println("sorted the slice is: ", s)
	return s
}

// 插入排序,核心是单个元素可以认为是有序的
func SortByInsert() []int {
	s := GenSliceByNum(12)
	lens := len(s)
	fmt.Println("origin slice is:", s)
	for i := 1; i < lens; i++ {
		insertindex := i
		temp := s[i]
		for insertindex >= 1 && temp < s[insertindex-1] {
			s[insertindex] = s[insertindex-1]
			insertindex -= 1
		}
		s[insertindex] = temp
	}
	fmt.Println("sorted the slice is: ", s)
	return s
}

// 希尔排序,特殊的插入排序,采用了分治法,按照不同的方法将一个大的数组划分为不同个小数组,然后逐次进行插入排序
// 这里划分小组的方式是使用间隔 gap 进行分隔
// 技巧，这里不用显式分组
func SortByShell() []int {
	s := GenSliceByNum(12)
	lens := len(s)
	fmt.Println("origin slice is:", s)
	gap := lens / 2
	for gap > 0 {
		for i := gap; i < lens; i++ {
			insertindex := i
			temp := s[i]
			for insertindex >= gap && temp < s[insertindex-gap] {
				s[insertindex] = s[insertindex-gap]
				insertindex -= gap
			}
			s[insertindex] = temp
		}
		gap = gap / 2
	}
	fmt.Println("sorted the slice is: ", s)
	return s
}

// 归并排序
// 算法思想；分治，将待排序数组逐个分组，直到为单个元素或者只有两个，然后排序后逐次合并
func merge(left, right []int) []int {
	sortedSlice := make([]int, 0, len(left)+len(right))
	i := 0
	j := 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			sortedSlice = append(sortedSlice, left[i])
			i++
		} else {
			sortedSlice = append(sortedSlice, right[j])
			j++
		}
	}
	sortedSlice = append(sortedSlice, left[i:]...)
	sortedSlice = append(sortedSlice, right[j:]...)
	return sortedSlice
}

func SortByMerge(s []int) []int {
	lens := len(s)
	if lens <= 1 {
		return s
	}
	left := SortByMerge(s[:lens/2])
	right := SortByMerge(s[lens/2:])
	return merge(left, right)
}

// 快速排序
// 递归+分治， 分组原理：选一个基准元素，分为左右两组，然后递归
// 快排属于原地排序
func SortByQuick(s []int, start, end int) {
	i := start
	j := end
	stand := s[(start+end)/2]
	if start >= end {
		return
	}
	for i <= j {
		for s[i] < stand {
			i++
		}
		for s[j] > stand {
			j--
		}
		if i <= j {
			s[i], s[j] = s[j], s[i]
			i++
			j--
		}
	}
	// 这里是因为上一步中的i++ 和 j-- ，循环结束后i 应该比j 大
	SortByQuick(s, start, j)
	SortByQuick(s, i, end)
}

// 计数排序
func SortByCount() {
	s := GenSliceByNum(12)
	fmt.Println(s)
	min, max := s[0], s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
		if v >= max {
			max = v
		}
	}
	// 创建计数数组
	count := make([]int, max-min+1)
	// 统计出现的次数
	for _, v := range s {
		count[v-min]++
	}
	index := 0

	for i, c := range count {
		for c > 0 {
			s[index] = c + i
			index++
			c--
		}
	}
	fmt.Println(s)
}

// 二分查找，有序数组
func BinarySearch(arr []int, target int) int {
	low, high := 0, len(arr)
	for low < high {
		mid := low + (high-low)/2 // 防止数据类型溢出
		midVal := arr[mid]
		if midVal < target {
			low = mid + 1
		} else if midVal > target {
			high = mid - 1
		} else {
			return midVal
		}
	}
	return -1
}

// 两数之和
func FindTwoNumSum(numbers []int, target int) (result [][2]int) {
	// 哈希表保存结果
	mapNumbers := make(map[int]int)
	for _, num := range numbers {
		complete := target - num
		if count, exist := mapNumbers[complete]; exist && count > 0 {
			result = append(result, [2]int{num, complete})
			mapNumbers[complete]--
		}
		mapNumbers[num]++
	}
	return
}

// 三数之和
func threeSum(nums []int) [][]int {
	sort.Ints(nums) // 对数组进行排序
	n := len(nums)
	var result [][]int

	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] { // 跳过重复的元素
			continue
		}
		left, right := i+1, n-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})
				left++
				right--
				// 跳过重复的元素
				for left < right && nums[left] == nums[left-1] {
					left++
				}
				for left < right && nums[right] == nums[right+1] {
					right--
				}
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return result
}

// grebecijohn@gmail.com

// AGP(u$8fN4trCLp
