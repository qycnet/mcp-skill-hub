import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { skillsApi } from '../api/client'
import { Upload, X, CheckCircle, AlertCircle, FileText } from 'lucide-react'

export default function PublishPage() {
  const navigate = useNavigate()
  const [dragOver, setDragOver] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  const [manualMode, setManualMode] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    display_name: '',
    description: '',
    category: '',
    tags: '',
    license: 'MIT',
    repository: '',
    homepage: '',
  })

  // 发布技能突变
  const publishMutation = useMutation({
    mutationFn: skillsApi.publish,
    onSuccess: (res) => {
      alert('技能发布成功！等待审核后上线')
      navigate(`/skills/${res.data.id}`)
    },
    onError: (error) => {
      alert(`发布失败：${error.response?.data?.error || error.message}`)
    },
  })

  const handleDragOver = (e) => {
    e.preventDefault()
    setDragOver(true)
  }

  const handleDragLeave = () => {
    setDragOver(false)
  }

  const handleDrop = (e) => {
    e.preventDefault()
    setDragOver(false)
    const file = e.dataTransfer.files[0]
    if (file && file.name.endsWith('.zip')) {
      setSelectedFile(file)
    } else {
      alert('请上传 ZIP 格式的技能包')
    }
  }

  const handleFileSelect = (e) => {
    const file = e.target.files[0]
    if (file && file.name.endsWith('.zip')) {
      setSelectedFile(file)
    } else {
      alert('请上传 ZIP 格式的技能包')
    }
  }

  const handleSubmit = (e) => {
    e.preventDefault()

    if (!manualMode && !selectedFile) {
      alert('请上传技能包或切换到手动填写模式')
      return
    }

    // 验证必填字段
    if (!formData.name || !formData.display_name || !formData.description || !formData.category) {
      alert('请填写所有必填字段')
      return
    }

    // 构建发布数据
    const publishData = {
      ...formData,
      tags: formData.tags.split(',').map(t => t.trim()).filter(t => t),
    }

    publishMutation.mutate(publishData)
  }

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    })
  }

  const categories = [
    { value: '', label: '选择分类', disabled: true },
    { value: 'developer-tools', label: '开发者工具' },
    { value: 'ai', label: '人工智能' },
    { value: 'security', label: '网络安全' },
    { value: 'productivity', label: '效率工具' },
    { value: 'data', label: '数据分析' },
    { value: 'automation', label: '自动化' },
    { value: 'other', label: '其他' },
  ]

  const licenses = [
    { value: 'MIT', label: 'MIT License' },
    { value: 'Apache-2.0', label: 'Apache License 2.0' },
    { value: 'GPL-3.0', label: 'GNU GPL v3' },
    { value: 'BSD-3-Clause', label: 'BSD 3-Clause' },
    { value: 'ISC', label: 'ISC License' },
    { value: 'UNLICENSED', label: '专有许可' },
  ]

  return (
    <div className="max-w-3xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          发布技能
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          分享你的 MCP 技能，帮助其他开发者提高效率
        </p>
      </div>

      <div className="card">
        {/* 上传方式选择 */}
        <div className="mb-6">
          <div className="flex gap-4">
            <button
              type="button"
              onClick={() => setManualMode(false)}
              className={`flex-1 py-3 px-4 rounded-lg border-2 transition-colors ${
                !manualMode
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 dark:border-gray-700'
              }`}
            >
              <Upload className="w-6 h-6 mx-auto mb-2" />
              <div className="font-semibold">上传技能包</div>
              <div className="text-xs text-gray-500">推荐：自动解析 manifest</div>
            </button>
            <button
              type="button"
              onClick={() => setManualMode(true)}
              className={`flex-1 py-3 px-4 rounded-lg border-2 transition-colors ${
                manualMode
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-200 dark:border-gray-700'
              }`}
            >
              <FileText className="w-6 h-6 mx-auto mb-2" />
              <div className="font-semibold">手动填写</div>
              <div className="text-xs text-gray-500">适合测试和学习</div>
            </button>
          </div>
        </div>

        {/* 上传区域 */}
        {!manualMode && (
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              技能包文件
            </label>
            <div
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
                dragOver
                  ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                  : 'border-gray-300 dark:border-gray-600'
              }`}
            >
              {selectedFile ? (
                <div className="flex items-center justify-center">
                  <CheckCircle className="w-12 h-12 text-green-500 mr-4" />
                  <div className="text-left">
                    <div className="font-semibold">{selectedFile.name}</div>
                    <div className="text-sm text-gray-500">
                      {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                    </div>
                  </div>
                  <button
                    onClick={() => setSelectedFile(null)}
                    className="ml-4 text-red-600 hover:text-red-700"
                  >
                    <X className="w-6 h-6" />
                  </button>
                </div>
              ) : (
                <div>
                  <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-700 dark:text-gray-300 mb-2">
                    拖拽 ZIP 文件到此处，或点击选择文件
                  </p>
                  <p className="text-sm text-gray-500">
                    技能包应包含 mcp-manifest.json 和必要的代码文件
                  </p>
                  <input
                    type="file"
                    accept=".zip"
                    onChange={handleFileSelect}
                    className="hidden"
                    id="file-upload"
                  />
                  <label
                    htmlFor="file-upload"
                    className="btn-primary inline-block mt-4 cursor-pointer"
                  >
                    选择文件
                  </label>
                </div>
              )}
            </div>
          </div>
        )}

        {/* 手动表单 */}
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 技能名称 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                技能名称 *
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                className="input-field"
                placeholder="my-awesome-skill"
              />
              <p className="text-xs text-gray-500 mt-1">
                只能包含小写字母、数字和连字符
              </p>
            </div>

            {/* 显示名称 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                显示名称 *
              </label>
              <input
                type="text"
                name="display_name"
                value={formData.display_name}
                onChange={handleChange}
                required
                className="input-field"
                placeholder="My Awesome Skill"
              />
            </div>
          </div>

          {/* 描述 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              描述 *
            </label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              required
              rows={4}
              className="input-field"
              placeholder="详细描述你的技能功能和用途..."
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 分类 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                分类 *
              </label>
              <select
                name="category"
                value={formData.category}
                onChange={handleChange}
                required
                className="input-field"
              >
                {categories.map((cat) => (
                  <option key={cat.value} value={cat.value} disabled={cat.disabled}>
                    {cat.label}
                  </option>
                ))}
              </select>
            </div>

            {/* 许可证 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                许可证
              </label>
              <select
                name="license"
                value={formData.license}
                onChange={handleChange}
                className="input-field"
              >
                {licenses.map((lic) => (
                  <option key={lic.value} value={lic.value}>
                    {lic.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* 标签 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              标签
            </label>
            <input
              type="text"
              name="tags"
              value={formData.tags}
              onChange={handleChange}
              className="input-field"
              placeholder="ai, productivity, automation（用逗号分隔）"
            />
          </div>

          {/* 外部链接 */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                代码仓库
              </label>
              <input
                type="url"
                name="repository"
                value={formData.repository}
                onChange={handleChange}
                className="input-field"
                placeholder="https://github.com/username/repo"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                主页
              </label>
              <input
                type="url"
                name="homepage"
                value={formData.homepage}
                onChange={handleChange}
                className="input-field"
                placeholder="https://example.com"
              />
            </div>
          </div>

          {/* 提交按钮 */}
          <div className="flex items-center gap-4">
            <button
              type="submit"
              disabled={publishMutation.isPending}
              className="btn-primary flex-1"
            >
              {publishMutation.isPending ? '发布中...' : '发布技能'}
            </button>
            <button
              type="button"
              onClick={() => navigate('/profile')}
              className="btn-secondary"
            >
              取消
            </button>
          </div>

          {/* 提示信息 */}
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <div className="flex items-start">
              <AlertCircle className="w-5 h-5 text-blue-600 dark:text-blue-400 mr-2 mt-0.5" />
              <div className="text-sm text-blue-800 dark:text-blue-300">
                <p className="font-semibold mb-1">发布须知</p>
                <ul className="list-disc list-inside space-y-1">
                  <li>技能发布后需要审核才能上线</li>
                  <li>确保你的技能包不包含恶意代码</li>
                  <li>建议使用开源许可证</li>
                  <li>提供清晰的文档和使用说明</li>
                </ul>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}
