package go_mdict

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildRangeTree(t *testing.T) {
	root := new(RecordBlockRangeTreeNode)
	list := []*MdictRecordBlockInfoListItem{
		{
			compressSize:                10,
			deCompressSize:              100,
			compressAccumulatorOffset:   0,
			deCompressAccumulatorOffset: 10,
		}, {
			compressSize:                10,
			deCompressSize:              100,
			compressAccumulatorOffset:   10,
			deCompressAccumulatorOffset: 10,
		}, {
			compressSize:                10,
			deCompressSize:              100,
			compressAccumulatorOffset:   20,
			deCompressAccumulatorOffset: 10,
		}, {
			compressSize:                10,
			deCompressSize:              100,
			compressAccumulatorOffset:   30,
			deCompressAccumulatorOffset: 10,
		},
	}

	BuildRangeTree(list, root)

	printRangeTreeOverLevel(t, root)

}

func TestBuildRangeTree2(t *testing.T) {
	root := new(RecordBlockRangeTreeNode)
	dict, err := New("./testdata/mdx/testdict.mdx")
	if err != nil {
		t.Fatal(err)
	}
	err = dict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for i := 0; i < len(dict.recordBlockInfo.recordInfoList)-1; i++ {
		curr := dict.recordBlockInfo.recordInfoList[i]
		next := dict.recordBlockInfo.recordInfoList[i+1]
		if curr.deCompressAccumulatorOffset+curr.deCompressSize != next.deCompressAccumulatorOffset {
			count++
		}

	}
	t.Logf("not equal %d", count)

	BuildRangeTree(dict.recordBlockInfo.recordInfoList, root)

	printRangeTreeOverLevel(t, root)

	data := QueryRangeData(root, 162512550)
	t.Logf("%v", data)
	count = 0
	for _, e := range dict.keyBlockData.keyEntries {
		data = QueryRangeData(root, e.RecordStartOffset)
		if data == nil {
			count++
		}
	}
	t.Logf("nil count: %d", count)
	assert.Equal(t, 0, count)

}

// 按层遍历二叉树

func printRangeTreeOverLevel(t *testing.T, root *RecordBlockRangeTreeNode) {
	t.Logf("level(%d)\t[%d,%d]", 0, root.startRange, root.endRange)
	level := 0
	levels := make([][]*RecordBlockRangeTreeNode, 0)
	print(&levels, root, level)
	for le, ls := range levels {
		line := ""
		line += fmt.Sprintf("level(%d)\t", le+1)
		for _, l := range ls {
			if l == nil {
				line += fmt.Sprintf("\t")
			} else {
				line += fmt.Sprintf("\t[%d, %d](%v)", l.startRange, l.endRange, l.data)
			}
		}
		t.Logf("%s", line)
	}

}

func print(levels *([][]*RecordBlockRangeTreeNode), root *RecordBlockRangeTreeNode, level int) {
	if root == nil {
		return
	}
	if len(*levels)-1 < level {
		(*levels) = append((*levels), make([]*RecordBlockRangeTreeNode, 0))
	}

	(*levels)[level] = append((*levels)[level], root.left)
	(*levels)[level] = append((*levels)[level], root.right)

	if root.left != nil {
		print(levels, root.left, level+1)
	}
	if root.right != nil {
		print(levels, root.right, level+1)
	}

}
