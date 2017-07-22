package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//"path/filepath"
	"syscall"
	"time"
)

//先假定上次更新时间为:
var last_synchronous_time time.Time

//操作记录,记录发现的所有文件变化

var ActionList []even

//-------------------------------------------------------------------文件操作类型------------------------------------------------------------------------------------
const (
	create = iota //0
	delet         //1
	modify        //2
	rename        //3
)

type even struct {
	Op    int
	level int
	Path  string
	temp  string //暂存一些数据,比如文件重命名后的新名字
}

//将该操作类型记录入全局变量中

func (eve *even) save() {
	ActionList = append(ActionList, *eve)
}

//--------------------------------------------------------------------历史树结构和配置结构-----------------------------------------------------------------------------
type History struct {
	Synchronous_time time.Time
	Dir_conternt     Dir
}

//--------------------------------------------------------------------建立目录树装结构,用做对比-------------------------------------------------------------------------
type Dir struct {
	Id         int64
	ParentId   int64
	Path       string
	Name       string
	Children   []Dir
	IsDir      bool
	ModifyTime time.Time
}

func (dir *Dir) getIDs() []int64 {
	var ids []int64
	for _, child := range dir.Children {
		ids = append(ids, child.Id)
	}
	return ids
}

func (dir *Dir) getChildById(id int64) Dir {
	for _, child := range dir.Children {
		if child.Id == id {
			return child
		}
	}
	//未找到指定dir
	return Dir{}
}

var pathSep string = string(os.PathSeparator)

//使用文件的创建时间(纳秒记录)作为文件的唯一id
var data syscall.Win32finddata

func getFileId(path string) int64 {
	uintpath := syscall.StringToUTF16Ptr(path)
	handle, err := syscall.FindFirstFile(uintpath, &data)
	defer syscall.CloseHandle(handle)
	if err != nil {
		fmt.Println("--->", err)
	}
	return data.CreationTime.Nanoseconds()
}

func pathWalk(path string, parentId int64) []Dir {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	var children []Dir
	var childPath string
	for _, dir := range dirs {
		childPath = path + pathSep + dir.Name()
		children = append(children, getFileInfo(childPath, parentId))
	}
	return children
}

//此函数作为扫描文件模块的入口
func getFileInfo(path string, parentId int64) Dir {
	fileInfo, errs := os.Stat(path)
	if errs != nil {
		fmt.Println(errs)
	}
	var self Dir
	self = Dir{
		Id:         getFileId(path),
		ParentId:   parentId,
		Path:       path,
		Name:       fileInfo.Name(),
		IsDir:      fileInfo.IsDir(),
		ModifyTime: fileInfo.ModTime(),
	}
	if self.IsDir {
		self.Children = pathWalk(path, self.Id)
	}
	return self
}

