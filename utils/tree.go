package utils

import (
	"cointrade/lib/db"
	"fmt"
)

type ListTree struct {
	Id       int         `json:"id"`
	Path     string      `json:"path"`
	Redirect string      `json:"redirect"`
	Level    int         `json:"level"`
	Icon     string      `json:"icon"`
	Topid    int         `json:"pid"`
	Children []*ListTree `json:"children"`
	Flag     string      `json:"flag"`
	Hidden   int         `json:"hidden"`
	Name     string      `json:"name"`
	Weight   int         `json:"weight"`
	Status   string      `json:"status"`
	Meta     Meta        `json:"meta"`
}

type Meta struct {
	Title string   `json:"title"`
	Icon  string   `json:"icon"`
	Roles []string `json:"roles"`
}

type Mean struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Hidden   int      `json:"hidden"`
	ParentId int      `json:"parent_id"`
	Icon     string   `json:"icon"`
	Status   string   `json:"status"`
	Weight   int      `json:"weight"`
	Roleid   []string `json:"role"`
}

type TreeTool struct {
	ListMean []*Mean
}

func (t *TreeTool) SetTreen(list []db.DBValues) *TreeTool {
	for _, tree := range list {
		m := new(Mean)
		tree.SetObj(m)
		t.ListMean = append(t.ListMean, m)
	}
	return t
}

func (t *TreeTool) GetTop(top int) []*Mean {
	rs := make([]*Mean, 0)
	for _, item := range t.ListMean {
		if item.ParentId == top {
			rs = append(rs, item)
		}
	}
	return rs
}

func (t *TreeTool) GetTree(topid int, level int) []*ListTree {
	rs := make([]*ListTree, 0)
	topList := t.GetTop(topid)
	level++
	for _, top := range topList {
		tree := new(ListTree)
		tree.Children = []*ListTree{}
		tree.Path = top.Path
		tree.Level = level
		tree.Name = top.Name
		tree.Topid = top.ParentId
		tree.Id = top.Id
		tree.Hidden = top.Hidden
		tree.Weight = top.Weight
		child := t.GetTree(top.Id, level)
		tree.Status = top.Status
		if len(child) > 0 {
			tree.Children = append(tree.Children, child...)
		}
		tree.Meta = Meta{
			Title: top.Name,
			Icon:  top.Icon,
			Roles: top.Roleid,
		}
		rs = append(rs, tree)
	}
	return rs
}

func (t *TreeTool) GetMeanList(list []*ListTree, level int) []*ListTree {
	rs := make([]*ListTree, 0)

	for _, tree := range list {
		tree.Flag = t.SpaceZ(tree.Level)
		if level == 0 {
			tree.Flag += "┌"
		} else {
			tree.Flag += "├"
		}

		rs = append(rs, tree)
		if len(tree.Children) > 0 {
			level++
			rs = append(rs, t.GetMeanList(tree.Children, level)...)
		}
		if tree.Level == 1 {
			level = 0
		}

	}
	return rs
}

func (t *TreeTool) SpaceZ(num int) string {
	str := ""
	if num == 1 {
		return ""
	}
	for i := 0; i < num; i++ {
		str += "    ."
	}
	return fmt.Sprintf("%s", str)
}
