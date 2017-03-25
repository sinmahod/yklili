package data

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"strings"
	"yklili/models"
)

type SystemController struct {
	DataController
}

type SysInfo struct {
	Name  string
	Value interface{}
}

func GetInfo(n string, v interface{}) *SysInfo {
	return &SysInfo{n, v}
}

func (c *SystemController) UsedPercent() {
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")

	c.put("MemoryUsedPercent", strconv.FormatFloat(v.UsedPercent, 'f', 2, 64))
	c.put("HeadUsedPercent", strconv.FormatFloat(d.UsedPercent, 'f', 2, 64))
	c.ServeJSON()
}

func (c *SystemController) List() {
	v, _ := mem.VirtualMemory()
	p, _ := cpu.Info()
	d, _ := disk.Usage("/")
	n, _ := host.Info()

	var cpu string

	if len(p) > 1 {
		for _, sub_cpu := range p {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			cpu = fmt.Sprintf("%v   %v cores ", modelname, cores)
			fmt.Println(sub_cpu)
		}
	} else {
		sub_cpu := p[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		cpu = fmt.Sprintf("%v   %v cores ", modelname, cores)
		fmt.Println(sub_cpu)
	}

	m := make([]*SysInfo, 0)
	m = append(m, GetInfo("系统信息", fmt.Sprintf("%v   %v ", n.OS, n.PlatformVersion)))
	m = append(m, GetInfo("系统总内存", fmt.Sprintf("%v GB", v.Total/1024/1024/1024)))
	m = append(m, GetInfo("剩余内存", fmt.Sprintf("%v MB", v.Free/1024/1024)))
	m = append(m, GetInfo("内存已使用", fmt.Sprintf("%f%%", v.UsedPercent)))
	m = append(m, GetInfo("CPU信息", cpu))

	m = append(m, GetInfo("硬盘总大小", fmt.Sprintf("%v GB", d.Total/1024/1024/1024)))
	m = append(m, GetInfo("硬盘剩余", fmt.Sprintf("%v GB", d.Free/1024/1024/1024)))
	m = append(m, GetInfo("硬盘已使用", fmt.Sprintf("%f%% ", d.UsedPercent)))

	m = append(m, GetInfo("服务器名称", fmt.Sprintf("%v ", n.Hostname)))

	m = append(m, GetInfo("缓存文件内存", beego.BConfig.MaxMemory/1024/1024))
	m = append(m, GetInfo("Beego版本", beego.BConfig.ServerName))
	m = append(m, GetInfo("应用所在路径", beego.AppPath))
	m = append(m, GetInfo("应用名称", beego.BConfig.AppName))
	m = append(m, GetInfo("端口", beego.BConfig.Listen.HTTPPort))
	m = append(m, GetInfo("运行模式", beego.BConfig.RunMode))

	sd := make([]string, 0)

	for k, _ := range beego.BConfig.WebConfig.StaticDir {
		sd = append(sd, k)
	}

	m = append(m, GetInfo("静态目录", strings.Join(sd, ";")))

	c.Data["json"] = models.GetDataGrid(m, 0, 0, 7)
	c.ServeJSON()
}
