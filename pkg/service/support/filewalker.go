package support

import (
	"github.com/terasum/medict/pkg/model"
	"io/fs"
	"path/filepath"
)

func WalkDir(dirpath string) ([]*model.DirItem, error) {
	list := make([]*model.DirItem, 0)
	err := filepath.WalkDir(dirpath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip self
		if dirpath == path {
			return nil
		}
		// skip non-dir
		if !d.IsDir() {
			return nil
		}
		item, err := innerWalker(dirpath, path, err)
		list = append(list, item)
		return nil
	})
	return list, err
}

func innerWalker(rootpath, subpath string, err error) (*model.DirItem, error) {
	if err != nil {
		return nil, err
	}
	item := &model.DirItem{
		BaseDir:    rootpath,
		CurrentDir: subpath,
		MdxAbsPath: "",
		MddAbsPath: make([]string, 0),
	}

	err = filepath.Walk(subpath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// if mdx
		if filepath.Ext(info.Name()) == ".mdx" {
			item.MdxAbsPath, _ = filepath.Abs(path)
		} else if filepath.Ext(info.Name()) == ".mdd" {
			// pop stack
			mddabs, _ := filepath.Abs(path)
			item.MddAbsPath = append(item.MddAbsPath, mddabs)
		} else {
			// TODO copy css/js files
		}
		return nil
	})
	return item, err
}
