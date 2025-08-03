package interactive

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/ciclebyte/template_starter/cli/internal/client"
)

// PreviewFiles 显示文件预览界面（简化版本，兼容Windows）
func PreviewFiles(files []client.RenderedFile, projectName string) (bool, error) {
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("Project Preview: %s\n", projectName)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	// 按路径排序文件
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	for {
		// 显示文件列表
		fmt.Println("Files to be generated:")
		fmt.Println(strings.Repeat("-", 40))
		
		fileMap := make(map[int]*client.RenderedFile)
		index := 1
		
		for i := range files {
			file := &files[i]
			prefix := "[FILE]"
			if file.IsDirectory {
				prefix = "[DIR] "
			}
			
			fmt.Printf("%2d. %s %s\n", index, prefix, file.Path)
			fileMap[index] = file
			index++
		}
		
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("Total: %d files/directories\n\n", len(files))

		// 提供选项
		fmt.Println("Options:")
		fmt.Println("  [number] - View file content")
		fmt.Println("  y - Confirm and generate project")
		fmt.Println("  n - Cancel")
		fmt.Println("  q - Quit")
		fmt.Print("\nEnter your choice: ")

		// 读取用户输入
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return true, nil
		case "n", "no", "q", "quit":
			return false, nil
		default:
			// 尝试解析为文件编号
			if num, err := strconv.Atoi(input); err == nil {
				if file, exists := fileMap[num]; exists {
					showFileContent(*file)
				} else {
					fmt.Printf("\nInvalid file number: %d\n\n", num)
				}
			} else {
				fmt.Printf("\nInvalid option: %s\n\n", input)
			}
		}
	}
}

// showFileContent 显示文件内容
func showFileContent(file client.RenderedFile) {
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("File: %s\n", file.Path)
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	
	if file.IsDirectory {
		fmt.Println("This is a directory - no content to display.")
	} else if file.Content == "" {
		fmt.Println("This file is empty.")
	} else {
		fmt.Println(file.Content)
	}
	
	fmt.Printf("\n" + strings.Repeat("-", 60) + "\n")
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
	fmt.Println()
}

// PreviewFilesPrompt 使用promptui的文件预览（备选方案）
func PreviewFilesPrompt(files []client.RenderedFile, projectName string) (bool, error) {
	// 构建选项列表
	items := []string{
		"View file list",
		"Browse files",
		"Confirm and generate",
		"Cancel",
	}

	for {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Project Preview: %s", projectName),
			Items: items,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return false, err
		}

		switch result {
		case "View file list":
			showFileList(files)
		case "Browse files":
			if err := browseFiles(files); err != nil {
				return false, err
			}
		case "Confirm and generate":
			return true, nil
		case "Cancel":
			return false, nil
		}
	}
}

// showFileList 显示文件列表
func showFileList(files []client.RenderedFile) {
	fmt.Printf("\n" + strings.Repeat("=", 50) + "\n")
	fmt.Println("Files to be generated:")
	fmt.Printf(strings.Repeat("-", 50) + "\n")
	
	// 按路径排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	
	for _, file := range files {
		prefix := "[FILE]"
		size := fmt.Sprintf("(%d bytes)", len(file.Content))
		if file.IsDirectory {
			prefix = "[DIR] "
			size = ""
		}
		
		fmt.Printf("%s %-40s %s\n", prefix, file.Path, size)
	}
	
	fmt.Printf(strings.Repeat("-", 50) + "\n")
	fmt.Printf("Total: %d files/directories\n", len(files))
	fmt.Print("\nPress Enter to continue...")
	fmt.Scanln()
}

// browseFiles 浏览文件内容
func browseFiles(files []client.RenderedFile) error {
	// 创建文件选择列表
	var fileItems []string
	fileMap := make(map[string]*client.RenderedFile)
	
	// 按路径排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	
	for i := range files {
		file := &files[i]
		prefix := "[FILE]"
		if file.IsDirectory {
			prefix = "[DIR] "
		}
		
		label := fmt.Sprintf("%s %s", prefix, file.Path)
		fileItems = append(fileItems, label)
		fileMap[label] = file
	}
	
	fileItems = append(fileItems, "← Back to main menu")

	for {
		prompt := promptui.Select{
			Label: "Select a file to view",
			Items: fileItems,
			Size:  15,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		if result == "← Back to main menu" {
			break
		}

		if file, exists := fileMap[result]; exists {
			showFileContent(*file)
		}
	}
	
	return nil
}