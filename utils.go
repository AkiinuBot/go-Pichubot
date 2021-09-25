package pichubot

import (
	"os"
)

// 判断文件夹是否存在 若不存在则新建
func CheckPath(path string) error {
	exist, err := pathExists(path)
	if !exist && err == nil {
		return (os.MkdirAll(path, os.ModePerm))
	}
	return err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*

func findFile(directoryPath string) []string {
	baseFile, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		Logger.Error("An error occurred while open file :[" + directoryPath + "] .")
		Logger.Error(err.Error())
		return nil
	}
	var res []string
	for _, fileItem := range baseFile {
		if fileItem.IsDir() {
			innerFiles := findFile(path.Join(directoryPath, fileItem.Name()))
			res = append(res, innerFiles...)
		} else {
			res = append(res, path.Join(directoryPath, fileItem.Name()))
		}
	}
	return res
}


//! WIP - plugin

// PluginItem 存储着插件的信息
// type PluginItem struct {
// 	Name       string
// 	TargetFunc plugin.Symbol
// }

// 所有插件必须实现该方法
const TargetFuncName = "RegFunction"

// LoadAllPlugin 将会过滤一次传入的targetFile,同时将so后缀的文件装载，并返回一个插件信息集合
func LoadAllPlugin(targetFile []string) []Plugin {
	var res []Plugin

	for _, fileItem := range targetFile {
		// 过滤插件文件
		if path.Ext(fileItem) == "so" {
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				Logger.Error("An error occurred while load plugin : [" + fileItem + "]")
				Logger.Error(err.Error())
			}

			//查找指定函数或符号
			targetFunc, err := pluginFile.Lookup(TargetFuncName)
			if err != nil {
				Logger.Error("An error occurred while search target func : [" + fileItem + "]")
				Logger.Error(err.Error())
			}

			//查找指定函数或符号
			listener, err := pluginFile.Lookup("Linsters")
			if err != nil {
				Logger.Error("An error occurred while search target func : [" + fileItem + "]")
				Logger.Error(err.Error())
			}

			//查找指定函数或符号
			version, err := pluginFile.Lookup("Version")
			if err != nil {
				Logger.Error("An error occurred while search target func : [" + fileItem + "]")
				Logger.Error(err.Error())
			}

			//采集插件信息
			pluginInfo := Plugin{
				Name:       fileItem,
				TargetFunc: targetFunc,
				Version:    version.(string),
				Listener:   listener.(listeners),
			}

			// 进行调用
			if f, ok := targetFunc.(func()); ok {
				f()
			}

			res = append(res, pluginInfo)
		}
	}
	return res
}

var PluginItems []Plugin

func init() {
	Logger.Info("Start Loading Plugins")
	CheckPath("./plugin")
	pluginFiles := findFile("./plugin")
	PluginItems = LoadAllPlugin(pluginFiles)
	Logger.Info("Appending Listeners...")
	for _, plugin := range PluginItems {
		Listeners.OnPrivateMsg = append(Listeners.OnPrivateMsg, plugin.Listener.OnPrivateMsg...)
		Listeners.OnGroupMsg = append(Listeners.OnGroupMsg, plugin.Listener.OnGroupMsg...)
		Listeners.OnGroupUpload = append(Listeners.OnGroupUpload, plugin.Listener.OnGroupUpload...)
		Listeners.OnGroupAdmin = append(Listeners.OnGroupAdmin, plugin.Listener.OnGroupAdmin...)
		Listeners.OnGroupDecrease = append(Listeners.OnGroupDecrease, plugin.Listener.OnGroupDecrease...)
		Listeners.OnGroupIncrease = append(Listeners.OnGroupIncrease, plugin.Listener.OnGroupIncrease...)
		Listeners.OnGroupBan = append(Listeners.OnGroupBan, plugin.Listener.OnGroupBan...)
		Listeners.OnFriendAdd = append(Listeners.OnFriendAdd, plugin.Listener.OnFriendAdd...)
		Listeners.OnGroupRecall = append(Listeners.OnGroupRecall, plugin.Listener.OnGroupRecall...)
		Listeners.OnFriendRecall = append(Listeners.OnFriendRecall, plugin.Listener.OnFriendRecall...)
		Listeners.OnNotify = append(Listeners.OnNotify, plugin.Listener.OnNotify...)
		Listeners.OnFriendRequest = append(Listeners.OnFriendRequest, plugin.Listener.OnFriendRequest...)
		Listeners.OnGroupRequest = append(Listeners.OnGroupRequest, plugin.Listener.OnGroupRequest...)
		Listeners.OnMetaLifecycle = append(Listeners.OnMetaLifecycle, plugin.Listener.OnMetaLifecycle...)
		Listeners.OnMetaHeartbeat = append(Listeners.OnMetaHeartbeat, plugin.Listener.OnMetaHeartbe
			at...)
	}
	Logger.Info("All Plugins' Listeners are loaded")
}
*/
