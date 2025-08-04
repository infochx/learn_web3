package main

import "fmt"

func main() {
	// ==========控制流程 考察：数字操作、条件判断=================================================================
	// 1、只出现一次的数字
	// nums :=[]int{1, 2, 3, 4, 3, 2, 1}
	// re := getSingleNumber(nums)
	re := getSingleNumber([]int{1, 2, 3, 4, 3, 2, 1})
	fmt.Println("只出现一次的数字：", re)

	// 2、判断一个整数是否是回文数
	result := isPalindromeNumber(12321)
	fmt.Println(result)

	// ==========字符串 考察：字符串处理、栈的使用=================================================================
	// 3、有效的括号
	isValidParenteseResult := isValidParentheses("()[]")
	fmt.Println(isValidParenteseResult)

	// 4、最长公共前缀
	longPrefix := getLongestCommonPrefix([]string{"abc", "abcd", "abcde"})
	fmt.Println(longPrefix)

	// 	==========基本值类型================================================================================
	// 5、加一 考察：数组操作、进位处理
	plusOne := plusOne([]int{9, 9, 5})
	fmt.Println("加一 结果：", plusOne)

	// 	==========引用类型：切片================================================================================
	// 6、删除有序数组中的重复项
	fmt.Println("不重复数据长度为：", removeDuplicates([]int{1, 2, 3, 4, 5, 5}))

	// 7、合并区间
	fmt.Println("合并后的区间：", merge([][]int{{1, 2}, {3, 5}, {4, 7}, {8, 9}, {9, 10}, {11, 12}}))

	// ==========基础==========================================================================================
	// 8、两数之和
	fmt.Println("两数下标：", twoSum([]int{1, 2, 3, 6}, 4))

}

func twoSum(nums []int, target int) []int {
	ma := make(map[int]int)
	for index, value := range nums {
		ma[value] = index
		targetValue := target - value
		targetIndex, flag := ma[targetValue]
		// if targetIndex != index && targetValue == nums[targetIndex] {
		if targetIndex != index && flag {
			return []int{index, targetIndex}
		}
	}
	return []int{-1, -1}

}

// 7、合并区间
func merge(nums [][]int) [][]int {
	// 从0开始迭代对比
	for i := 0; i < len(nums); i++ {
		start := i
		for (i+1) < len(nums) && nums[i][1] >= nums[i+1][0] {
			i++
			// 需要合并的数组
			add := [][]int{{nums[start][0], nums[i][1]}}
			// nums[i+1]不在数组范围内时
			if (i + 1) == len(nums) {
				nums = append(nums[:start], add...)
				return nums
			}
			// 数组合并
			nums = append(nums[:start], append(add, nums[i+1:]...)...)
			// 下一次迭代对比从start位置开始
			i = start
		}
	}
	return nums
}

// 6、删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	var lenth int
	for i, j := 0, 1; i < len(nums); i++ {
		// i 记录不重复数字，j记录与i下标不相等的数字下标
		for nums[i] == nums[j] {
			j++
			// 数组已经遍历完
			if j == len(nums) {
				nums = nums[:i+1]
				return i + 1
			}
		}
		nums[i+1] = nums[j]
		// 不重复的数字+1
		lenth = i + 1
	}
	nums = nums[:lenth]
	return lenth
}

// 5、加一 考察：数组操作、进位处理
func plusOne(nums []int) []int {
	for i := len(nums) - 1; i >= 0; i-- {
		// 加1操作
		nums[i]++
		// 大于10的时候取个位数
		if nums[i] >= 10 {
			nums[i] %= 10
			if i == 0 {
				// 数组前加一位
				nums = append([]int{1}, nums...)
			} else {
				// 前一位数加1
				nums[i-1]++
			}
		}
	}
	return nums
}

// 4、最长公共前缀
func getLongestCommonPrefix(strs []string) string {
	// 数组长度为0,返回“”
	if len(strs) == 0 {
		return ""
	}

	// 假设第1个字符串为最长公共子串
	commonPrefix := strs[0]
	for _, str := range strs {
		// 公共子串长度不能超过其它字符串的长度
		if len(commonPrefix) > len(str) {
			commonPrefix = commonPrefix[:len(str)]
		}
		// 公共串长度为0时返回“”
		if len(commonPrefix) == 0 {
			return ""
		}

		// 下标对应的byte不相同时，截取公共串长度并终止当前循环
		for i := 0; i < len(commonPrefix); i++ {
			if commonPrefix[i] != str[i] {
				commonPrefix = commonPrefix[:i]
				break
			}
		}
	}
	return commonPrefix
}

// 3、给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
func isValidParentheses(str string) bool {
	// 如果长度为奇数，则返回false
	if len(str)%2 == 1 {
		return false
	}
	// 定义map,右括号 -> 左括号
	mp := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	// 初始化切片，记录左括号
	stack := []rune{}
	// 遍历字符串的每个字符
	for _, ch := range str {
		// 不是右括号，则添加到切片中
		if mp[ch] == 0 {
			// 将左括号添加到切片中
			stack = append(stack, ch)
		} else {
			// 如果长度为0,或者最新添加的左括号与右括号不匹配时返回false
			if len(stack) == 0 || stack[len(stack)-1] != mp[ch] {
				return false
			}
			// 出栈，移除对应的左括号
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

// 2、判断一个整数是否是回文数
func isPalindromeNumber(orgNum int) string {
	num := orgNum
	reverse := 0
	for num > 0 {
		reverse = reverse*10 + num%10
		num /= 10
	}
	if reverse == orgNum {
		return "是回文数"
	} else {
		return "不是回文数"
	}
}

// 1、只出现一次的数字
func getSingleNumber(nums []int) int {
	resultMap := make(map[int]int)
	for _, value := range nums {
		// resultMap[value]++
		resultMap[value] += 1
	}

	for key, value := range resultMap {
		if value == 1 {
			return key
		}
	}
	return -1
}
