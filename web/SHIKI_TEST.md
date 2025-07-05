# Shiki + Monaco Editor 集成测试指南

## 测试环境

- 开发服务器已启动：`http://localhost:5173/`
- 依赖已更新为 ESM 兼容版本
- Monaco Editor Worker 配置已完成

## 测试步骤

### 1. 访问模板编辑页面

1. 打开浏览器访问 `http://localhost:5173/`
2. 导航到模板列表页面
3. 点击任意模板的编辑按钮，进入模板编辑页面

### 2. 验证 Monaco Editor 加载

1. 在文件树中点击任意文件
2. 检查 Monaco Editor 是否正常显示
3. 确认编辑器有深色主题背景

### 3. 验证语法高亮

测试以下文件类型的语法高亮：

#### 前端文件
- `.vue` - Vue 单文件组件
- `.js` - JavaScript 文件
- `.ts` - TypeScript 文件
- `.json` - JSON 配置文件
- `.html` - HTML 文件
- `.md` - Markdown 文件

#### 后端文件
- `.go` - Go 语言文件
- `.java` - Java 文件
- `.py` - Python 文件
- `.c` - C 语言文件
- `.cpp` - C++ 文件
- `.cs` - C# 文件
- `.php` - PHP 文件
- `.rb` - Ruby 文件

#### 配置文件
- `.sh` - Shell 脚本
- `.xml` - XML 文件
- `.yml` / `.yaml` - YAML 配置文件

### 4. 验证多标签功能

1. 打开多个不同类型的文件
2. 检查标签页是否正确显示文件名
3. 切换标签页，确认语法高亮正确切换
4. 关闭标签页，确认功能正常

### 5. 验证保存功能

1. 在编辑器中修改文件内容
2. 使用 Ctrl+S 快捷键保存
3. 检查通知提示是否显示
4. 右键菜单中是否显示"保存"选项

### 6. 验证深色主题

1. 检查编辑器背景是否为深色
2. 确认标签页样式与编辑器主题一致
3. 验证整体界面无缝衔接

## 预期结果

### ✅ 正常情况

- Monaco Editor 正常加载，无控制台错误
- 语法高亮正确显示，颜色丰富
- 多标签页功能正常
- 保存功能正常工作
- 深色主题统一美观

### ❌ 可能的问题

1. **Worker 加载警告**
   - 检查 `public/monaco/` 目录下的 worker 文件是否存在
   - 确认 `main.js` 中的 `MonacoEnvironment` 配置正确

2. **语法高亮不生效**
   - 检查 `tab.name` 是否为完整文件名（包含扩展名）
   - 确认 `getLanguage` 函数正确识别文件类型
   - 验证 shiki 高亮器初始化成功

3. **主题显示异常**
   - 确认 `monaco.editor.setTheme('github-dark')` 已调用
   - 检查 CSS 样式是否正确应用

## 调试技巧

### 控制台调试

在浏览器控制台中检查：

```javascript
// 检查 Monaco 是否加载
console.log(window.monaco)

// 检查当前主题
console.log(monaco.editor.getTheme())

// 检查语言识别
console.log('当前文件语言:', lang)
```

### 网络调试

检查 Network 标签页：
- Monaco Editor worker 文件是否正确加载
- 是否有 404 错误

## 常见问题解决

### 1. ESM 兼容性问题

如果遇到 ESM 相关错误：
- 确保 `package.json` 中 `"type": "module"`
- 使用 `.js` 扩展名的配置文件
- 使用 `import/export` 语法

### 2. 依赖版本问题

如果遇到依赖冲突：
```bash
npm install
npm audit fix
```

### 3. Worker 路径问题

如果 worker 文件加载失败：
- 确认 `public/monaco/` 目录存在
- 检查文件路径配置
- 验证静态文件服务正常

## 性能优化建议

1. **按需加载语言**
   - 只加载需要的语言支持
   - 减少初始包大小

2. **缓存优化**
   - 利用浏览器缓存 worker 文件
   - 避免重复加载

3. **内存管理**
   - 及时销毁不需要的编辑器实例
   - 清理事件监听器 