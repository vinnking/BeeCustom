package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"BeeCustom/controllers"
	"BeeCustom/enums"
	"BeeCustom/file"
	"BeeCustom/models"
	"BeeCustom/utils"
	"BeeCustom/xmlTemplate"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

//回执解释
func annotationCReturnXmlParse() *toolbox.Task {

	task := toolbox.NewTask("task", "* * * * * *", func() error {
		parseAnnotationReturns("annotation_return_path_c", "annotation_history_path_c")
		return nil
	})

	err := task.Run()
	if err != nil {
		utils.LogError(fmt.Sprintf("tk.Run error :%v", err))
		return nil
	}

	return task
}

//解析回执
func parseAnnotationReturns(returnPathCofig, historyPathCofig string) {

	returnPath := beego.AppConfig.String(returnPathCofig)
	historyPath := beego.AppConfig.String(historyPathCofig)
	pathCfiles, err := ioutil.ReadDir(returnPath)
	if err != nil {
		utils.LogError(fmt.Sprintf("获取数据列表和总数 error:%v", err))
	}

	fileIndex := 0
	if len(pathCfiles) >= fileIndex {
		f := pathCfiles[fileIndex]
		fullPath := returnPath + f.Name()
		xmlFile, err := os.Open(fullPath) // For read access.
		if err != nil {
			utils.LogError(fmt.Sprintf("os.Open :%v", err))
			fileIndex++
		}

		defer xmlFile.Close()

		data, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			utils.LogError(fmt.Sprintf(" ioutil.ReadAll :%v", err))
			fileIndex++
		}

		prefix, ext, failedName := getNameExts(f)
		if len(prefix) > 0 && len(ext) > 0 {
			if ext == "INV" {
				if prefix == "Successed" {
					v := xmlTemplate.CommonResponeMessage{}
					err := xml.Unmarshal(data, &v)
					if err != nil {
						utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
						fileIndex++
					}

					annotation, err := models.GetAnnotationByEtpsInnerInvtNo(v.EtpsPreentNo)
					if err != nil {
						utils.LogError(fmt.Sprintf(" models.GetAnnotationByEtpsInnerInvtNo :%v,filename:%v", err, f.Name()))
						fileIndex++
					}

					aReturn := models.NewAnnotationReturn(0)
					aReturn.EtpsPreentNo = v.EtpsPreentNo
					aReturn.CheckInfo = v.CheckInfo
					aReturn.DealFlag = v.DealFlag
					aReturn.Annotation = annotation

					if err = models.AnnotationReturnSave(&aReturn); err != nil {
						utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
						fileIndex++
					}

					//更新状态
					if err = controllers.UpdateAnnotationStatus(annotation, "单一处理中"); err != nil {
						utils.LogError(fmt.Sprintf("controllers.UpdateAnnotationStatus :%v,filename:%v", err, f.Name()))
					}

					err = moveFile(historyPath, v.EtpsPreentNo, fullPath, f)
					if err != nil {
						utils.LogError(fmt.Sprintf("moveFile:%v,filename:%v", err, f.Name()))
						fileIndex++
					}

					//ws 自动更新
					wsPush()

				} else if prefix == "Receipt" {
					err = moveFile(historyPath, "Receipt", fullPath, f)
					if err != nil {
						utils.LogError(fmt.Sprintf("moveFile :%v,filename:%v", err, f.Name()))
						fileIndex++
					}
				}
			} else if ext == "INVT" && prefix == "Receipt" {
				v := xmlTemplate.ReturnPackage{}
				err := xml.Unmarshal(data, &v)
				if err != nil {
					utils.LogError(fmt.Sprintf(" xml.Unmarshal :%v,filename:%v", err, f.Name()))
					fileIndex++
				}

				annotation, err := models.GetAnnotationBySeqNo(v.InvPreentNo)
				if err != nil {
					utils.LogError(fmt.Sprintf("models.GetAnnotationBySeqNo :%v,filename:%v", err, f.Name()))
					fileIndex++
				}

				aReturn := models.NewAnnotationReturn(0)
				aReturn.BusinessId = v.DataInfoBusinessId
				aReturn.ManageResult = v.ManageResult
				aReturn.Reason = v.Reason
				aReturn.Annotation = annotation

				if err = models.AnnotationReturnSave(&aReturn); err != nil {
					utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
					fileIndex++
				}

				//更新状态
				if err = controllers.UpdateAnnotationStatus(annotation, "已完成"); err != nil {
					utils.LogError(fmt.Sprintf("enums.UpdateAnnotationStatus :%v,filename:%v", err, f.Name()))
				}

				err = moveFile(historyPath, annotation.EtpsInnerInvtNo, fullPath, f)
				if err != nil {
					utils.LogError(fmt.Sprintf("os.Rename :%v,filename:%v", err, f.Name()))
					fileIndex++
				}

				//ws 自动更新
				msg := utils.Message{"清单状态更新", true}
				utils.Broadcast <- msg

			} else if ext == "xml" {
				if prefix == "Failed" {
					v := xmlTemplate.Failed{}
					err := xml.Unmarshal(data, &v)
					if err != nil {
						utils.LogError(fmt.Sprintf("xml.Unmarshal :%v ,filename:%v", err, f.Name()))
						fileIndex++
					}

					annotation, err := models.GetAnnotationByEtpsInnerInvtNo(failedName)
					if err != nil {
						utils.LogError(fmt.Sprintf(" models.GetAnnotationByEtpsInnerInvtNo :%v,filename:%v,failedName:%v", err, f.Name(), failedName))
						fileIndex++
					}

					aReturn := models.NewAnnotationReturn(0)
					aReturn.CheckInfo = v.FailInfo
					aReturn.DealFlag = v.FailCode
					aReturn.Annotation = annotation

					if err = models.AnnotationReturnSave(&aReturn); err != nil {
						utils.LogError(fmt.Sprintf("models.AnnotationReturnSave :%v,filename:%v", err, f.Name()))
						fileIndex++
					}

					fileIndex++
				}
			}
		}

		err = moveFile(historyPath, "Error", fullPath, f)
		if err != nil {
			utils.LogError(fmt.Sprintf("moveFile:%v,filename:%v", err, f.Name()))
			fileIndex++
		}
		fileIndex++
	}

}

//ws 自动更新
func wsPush() {
	msg := utils.Message{"清单状态更新", true}
	utils.Broadcast <- msg
}

//移动文件
func moveFile(historyPath, v, fullPath string, f os.FileInfo) error {
	path := historyPath + time.Now().Format(enums.BaseDateFormat) + "/" + v + "/"
	if err := file.CreateFile(path); err != nil {
		utils.LogError(fmt.Sprintf("文件夹创建失败:%v", err))
	}

	err := os.Rename(fullPath, path+f.Name())
	if err != nil {
		return err
	}

	return nil
}

//回执文件的前缀，和后缀
func getNameExts(f os.FileInfo) (string, string, string) {
	extNames := strings.Split(f.Name(), ".")
	if len(extNames) > 1 && len(extNames[1]) > 0 && len(extNames[0]) > 0 {
		eNames := strings.Split(extNames[0], "_")
		if len(eNames) > 1 && len(eNames[0]) > 0 {
			if len(eNames) > 3 && len(eNames[3]) > 0 {
				return eNames[0], extNames[1], eNames[3]
			} else {
				return eNames[0], extNames[1], ""
			}
		}

	} else {
		return "", extNames[1], ""
	}

	return "", "", ""
}