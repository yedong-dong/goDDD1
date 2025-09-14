package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VueController Vue 控制器
type VueController struct {
}

// NewVueController 创建 Vue 控制器
func NewVueController() *VueController {
	return &VueController{}
}

// Index 首页
func (vc *VueController) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"roles":        []string{"admin"},
			"introduction": "I am a super administrator",
			"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			"name":         "Super Admin",
		},
	})
}

// TableItem 表格数据项结构
type TableItem struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Pageviews   int    `json:"pageviews"`
	Status      string `json:"status"`
	DisplayTime string `json:"display_time"`
}

func (vc *VueController) Table(c *gin.Context) {
	// 模拟表格数据
	items := []TableItem{
		{
			Title:       "Vue.js 入门教程",
			Author:      "张三",
			Pageviews:   1024,
			Status:      "published",
			DisplayTime: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			Title:       "Go语言实战指南",
			Author:      "李四",
			Pageviews:   2048,
			Status:      "published",
			DisplayTime: time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			Title:       "React 组件开发",
			Author:      "王五",
			Pageviews:   512,
			Status:      "draft",
			DisplayTime: time.Now().Add(-48 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			Title:       "数据库设计原理",
			Author:      "赵六",
			Pageviews:   3072,
			Status:      "published",
			DisplayTime: time.Now().Add(-72 * time.Hour).Format("2006-01-02 15:04:05"),
		},
		{
			Title:       "微服务架构实践",
			Author:      "孙七",
			Pageviews:   128,
			Status:      "deleted",
			DisplayTime: time.Now().Add(-96 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"items": items,
		},
	})
}