//---------------------------------------------------------------------------------------------------------------------------------------------
//将扫描结果存储在本地文件中
func buildTree(history History, path string) {
	jsonbyte, err := json.Marshal(history)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(path, jsonbyte, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func loadTree(path string) History {
	jsonbyte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	var history History
	err = json.Unmarshal(jsonbyte, &history)
	if err != nil {
		fmt.Println(err)
	}
	return history
}

//---------------------------------------------------------------载入上次同步完成后的目录树状结构,开始对比------------------------------------------------------------------------------

//传入当前和上次的文件列表,对比出新增和删除的文件id,以及原有的文件
func getADDandRemove(last, now []int64) (add []int64, remove []int64, origin []int64) {
	var checkBox map[int64]int = make(map[int64]int)
	var addFiles, reMoveFiles, originFiles []int64
	var resultContainer [][]int64 = [][]int64{}
	resultContainer = append(resultContainer, reMoveFiles, originFiles, addFiles)
	//上次同步后的目录标记为-1
	for _, fileId := range last {
		checkBox[fileId] -= 1

	}
	//本次同步发现的目录标记为1
	for _, fileId := range now {
		checkBox[fileId] += 1
	}
	//1为删除文件列表,2为新增文件列表,3为原有文件列表
	for fileId, value := range checkBox {
		resultContainer[value+1] = append(resultContainer[value+1], fileId)
	}
	return resultContainer[2], resultContainer[0], resultContainer[1]
}

func compare(base, current Dir, level int) {
	current_level_ids := current.getIDs()
	base_level_ids := base.getIDs()
	//add, remove, origin := getADDandRemove(base_level_ids, current_level_ids)
	add, remove, origin := getADDandRemove(base_level_ids, current_level_ids)
	createHandle(current, add, level)
	removeHandle(base, remove, level)
	modifyHandle(base, current, origin, level)

}

//新增处理
func createHandle(current Dir, ids []int64, level int) {
	for _, id := range ids {
		current_dir := current.getChildById(id)
		tempEven := even{
			Op:    create,
			Path:  current_dir.Path,
			level: level,
		}

		tempEven.save()
	}

}

//删除处理
func removeHandle(base Dir, ids []int64, level int) {
	for _, id := range ids {
		base_dir := base.getChildById(id)
		tempEven := even{
			Op:    delet,
			Path:  base_dir.Path,
			level: level,
		}

		tempEven.save()
	}

}

//修改处理,包括重命名
func modifyHandle(base, current Dir, ids []int64, level int) {
	//在原有文件中循环查找更改过的文件(修改或重命名)
	for _, id := range ids {
		//根据id对比操作文件
		base_dir := base.getChildById(id)
		current_dir := current.getChildById(id)
		//先判断原有文件是否被重命名过
		var tempEven even
		if base_dir.Name != current_dir.Name {
			//确定文件名被更改过,记录文件与对应操作类型
			//构成当前文件的完整路径
			tempEven = even{
				Op:    rename,
				Path:  base_dir.Path,
				temp:  current_dir.Name,
				level: level,
			}
			//fmt.Println("tmpstruct:", tempEven)
			//将自身加入到待更新列表中
			tempEven.save()
		}
		//比较原有的修改时间,确定处理方式
		//发现文件的最后更新时间大于上次记录的更新时间,确其被更改过
		//暂时放弃单独筛选出被修改文件,采取遍历全部文件
		//if last_synchronous_time.Before(current_dir.ModifyTime) {
		//fmt.Println("find file has been modify:" + current_dir.Name)
		//确定文件被更改过,根据文件类型区分下次操作
		if base_dir.IsDir {
			//fmt.Println("into circle")
			//该文件类型是--->文件夹.递归检测其子目录下的更新情况
			compare(base_dir, current_dir, level+1)
		} else {
			//该文件不是文件夹,直接记录文件路径,等待整理
			//暂时放弃单独筛选出被修改文件,采取遍历全部文件
			if last_synchronous_time.Before(current_dir.ModifyTime) {
				tempEven = even{
					Op:   modify,
					Path: current_dir.Path + pathSep + current_dir.Name,
				}
				//将自身加入到待更新列表中
				tempEven.save()
			}
		}
		//}
	}
}

//---------------------------查看文件是否存在------------------------------------------
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//---------------------------------------------------------------------------------------------------------------------------------------------

//生成待处理序列,准备发送
func scanf(dir_path, tree_path string) []even {
	exist_flag := PathExists(tree_path)
	var history History
	dirs := getFileInfo(dir_path, -1)
	if !exist_flag {
		fmt.Println("config file not exist,build it")
		history = History{
			Synchronous_time: time.Now(),
			Dir_conternt:     dirs,
		}

		buildTree(history, tree_path)
		fmt.Println("build tree success")
		return nil
	} else {
		fmt.Println("start to load tree")
		history = loadTree(tree_path)
		last_synchronous_time = history.Synchronous_time
		fmt.Println("start to compare dir")
		compare(history.Dir_conternt, dirs, 1)

		return ActionList
	}

}

/*func draw(eve []even) {
	for _, e := range eve {
		fmt.Println("op:", e.Op, ",level", e.level, ",path", e.Path, ",temp:", e.temp)
	}
}*/

func main() {
	tmp := scanf(".\\tt", ".\\Config.Ame")
	if tmp != nil {
		for _, m := range tmp {
			fmt.Println(m)
		}
	}
}
