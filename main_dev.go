//go:build dev
// +build dev

/*
 * @Description: 开发模式入口
 * @Author: 安知鱼
 * @Date: 2025-01-23
 */
package main

import (
	"embed"
	"flag"
	"log"
	"os"

	"github.com/anzhiyu-c/anheyu-app/cmd/server"
)

// 开发模式下使用空的 embed.FS
var content embed.FS

func main() {
	// 解析命令行参数
	var exportAssetsDir string
	flag.StringVar(&exportAssetsDir, "export-assets", "", "导出静态资源到指定目录（用于自定义静态资源）")
	flag.Parse()

	// 开发模式不支持导出资源
	if exportAssetsDir != "" {
		log.Println("⚠️  开发模式不支持导出静态资源，请先构建前端后使用生产模式")
		return
	}

	log.Println("🔧 开发模式启动 - 前端请单独运行 npm run serve")
	log.Println("💡 提示：前端开发服务器通常运行在 http://localhost:5173 或 http://localhost:8080")

	// 调用位于 cmd/server 包中的 NewApp 函数来构建整个应用
	// 开发模式下传入空的 embed.FS
	app, cleanup, err := server.NewApp(content)
	if err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}

	// 使用 defer 来确保 cleanup 函数在 main 退出时被调用
	defer cleanup()

	// 确保后台任务在程序退出时被停止
	defer app.Stop()

	if os.Getenv("ANHEYU_LICENSE_KEY") == "" {
		app.PrintBanner()
	}

	// 启动应用（只提供 API 服务）
	if err := app.Run(); err != nil {
		log.Fatalf("应用运行失败: %v", err)
	}
}
