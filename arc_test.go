package arc

import (
	"container/list"
	"reflect"
	"testing"
)

func TestTable(t *testing.T) {
	type result struct {
		k      int
		rl, rd []int
		fl, fd []int
	}
	table := []result{
		{k: 5, rl: []int{5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 17, rl: []int{17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 8, rl: []int{8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 0, rl: []int{0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 13, rl: []int{13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 25, rl: []int{25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 27, rl: []int{27, 25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 26, rl: []int{26, 27, 25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}},
		{k: 8, rl: []int{26, 27, 25, 13, 0, 17, 5}, rd: []int{}, fl: []int{8}, fd: []int{}},
		{k: 13, rl: []int{26, 27, 25, 0, 17, 5}, rd: []int{}, fl: []int{13, 8}, fd: []int{}},
		{k: 19, rl: []int{19, 26, 27, 25, 0, 17, 5}, rd: []int{}, fl: []int{13, 8}, fd: []int{}},
		{k: 17, rl: []int{19, 26, 27, 25, 0, 5}, rd: []int{}, fl: []int{17, 13, 8}, fd: []int{}},
		{k: 15, rl: []int{15, 19, 26, 27, 25, 0, 5}, rd: []int{}, fl: []int{17, 13, 8}, fd: []int{}},
		{k: 28, rl: []int{28, 15, 19, 26, 27, 25, 0}, rd: []int{5}, fl: []int{17, 13, 8}, fd: []int{}},
		{k: 2, rl: []int{2, 28, 15, 19, 26, 27, 25}, rd: []int{0, 5}, fl: []int{17, 13, 8}, fd: []int{}},
		{k: 3, rl: []int{3, 2, 28, 15, 19, 26, 27}, rd: []int{25, 0, 5}, fl: []int{17, 13, 8}, fd: []int{}},
		{k: 2, rl: []int{3, 28, 15, 19, 26, 27}, rd: []int{25, 0, 5}, fl: []int{2, 17, 13, 8}, fd: []int{}},
		{k: 11, rl: []int{11, 3, 28, 15, 19, 26}, rd: []int{27, 25, 0, 5}, fl: []int{2, 17, 13, 8}, fd: []int{}},
		{k: 5, rl: []int{11, 3, 28, 15, 19}, rd: []int{26, 27, 25, 0}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}},
		{k: 14, rl: []int{14, 11, 3, 28, 15}, rd: []int{19, 26, 27, 25, 0}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}},
		{k: 22, rl: []int{22, 14, 11, 3, 28}, rd: []int{15, 19, 26, 27, 25}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}},
		{k: 24, rl: []int{24, 22, 14, 11, 3}, rd: []int{28, 15, 19, 26, 27}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}},
		{k: 4, rl: []int{4, 24, 22, 14, 11}, rd: []int{3, 28, 15, 19, 26}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}},
		{k: 15, rl: []int{4, 24, 22, 14}, rd: []int{11, 3, 28, 19, 26}, fl: []int{15, 5, 2, 17, 13, 8}, fd: []int{}},
		{k: 4, rl: []int{24, 22, 14}, rd: []int{11, 3, 28, 19, 26}, fl: []int{4, 15, 5, 2, 17, 13, 8}, fd: []int{}},
		{k: 11, rl: []int{24, 22, 14}, rd: []int{3, 28, 19, 26}, fl: []int{11, 4, 15, 5, 2, 17, 13}, fd: []int{8}},
		{k: 2, rl: []int{24, 22, 14}, rd: []int{3, 28, 19, 26}, fl: []int{2, 11, 4, 15, 5, 17, 13}, fd: []int{8}},
		{k: 19, rl: []int{24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15, 5, 17}, fd: []int{13, 8}},
		{k: 7, rl: []int{7, 24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15, 5}, fd: []int{17, 13, 8}},
		{k: 29, rl: []int{29, 7, 24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 12, rl: []int{12, 29, 7, 24, 22}, rd: []int{14, 3, 28, 26}, fl: []int{19, 2, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 2, rl: []int{12, 29, 7, 24, 22}, rd: []int{14, 3, 28, 26}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 25, rl: []int{25, 12, 29, 7, 24}, rd: []int{22, 14, 3, 28, 26}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 0, rl: []int{0, 25, 12, 29, 7}, rd: []int{24, 22, 14, 3, 28}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 29, rl: []int{0, 25, 12, 7}, rd: []int{24, 22, 14, 3, 28}, fl: []int{29, 2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}},
		{k: 28, rl: []int{0, 25, 12, 7}, rd: []int{24, 22, 14, 3}, fl: []int{28, 29, 2, 19, 11, 4}, fd: []int{15, 5, 17, 13, 8}},
		{k: 24, rl: []int{0, 25, 12, 7}, rd: []int{22, 14, 3}, fl: []int{24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}},
		{k: 7, rl: []int{0, 25, 12}, rd: []int{22, 14, 3}, fl: []int{7, 24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}},
		{k: 12, rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{12, 7, 24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}},
		{k: 17, rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{17, 12, 7, 24, 28, 29, 2, 19}, fd: []int{11, 4, 15, 5, 13, 8}},
		{k: 5, rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{5, 17, 12, 7, 24, 28, 29, 2}, fd: []int{19, 11, 4, 15, 13, 8}},
		{k: 7, rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{7, 5, 17, 12, 24, 28, 29, 2}, fd: []int{19, 11, 4, 15, 13, 8}},
		{k: 14, rl: []int{0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24, 28, 29}, fd: []int{2, 19, 11, 4, 15, 13, 8}},
		{k: 20, rl: []int{20, 0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24, 28}, fd: []int{29, 2, 19, 11, 4, 15, 13, 8}},
		{k: 26, rl: []int{26, 20, 0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24}, fd: []int{28, 29, 2, 19, 11, 4, 15, 13}},
		{k: 22, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{22, 14, 7, 5, 17, 12}, fd: []int{24, 28, 29, 2, 19, 11, 4, 15, 13}},
		{k: 24, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{24, 22, 14, 7, 5, 17}, fd: []int{12, 28, 29, 2, 19, 11, 4, 15, 13}},
		{k: 12, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{12, 24, 22, 14, 7, 5}, fd: []int{17, 28, 29, 2, 19, 11, 4, 15, 13}},
		{k: 9, rl: []int{9, 26, 20, 0, 25}, rd: []int{3}, fl: []int{12, 24, 22, 14, 7}, fd: []int{5, 17, 28, 29, 2, 19, 11, 4, 15}},
		{k: 9, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{9, 12, 24, 22, 14, 7}, fd: []int{5, 17, 28, 29, 2, 19, 11, 4, 15}},
		{k: 28, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{28, 9, 12, 24, 22, 14}, fd: []int{7, 5, 17, 29, 2, 19, 11, 4, 15}},
		{k: 22, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{22, 28, 9, 12, 24, 14}, fd: []int{7, 5, 17, 29, 2, 19, 11, 4, 15}},
		{k: 29, rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{29, 22, 28, 9, 12, 24}, fd: []int{14, 7, 5, 17, 2, 19, 11, 4, 15}},
		{k: 21, rl: []int{21, 26, 20, 0, 25}, rd: []int{3}, fl: []int{29, 22, 28, 9, 12}, fd: []int{24, 14, 7, 5, 17, 2, 19, 11, 4}},
		{k: 2, rl: []int{21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28, 9, 12}, fd: []int{24, 14, 7, 5, 17, 19, 11, 4}},
		{k: 6, rl: []int{6, 21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28, 9}, fd: []int{12, 24, 14, 7, 5, 17, 19, 11}},
		{k: 18, rl: []int{18, 6, 21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}},
		{k: 21, rl: []int{18, 6, 26, 20, 0}, rd: []int{25, 3}, fl: []int{21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}},
		{k: 6, rl: []int{18, 26, 20, 0}, rd: []int{25, 3}, fl: []int{6, 21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}},
		{k: 17, rl: []int{18, 26, 20}, rd: []int{0, 25, 3}, fl: []int{17, 6, 21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 19}},
		{k: 4, rl: []int{4, 18, 26, 20}, rd: []int{0, 25, 3}, fl: []int{17, 6, 21, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}},
		{k: 18, rl: []int{4, 26, 20}, rd: []int{0, 25, 3}, fl: []int{18, 17, 6, 21, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}},
		{k: 21, rl: []int{4, 26, 20}, rd: []int{0, 25, 3}, fl: []int{21, 18, 17, 6, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}},
		{k: 26, rl: []int{4, 20}, rd: []int{0, 25, 3}, fl: []int{26, 21, 18, 17, 6, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}},
		{k: 15, rl: []int{15, 4, 20}, rd: []int{0, 25, 3}, fl: []int{26, 21, 18, 17, 6, 2, 29}, fd: []int{22, 28, 9, 12, 24, 14, 7}},
		{k: 25, rl: []int{15, 4, 20}, rd: []int{0, 3}, fl: []int{25, 26, 21, 18, 17, 6, 2}, fd: []int{29, 22, 28, 9, 12, 24, 14, 7}},
		{k: 26, rl: []int{15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17, 6, 2}, fd: []int{29, 22, 28, 9, 12, 24, 14, 7}},
		{k: 8, rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17, 6}, fd: []int{2, 29, 22, 28, 9, 12, 24, 14}},
		{k: 16, rl: []int{16, 8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}},
		{k: 25, rl: []int{16, 8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{25, 26, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}},
		{k: 16, rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{16, 25, 26, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}},
		{k: 17, rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{17, 16, 25, 26, 21, 18}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}},
		{k: 0, rl: []int{8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25, 26, 21}, fd: []int{18, 6, 2, 29, 22, 28, 9, 12, 24}},
		{k: 27, rl: []int{27, 8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25, 26}, fd: []int{21, 18, 6, 2, 29, 22, 28, 9, 12}},
		{k: 10, rl: []int{10, 27, 8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25}, fd: []int{26, 21, 18, 6, 2, 29, 22, 28, 9}},
		{k: 4, rl: []int{10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17, 16, 25}, fd: []int{26, 21, 18, 6, 2, 29, 22, 28, 9}},
		{k: 11, rl: []int{11, 10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17, 16}, fd: []int{25, 26, 21, 18, 6, 2, 29, 22, 28}},
		{k: 14, rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17}, fd: []int{16, 25, 26, 21, 18, 6, 2, 29, 22}},
		{k: 3, rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{3, 4, 0}, fd: []int{17, 16, 25, 26, 21, 18, 6, 2, 29, 22}},
		{k: 26, rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26, 3, 4}, fd: []int{0, 17, 16, 25, 21, 18, 6, 2, 29, 22}},
		{k: 9, rl: []int{9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26, 3}, fd: []int{4, 0, 17, 16, 25, 21, 18, 6, 2, 29}},
		{k: 23, rl: []int{23, 9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26}, fd: []int{3, 4, 0, 17, 16, 25, 21, 18, 6, 2}},
		{k: 13, rl: []int{13, 23, 9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{}, fd: []int{26, 3, 4, 0, 17, 16, 25, 21, 18, 6}},
		{k: 19, rl: []int{19, 13, 23, 9, 14, 11, 10, 27, 8, 15}, rd: []int{}, fl: []int{}, fd: []int{26, 3, 4, 0, 17, 16, 25, 21, 18, 6}},
		{k: 26, rl: []int{19, 13, 23, 9, 14, 11, 10, 27, 8}, rd: []int{15}, fl: []int{26}, fd: []int{3, 4, 0, 17, 16, 25, 21, 18, 6}},
		{k: 4, rl: []int{19, 13, 23, 9, 14, 11, 10, 27}, rd: []int{8, 15}, fl: []int{4, 26}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}},
		{k: 5, rl: []int{5, 19, 13, 23, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{4, 26}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}},
		{k: 26, rl: []int{5, 19, 13, 23, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{26, 4}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}},
		{k: 23, rl: []int{5, 19, 13, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{23, 26, 4}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}},
		{k: 28, rl: []int{28, 5, 19, 13, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}},
		{k: 22, rl: []int{22, 28, 5, 19, 13, 9, 14, 11}, rd: []int{10, 27}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}},
		{k: 8, rl: []int{8, 22, 28, 5, 19, 13, 9, 14}, rd: []int{11, 10}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}},
		{k: 1, rl: []int{1, 8, 22, 28, 5, 19, 13, 9}, rd: []int{14, 11}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}},
		{k: 27, rl: []int{27, 1, 8, 22, 28, 5, 19, 13}, rd: []int{9, 14}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}},
		{k: 3, rl: []int{27, 1, 8, 22, 28, 5, 19}, rd: []int{13, 9, 14}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}},
		{k: 11, rl: []int{11, 27, 1, 8, 22, 28, 5}, rd: []int{19, 13, 9}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}},
		{k: 10, rl: []int{10, 11, 27, 1, 8, 22, 28}, rd: []int{5, 19, 13}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}},
		{k: 25, rl: []int{10, 11, 27, 1, 8, 22}, rd: []int{28, 5, 19, 13}, fl: []int{25, 3, 23, 26}, fd: []int{4, 0, 17, 16, 21, 18}},
		{k: 26, rl: []int{10, 11, 27, 1, 8, 22}, rd: []int{28, 5, 19, 13}, fl: []int{26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21, 18}},
		{k: 12, rl: []int{12, 10, 11, 27, 1, 8}, rd: []int{22, 28, 5, 19}, fl: []int{26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21, 18}},
		{k: 18, rl: []int{12, 10, 11, 27, 1}, rd: []int{8, 22, 28, 5, 19}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}},
		{k: 24, rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28, 5}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}},
		{k: 18, rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28, 5}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}},
		{k: 5, rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28}, fl: []int{5, 18, 26, 25, 3}, fd: []int{23, 4, 0, 17, 16, 21}},
		{k: 8, rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}},
		{k: 8, rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}},
		{k: 12, rl: []int{24, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{12, 8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}},
		{k: 17, rl: []int{24, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{17, 12, 8, 5, 18, 26}, fd: []int{25, 3, 23, 4, 0, 16, 21}},
		{k: 0, rl: []int{24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{0, 17, 12, 8, 5, 18, 26}, fd: []int{25, 3, 23, 4, 16, 21}},
		{k: 7, rl: []int{7, 24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{0, 17, 12, 8, 5, 18}, fd: []int{26, 25, 3, 23, 4, 16}},
		{k: 8, rl: []int{7, 24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{8, 0, 17, 12, 5, 18}, fd: []int{26, 25, 3, 23, 4, 16}},
		{k: 1, rl: []int{7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17, 12, 5}, fd: []int{18, 26, 25, 3, 23, 4, 16}},
		{k: 9, rl: []int{9, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17, 12}, fd: []int{5, 18, 26, 25, 3, 23, 4}},
		{k: 29, rl: []int{29, 9, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17}, fd: []int{12, 5, 18, 26, 25, 3, 23}},
		{k: 9, rl: []int{29, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{9, 1, 8, 0, 17}, fd: []int{12, 5, 18, 26, 25, 3, 23}},
		{k: 28, rl: []int{29, 7, 24, 10, 11}, rd: []int{27, 22}, fl: []int{28, 9, 1, 8, 0}, fd: []int{17, 12, 5, 18, 26, 25, 3, 23}},
		{k: 27, rl: []int{29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9, 1, 8}, fd: []int{0, 17, 12, 5, 18, 26, 25, 3, 23}},
		{k: 13, rl: []int{13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9, 1}, fd: []int{8, 0, 17, 12, 5, 18, 26, 25, 3}},
		{k: 21, rl: []int{21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9}, fd: []int{1, 8, 0, 17, 12, 5, 18, 26, 25}},
		{k: 8, rl: []int{21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8, 27, 28}, fd: []int{9, 1, 0, 17, 12, 5, 18, 26, 25}},
		{k: 3, rl: []int{3, 21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8, 27}, fd: []int{28, 9, 1, 0, 17, 12, 5, 18, 26}},
		{k: 2, rl: []int{2, 3, 21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8}, fd: []int{27, 28, 9, 1, 0, 17, 12, 5, 18}},
		{k: 0, rl: []int{2, 3, 21, 13, 29, 7, 24, 10}, rd: []int{11, 22}, fl: []int{0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5, 18}},
		{k: 18, rl: []int{2, 3, 21, 13, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}},
		{k: 13, rl: []int{2, 3, 21, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{13, 18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}},
		{k: 21, rl: []int{2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{21, 13, 18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}},
		{k: 5, rl: []int{2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{5, 21, 13, 18, 0}, fd: []int{8, 27, 28, 9, 1, 17, 12}},
		{k: 14, rl: []int{14, 2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{5, 21, 13, 18}, fd: []int{0, 8, 27, 28, 9, 1, 17}},
		{k: 3, rl: []int{14, 2, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{3, 5, 21, 13, 18}, fd: []int{0, 8, 27, 28, 9, 1, 17}},
		{k: 4, rl: []int{4, 14, 2, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{3, 5, 21, 13}, fd: []int{18, 0, 8, 27, 28, 9, 1}},
		{k: 8, rl: []int{4, 14, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{8, 3, 5, 21, 13}, fd: []int{18, 0, 27, 28, 9, 1}},
		{k: 19, rl: []int{19, 4, 14, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{8, 3, 5, 21}, fd: []int{13, 18, 0, 27, 28, 9}},
		{k: 14, rl: []int{19, 4, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{14, 8, 3, 5, 21}, fd: []int{13, 18, 0, 27, 28, 9}},
		{k: 11, rl: []int{19, 4, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{11, 14, 8, 3, 5}, fd: []int{21, 13, 18, 0, 27, 28, 9}},
		{k: 4, rl: []int{19, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{4, 11, 14, 8, 3, 5}, fd: []int{21, 13, 18, 0, 27, 28, 9}},
		{k: 17, rl: []int{17, 19, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{4, 11, 14, 8, 3}, fd: []int{5, 21, 13, 18, 0, 27, 28}},
		{k: 19, rl: []int{17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14, 8, 3}, fd: []int{5, 21, 13, 18, 0, 27, 28}},
		{k: 12, rl: []int{12, 17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14, 8}, fd: []int{3, 5, 21, 13, 18, 0, 27}},
		{k: 23, rl: []int{23, 12, 17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14}, fd: []int{8, 3, 5, 21, 13, 18, 0}},
		{k: 29, rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{29, 19, 4, 11, 14}, fd: []int{8, 3, 5, 21, 13, 18, 0}},
		{k: 25, rl: []int{25, 23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{29, 19, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}},
		{k: 25, rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{25, 29, 19, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}},
		{k: 19, rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{19, 25, 29, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}},
		{k: 23, rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{23, 19, 25, 29, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}},
		{k: 11, rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{11, 23, 19, 25, 29, 4}, fd: []int{14, 8, 3, 5, 21, 13, 18}},
		{k: 18, rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{18, 11, 23, 19, 25, 29}, fd: []int{4, 14, 8, 3, 5, 21, 13}},
		{k: 24, rl: []int{12, 17, 2, 7}, rd: []int{10, 22}, fl: []int{24, 18, 11, 23, 19, 25}, fd: []int{29, 4, 14, 8, 3, 5, 21, 13}},
		{k: 13, rl: []int{12, 17, 2, 7}, rd: []int{10, 22}, fl: []int{13, 24, 18, 11, 23, 19}, fd: []int{25, 29, 4, 14, 8, 3, 5, 21}},
		{k: 12, rl: []int{17, 2, 7}, rd: []int{10, 22}, fl: []int{12, 13, 24, 18, 11, 23, 19}, fd: []int{25, 29, 4, 14, 8, 3, 5, 21}},
		{k: 27, rl: []int{27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{12, 13, 24, 18, 11, 23}, fd: []int{19, 25, 29, 4, 14, 8, 3, 5}},
		{k: 19, rl: []int{27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{19, 12, 13, 24, 18, 11}, fd: []int{23, 25, 29, 4, 14, 8, 3, 5}},
		{k: 16, rl: []int{16, 27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{19, 12, 13, 24, 18}, fd: []int{11, 23, 25, 29, 4, 14, 8, 3}},
		{k: 3, rl: []int{16, 27, 17, 2}, rd: []int{7, 10, 22}, fl: []int{3, 19, 12, 13, 24, 18}, fd: []int{11, 23, 25, 29, 4, 14, 8}},
		{k: 7, rl: []int{16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12, 13, 24}, fd: []int{18, 11, 23, 25, 29, 4, 14, 8}},
		{k: 1, rl: []int{1, 16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12, 13}, fd: []int{24, 18, 11, 23, 25, 29, 4, 14}},
		{k: 26, rl: []int{26, 1, 16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12}, fd: []int{13, 24, 18, 11, 23, 25, 29, 4}},
		{k: 2, rl: []int{26, 1, 16, 27, 17}, rd: []int{10, 22}, fl: []int{2, 7, 3, 19, 12}, fd: []int{13, 24, 18, 11, 23, 25, 29, 4}},
		{k: 15, rl: []int{15, 26, 1, 16, 27, 17}, rd: []int{10, 22}, fl: []int{2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 23, 25, 29}},
		{k: 23, rl: []int{15, 26, 1, 16, 27}, rd: []int{17, 10, 22}, fl: []int{23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 25, 29}},
		{k: 23, rl: []int{15, 26, 1, 16, 27}, rd: []int{17, 10, 22}, fl: []int{23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 25, 29}},
		{k: 18, rl: []int{15, 26, 1, 16}, rd: []int{27, 17, 10, 22}, fl: []int{18, 23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 11, 25, 29}},
		{k: 1, rl: []int{15, 26, 16}, rd: []int{27, 17, 10, 22}, fl: []int{1, 18, 23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 11, 25, 29}},
		{k: 12, rl: []int{15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7, 3, 19}, fd: []int{13, 24, 11, 25, 29}},
		{k: 20, rl: []int{20, 15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7, 3}, fd: []int{19, 13, 24, 11, 25}},
		{k: 21, rl: []int{21, 20, 15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7}, fd: []int{3, 19, 13, 24, 11}},
		{k: 6, rl: []int{6, 21, 20, 15}, rd: []int{26, 16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7}, fd: []int{3, 19, 13, 24}},
		{k: 3, rl: []int{6, 21, 20}, rd: []int{15, 26, 16, 27, 17, 10, 22}, fl: []int{3, 12, 1, 18, 23, 2, 7}, fd: []int{19, 13, 24}},
		{k: 19, rl: []int{6, 21}, rd: []int{20, 15, 26, 16, 27, 17, 10, 22}, fl: []int{19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}},
		{k: 19, rl: []int{6, 21}, rd: []int{20, 15, 26, 16, 27, 17, 10, 22}, fl: []int{19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}},
		{k: 15, rl: []int{6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{15, 19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}},
		{k: 8, rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{15, 19, 3, 12, 1, 18, 23, 2}, fd: []int{7, 13}},
		{k: 3, rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{3, 15, 19, 12, 1, 18, 23, 2}, fd: []int{7, 13}},
		{k: 23, rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{7, 13}},
		{k: 29, rl: []int{29, 8}, rd: []int{6, 21, 20, 26, 16, 27, 17, 10}, fl: []int{23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{7, 13}},
		{k: 7, rl: []int{29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17, 10}, fl: []int{7, 23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{13}},
		{k: 10, rl: []int{29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1, 18}, fd: []int{2, 13}},
		{k: 9, rl: []int{9, 29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1}, fd: []int{18, 2}},
		{k: 28, rl: []int{28, 9}, rd: []int{29, 8, 6, 21, 20, 26, 16, 27}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1}, fd: []int{18, 2}},
		{k: 6, rl: []int{28, 9}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{6, 10, 7, 23, 3, 15, 19, 12}, fd: []int{1, 18, 2}},
		{k: 15, rl: []int{28, 9}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{15, 6, 10, 7, 23, 3, 19, 12}, fd: []int{1, 18, 2}},
		{k: 9, rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{9, 15, 6, 10, 7, 23, 3, 19, 12}, fd: []int{1, 18, 2}},
		{k: 3, rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{3, 9, 15, 6, 10, 7, 23, 19, 12}, fd: []int{1, 18, 2}},
		{k: 6, rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{6, 3, 9, 15, 10, 7, 23, 19, 12}, fd: []int{1, 18, 2}},
		{k: 20, rl: []int{28}, rd: []int{29, 8, 21, 26, 16, 27}, fl: []int{20, 6, 3, 9, 15, 10, 7, 23, 19}, fd: []int{12, 1, 18, 2}},
		{k: 8, rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{8, 20, 6, 3, 9, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}},
		{k: 6, rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{6, 8, 20, 3, 9, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}},
		{k: 9, rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{9, 6, 8, 20, 3, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}},
		{k: 3, rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{3, 9, 6, 8, 20, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}},
		{k: 26, rl: []int{28}, rd: []int{29, 21, 16, 27}, fl: []int{26, 3, 9, 6, 8, 20, 15, 10, 7}, fd: []int{23, 19, 12, 1, 18, 2}},
		{k: 20, rl: []int{28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15, 10, 7}, fd: []int{23, 19, 12, 1, 18, 2}},
		{k: 17, rl: []int{17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15, 10}, fd: []int{7, 23, 19, 12, 1, 18}},
		{k: 22, rl: []int{22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15}, fd: []int{10, 7, 23, 19, 12, 1}},
		{k: 10, rl: []int{22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{10, 20, 26, 3, 9, 6, 8}, fd: []int{15, 7, 23, 19, 12, 1}},
		{k: 2, rl: []int{2, 22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 7, 23, 19, 12}},
		{k: 7, rl: []int{2, 22, 17}, rd: []int{28, 29, 21, 16, 27}, fl: []int{7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 23, 19, 12}},
		{k: 7, rl: []int{2, 22, 17}, rd: []int{28, 29, 21, 16, 27}, fl: []int{7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 23, 19, 12}},
		{k: 15, rl: []int{2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{15, 7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 23, 19, 12}},
		{k: 1, rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{15, 7, 10, 20, 26, 3, 9}, fd: []int{6, 8, 23, 19}},
		{k: 26, rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{26, 15, 7, 10, 20, 3, 9}, fd: []int{6, 8, 23, 19}},
		{k: 21, rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 16, 27}, fl: []int{21, 26, 15, 7, 10, 20, 3}, fd: []int{9, 6, 8, 23, 19}},
	}
	ints := func(l *list.List) []int {
		a := make([]int, l.Len())
		el := l.Front()
		for i := range a {
			a[i] = el.Value.(*item).key.(int)
			el = el.Next()
		}
		return a
	}
	c := New(10)
	for i, want := range table {
		k := want.k
		if _, ok := c.Get(k); !ok {
			c.Set(k, k)
		}
		got := result{k: k, rl: ints(c.rl), rd: ints(c.rd), fl: ints(c.fl), fd: ints(c.fd)}
		if !reflect.DeepEqual(got, want) {
			var prev result
			if i > 0 {
				prev = table[i-1]
			}
			t.Errorf("unexpected state: step %d\nprev %+v\ngot  %+v\nwant %+v", i, prev, got, want)
		}
	}
}
