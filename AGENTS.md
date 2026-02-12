# Opencode Skills 指南

本文档整合了 Anthropic Skills 的核心功能，可在项目根目录执行各种专业任务。

## 可用技能概览

### 1. DOCX - Word 文档处理

用于创建、读取、编辑或操作 Word 文档（.docx 文件）。

**何时使用：**
- 用户提到 "Word doc"、"word document"、".docx"
- 创建带有目录、标题、页码或信头的专业文档
- 从 .docx 文件提取或重新组织内容
- 在文档中插入或替换图片
- 在 Word 文件中执行查找和替换
- 处理修订记录或评论
- 用户请求以 Word 或 .docx 格式创建报告、备忘录、信件、模板

**核心要点：**
- .docx 文件是一个包含 XML 文件的 ZIP 归档
- **创建新文档：** 使用 `docx` npm 包
- **编辑现有文档：** 解包 → 编辑 XML → 重新打包

**快速示例：**
```javascript
const { Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell } = require('docx');
const fs = require('fs');

const doc = new Document({
  sections: [{
    properties: {
      page: {
        size: { width: 12240, height: 15840 }, // US Letter
        margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 }
      }
    },
    children: [
      new Paragraph({
        children: [new TextRun({ text: "Hello World", bold: true, size: 32 })]
      })
    ]
  }]
});

Packer.toBuffer(doc).then(buffer => {
  fs.writeFileSync("document.docx", buffer);
});
```

**关键规则：**
- **显式设置页面大小** - docx-js 默认使用 A4；美国文档使用 US Letter (12240 x 15840 DXA)
- **永远不要使用 `\n`** - 使用单独的 Paragraph 元素
- **ImageRun 需要 `type`** - 始终指定 png/jpg 等类型
- **表格需要双重宽度** - `columnWidths` 数组 AND 单元格 `width`，两者必须匹配
- **使用 `ShadingType.CLEAR`** - 永远不要用 SOLID 作为表格底纹

**完整文档：** 参见 `.opencode/skills/anthropic-skills/skills/docx/SKILL.md`

---

### 2. PDF - PDF 文档处理

用于读取、提取文本/表格、合并、拆分、旋转、添加水印、创建新 PDF、填写表单、加密/解密、提取图片、OCR 等。

**何时使用：**
- 用户想要对 PDF 文件进行任何操作
- 读取或提取文本/表格
- 合并多个 PDF
- 拆分 PDF
- 旋转页面
- 添加水印
- 创建新 PDF
- 填写 PDF 表单
- 加密/解密 PDF
- 从扫描的 PDF 进行 OCR

**核心工具：**

**Python 库：**
- **pypdf** - 基本操作（合并、拆分、旋转）
- **pdfplumber** - 文本和表格提取
- **reportlab** - 创建 PDF

**命令行工具：**
- **pdftotext** (poppler-utils) - 提取文本
- **qpdf** - 合并、拆分、旋转、解密

**快速示例：**

```python
# 合并 PDF
from pypdf import PdfWriter, PdfReader

writer = PdfWriter()
for pdf_file in ["doc1.pdf", "doc2.pdf"]:
    reader = PdfReader(pdf_file)
    for page in reader.pages:
        writer.add_page(page)

with open("merged.pdf", "wb") as output:
    writer.write(output)

# 提取文本
import pdfplumber
with pdfplumber.open("document.pdf") as pdf:
    for page in pdf.pages:
        text = page.extract_text()
        print(text)

# 创建 PDF
from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas

c = canvas.Canvas("hello.pdf", pagesize=letter)
c.drawString(100, 700, "Hello World!")
c.save()
```

**表格提取：**
```python
import pdfplumber
import pandas as pd

with pdfplumber.open("document.pdf") as pdf:
    for page in pdf.pages:
        tables = page.extract_tables()
        for table in tables:
            df = pd.DataFrame(table[1:], columns=table[0])
            print(df)
```

**扫描 PDF 的 OCR：**
```python
import pytesseract
from pdf2image import convert_from_path

images = convert_from_path('scanned.pdf')
for i, image in enumerate(images):
    text = pytesseract.image_to_string(image)
    print(f"Page {i+1}:\n{text}\n")
```

**完整文档：** 参见 `.opencode/skills/anthropic-skills/skills/pdf/SKILL.md`

---

### 3. XLSX - Excel 电子表格处理

用于打开、读取、编辑或修复现有的 .xlsx、.xlsm、.csv 或 .tsv 文件；从头开始创建新电子表格；或在表格文件格式之间转换。

**何时使用：**
- 添加列、计算公式、格式化、图表、清理杂乱数据
- 从其他数据源创建新电子表格
- 在表格文件格式之间转换
- 用户引用电子表格文件名或路径

**核心工具：**
- **pandas** - 数据分析、批量操作、简单数据导出
- **openpyxl** - 复杂格式化、公式、Excel 特定功能

**关键要求：**

⚠️ **使用公式，而非硬编码值**
```python
# ❌ 错误 - 硬编码计算值
sheet['B10'] = 5000

# ✅ 正确 - 使用 Excel 公式
sheet['B10'] = '=SUM(B2:B9)'
```

**快速示例：**

```python
# 读取 Excel
import pandas as pd
df = pd.read_excel('file.xlsx')
df.head()
df.info()

# 创建新 Excel
from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment

wb = Workbook()
sheet = wb.active
sheet['A1'] = 'Hello'
sheet['B1'] = 'World'
sheet['B2'] = '=SUM(A1:A10)'

# 格式化
sheet['A1'].font = Font(bold=True, color='FF0000')
sheet['A1'].alignment = Alignment(horizontal='center')
sheet.column_dimensions['A'].width = 20

wb.save('output.xlsx')
```

