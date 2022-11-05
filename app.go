package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	"packencryptor/encryptor"
	runtime2 "runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectPathDialog() string {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Minecraft Bedrock Resource Pack (*.mcpack, *.zip)",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return filePath
}

func (a *App) SelectDestDialog() string {
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Minecraft Bedrock Resource Pack (*.mcpack, *.zip)",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return filePath
}

func (a *App) EncryptPack(path string, dest string) (key string, e error) {
	defer func() {
		if r, ok := recover().(error); ok {
			key = ""
			e = r
		}
	}()
	key = encryptor.RandomKey()
	encryptor.EncryptPack(path, dest, key)

	if runtime2.GOOS == "windows" {
		cmd := exec.Command("explorer", "/select,", dest)
		_ = cmd.Run()
	}

	return key, nil
}
