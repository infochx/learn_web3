// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Task1{
    // 1.创建一个名为Voting的合约，包含以下功能：
    //  - 一个mapping来存储候选人的得票数
    mapping(address => uint256) public map;
    address[] public accounts; // = new uint[](1);

    // - 一个vote函数，允许用户投票给某个候选人
    function vote(address _account) public {
        if (map[_account] == 0) {
            accounts.push(_account);
        }
        map[_account] = map[_account] + 1;
    }

    // - 一个getVotes函数，返回某个候选人的得票数
    function getVotes(address _account) public view returns (uint256) {
        return map[_account];
    }

    // - 一个resetVotes函数，重置所有候选人的得票数
    function resetVotes() public {
        for (uint256 i = 0; i < accounts.length; i++) {
            delete map[accounts[i]];
        }
    }

    // 2. 反转字符串 (Reverse String)
    // - 题目描述：反转一个字符串。输入 "abcde"，输出 "edcba"
    function ReverseString(string memory str) public pure returns (string memory) {
        bytes memory arr = bytes(str);
        for (uint256 i = 0; i < arr.length / 2; i++) {
            uint256 j = arr.length - 1 - i;
            (arr[i], arr[j]) = (arr[j], arr[i]);
        }
        return string(arr);
    }

    // 3. 用 solidity 实现整数转罗马数字
    // - 题目描述在 https://leetcode.cn/problems/integer-to-roman/description/
    function integerToRoman(uint256 num) public pure returns (string memory) {
        uint16[13] memory values = [1000,900,500,400,100,90,50,40,10,9,5,4,1];
        string[13] memory symbols = ["M","CM","D","CD","C","XC","L","XL","X","IX","V","IV","I"];
        bytes memory result = new bytes(0);
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = abi.encodePacked(result, symbols[i]);
            }
        }
        return string(result);
    }

    // 4. 用 solidity 实现罗马数字转数整数
    // - 题目描述在 https://leetcode.cn/problems/roman-to-integer/description/3.
    mapping (bytes1 => uint) public romanToIntMap;
    constructor() {
        romanToIntMap[bytes1("I")] = 1;
        romanToIntMap[bytes1("V")] = 5;
        romanToIntMap[bytes1("X")] = 10;
        romanToIntMap[bytes1("L")] = 50;
        romanToIntMap[bytes1("C")] = 100;
        romanToIntMap[bytes1("D")] = 500;
        romanToIntMap[bytes1("M")] = 1000;
    }
    function romanToInteger(string memory str) public view  returns (uint256 result) {
        bytes memory romanBytes = bytes(str);
        for (uint256 i = 0; i < romanBytes.length; i++) {
            if (i < romanBytes.length - 1 && romanToIntMap[bytes1(romanBytes[i])] < romanToIntMap[bytes1(romanBytes[i + 1])]) {
                result -= romanToIntMap[bytes1(romanBytes[i])];
            } else {
                result += romanToIntMap[bytes1(romanBytes[i])];
            }
        }
        return result;
    }

    // 5. 合并两个有序数组 (Merge Sorted Array)
    // - 题目描述：将两个有序数组合并为一个有序数组。
    function MergeSortedArray(uint256[] calldata nums1, uint256[] calldata nums2) public pure returns (uint256[] memory)
    {
        uint len1 = nums1.length;
        uint len2 = nums2.length;
        uint256[] memory result = new uint256[](len1+len2);
        uint256 i = 0;
        uint256 j = 0;
        for (uint256 k = 0; k < result.length; k++) {
            if (i == len1) {
                // 如果 nums1 已经遍历完，直接将 nums2 的元素添加到结果数组中
                while (j < len2){
                    result[k] = nums2[j];
                    j++;
                    k++;
                }
                return result;
            }
            if (j == len2) {
                // 如果 nums2 已经遍历完，直接将 nums1 的元素添加到结果数组中
                while (i < len1){
                    result[k]=nums1[i];
                    i++;
                    k++;
                }
                return result;
            }
            if (nums1[i] < nums2[j]) {
                result[k] = nums1[i];
                i++;
            } else {
                result[k] = nums2[j];
                j++;
            }
        }
        return result;
    }

    // 6. 二分查找 (Binary Search)
    // - 题目描述：在一个有序数组中查找目标值。
    function binarySearch(int[] calldata nums, int target) public pure returns (uint){
        uint left = 0;
        uint right = nums.length - 1;
        while (left <= right) {
            uint mid = left + (right - left) / 2;
            if (nums[mid] == target) {
                return mid;
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        return 0;
    }
}