**编辑现有文件：**
```python
from openpyxl import load_workbook

wb = load_workbook('existing.xlsx')
sheet = wb['SheetName']
sheet['A1'] = 'New Value'
sheet.insert_rows(2)
wb.save('modified.xlsx')
```

**财务模型颜色标准：**
- **蓝色文本 (RGB: 0,0,255)**：硬编码输入，用户会更改的场景数字
- **黑色文本 (RGB: 0,0,0)**：所有公式和计算
- **绿色文本 (RGB: 0,128,0)**：从同一工作簿其他工作表提取的链接
- **红色文本 (RGB: 255,0,0)**：指向其他文件的外部链接
- **黄色背景 (RGB: 255,255,0)**：需要注意的关键假设或需要更新的单元格

**数字格式标准：**
- **年份**：格式化为文本字符串（例如 "2024" 而非 "2,024"）
- **货币**：使用 $#,##0 格式；始终在标题中指定单位（"Revenue ($mm)"）
- **零**：使用数字格式使所有零显示为 "-"
- **百分比**：默认为 0.0% 格式（一位小数）

**完整文档：** 参见 `.opencode/skills/anthropic-skills/skills/xlsx/SKILL.md`

---

### 4. Webapp Testing - Web 应用测试

用于与本地 Web 应用交互和测试，使用 Playwright。支持验证前端功能、调试 UI 行为、捕获浏览器截图和查看浏览器日志。

**何时使用：**
- 测试本地 Web 应用
- 验证前端功能
- 调试 UI 行为
- 捕获浏览器截图
- 查看浏览器日志

**决策树：**
```
用户任务 → 是静态 HTML 吗？
    ├─ 是 → 直接读取 HTML 文件以识别选择器
    │         ├─ 成功 → 使用选择器编写 Playwright 脚本
    │         └─ 失败/不完整 → 视为动态（见下文）
    │
    └─ 否（动态 Web 应用）→ 服务器已经在运行了吗？
        ├─ 否 → 需要先启动服务器
        │
        └─ 是 → 侦察-然后-行动：
            1. 导航并等待 networkidle
            2. 截图或检查 DOM
            3. 从渲染状态识别选择器
            4. 使用发现的选择器执行操作
```

**快速示例：**

```python
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page()
    page.goto('http://localhost:5173')
    page.wait_for_load_state('networkidle')  # 关键：等待 JS 执行
    
    # 截图检查
    page.screenshot(path='/tmp/screenshot.png', full_page=True)
    
    # 与元素交互
    page.click('button[type="submit"]')
    page.fill('input[name="username"]', 'testuser')
    
    browser.close()
```

**可用脚本：**
- `scripts/with_server.py` - 管理服务器生命周期（支持多个服务器）

**使用 with_server.py：**
```bash
# 单服务器
python scripts/with_server.py --server "npm run dev" --port 5173 -- python your_automation.py

# 多服务器（例如后端 + 前端）
python scripts/with_server.py \
  --server "cd backend && python server.py" --port 3000 \
  --server "cd frontend && npm run dev" --port 5173 \
  -- python your_automation.py
```

**最佳实践：**
- 始终在检查动态应用的 DOM 之前等待 `networkidle`
- 使用描述性选择器：`text=`、`role=`、CSS 选择器或 ID
- 添加适当的等待：`page.wait_for_selector()` 或 `page.wait_for_timeout()`
- 完成后始终关闭浏览器

**完整文档：** 参见 `.opencode/skills/anthropic-skills/skills/webapp-testing/SKILL.md`

---

### 5. MCP Builder - MCP 服务器生成器

用于生成 Model Context Protocol (MCP) 服务器。这是一个创建自定义工具集成的方式。

**何时使用：**
- 用户想要创建自定义工具集成
- 构建 MCP 服务器以连接外部 API
- 为特定工作流程创建工具

**完整文档：** 参见 `.opencode/skills/anthropic-skills/skills/mcp-builder/SKILL.md`

---

### 6. 其他技能

**技能列表：**
- **algorithmic-art** - 算法艺术生成
- **brand-guidelines** - 品牌指南
- **canvas-design** - 画布设计
- **doc-coauthoring** - 文档协作
- **frontend-design** - 前端设计
- **internal-comms** - 内部沟通
- **skill-creator** - 技能创建辅助
- **slack-gif-creator** - Slack GIF 创建器
- **theme-factory** - 主题工厂
- **web-artifacts-builder** - Web 工件构建器

**完整技能文档位置：** `.opencode/skills/anthropic-skills/skills/`

---

## 如何使用这些技能

1. **识别用户需求** - 确定需要哪个技能（DOCX、PDF、XLSX 等）

2. **参考相关技能部分** - 阅读上述相应部分的关键要点

3. **查阅完整文档** - 对于复杂任务，参考完整的 SKILL.md 文件：
   ```bash
   cat .opencode/skills/anthropic-skills/skills/<skill-name>/SKILL.md
   ```

4. **执行最佳实践** - 遵循每个技能中概述的具体指南和代码示例

5. **验证结果** - 确保输出符合技能要求（例如，XLSX 文件无公式错误）

---

## 技能来源

这些技能改编自 Anthropic 的公开技能仓库：
- **仓库：** https://github.com/anthropics/skills
- **许可证：** Apache 2.0（示例技能），部分为专有（文档技能）
- **用途：** 演示和教育目的

**免责声明：** 这些技能仅供演示和教育目的。虽然某些功能可能在 Claude 中可用，但你从此处实现获得的行为可能与技能中显示的不同。在依赖关键任务之前，请始终在你自己的环境中彻底测试技能。
