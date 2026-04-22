import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { skillsApi } from '../api/client'
import { Star, Download, Filter } from 'lucide-react'

export default function SkillListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [showFilters, setShowFilters] = useState(false)

  const query = searchParams.get('q') || ''
  const category = searchParams.get('category') || ''
  const sort = searchParams.get('sort') || 'downloads'
  const page = parseInt(searchParams.get('page') || '1')

  const { data, isLoading } = useQuery({
    queryKey: ['skills', { q: query, category, sort, page }],
    queryFn: () => skillsApi.list({
      q: query,
      category,
      sort,
      page,
      page_size: 20,
    }),
    select: (res) => res.data,
  })

  const handleFilterChange = (key, value) => {
    const newParams = new URLSearchParams(searchParams)
    if (value) {
      newParams.set(key, value)
    } else {
      newParams.delete(key)
    }
    newParams.set('page', '1') // 重置页码
    setSearchParams(newParams)
  }

  const sortOptions = [
    { value: 'downloads', label: '下载量' },
    { value: 'rating', label: '评分' },
    { value: 'quality_score', label: '质量评分' },
    { value: 'created_at', label: '最新发布' },
  ]

  const categories = [
    { value: '', label: '全部分类' },
    { value: 'developer-tools', label: '开发者工具' },
    { value: 'ai', label: '人工智能' },
    { value: 'security', label: '网络安全' },
    { value: 'productivity', label: '效率工具' },
    { value: 'data', label: '数据分析' },
  ]

  return (
    <div className="space-y-6">
      {/* 页面标题 */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            技能库
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            {query ? `搜索 "${query}"` : '发现和安装 MCP 技能'}
            {data && ` - 共 ${data.total} 个技能`}
          </p>
        </div>

        <button
          onClick={() => setShowFilters(!showFilters)}
          className="flex items-center space-x-2 px-4 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
        >
          <Filter className="w-5 h-5" />
          <span>筛选</span>
        </button>
      </div>

      {/* 筛选器 */}
      {showFilters && (
        <div className="card space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* 分类筛选 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                分类
              </label>
              <select
                value={category}
                onChange={(e) => handleFilterChange('category', e.target.value)}
                className="input-field"
              >
                {categories.map((cat) => (
                  <option key={cat.value} value={cat.value}>
                    {cat.label}
                  </option>
                ))}
              </select>
            </div>

            {/* 排序 */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                排序
              </label>
              <select
                value={sort}
                onChange={(e) => handleFilterChange('sort', e.target.value)}
                className="input-field"
              >
                {sortOptions.map((opt) => (
                  <option key={opt.value} value={opt.value}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
      )}

      {/* 技能列表 */}
      {isLoading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-600 dark:text-gray-400">加载中...</p>
        </div>
      ) : data?.skills?.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-xl text-gray-600 dark:text-gray-400">未找到相关技能</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {data?.skills?.map((skill) => (
              <SkillListItem key={skill.id} skill={skill} />
            ))}
          </div>

          {/* 分页 */}
          {data && data.total > 20 && (
            <div className="flex justify-center space-x-2">
              {Array.from({ length: Math.ceil(data.total / 20) }, (_, i) => i + 1).map((pageNum) => (
                <button
                  key={pageNum}
                  onClick={() => handleFilterChange('page', String(pageNum))}
                  className={`px-4 py-2 rounded-lg ${
                    pageNum === page
                      ? 'bg-primary-600 text-white'
                      : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  }`}
                >
                  {pageNum}
                </button>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}

function SkillListItem({ skill }) {
  return (
    <a
      href={`/skills/${skill.id}`}
      className="card block hover:shadow-xl transition-shadow"
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          {skill.display_name}
        </h3>
        <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 text-xs rounded">
          {skill.category}
        </span>
      </div>

      <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">
        {skill.description}
      </p>

      <div className="flex items-center justify-between text-sm">
        <div className="flex items-center space-x-3 text-gray-500">
          <span className="flex items-center">
            <Star className="w-4 h-4 mr-1 text-yellow-500" />
            {skill.rating.toFixed(1)}
          </span>
          <span className="flex items-center">
            <Download className="w-4 h-4 mr-1" />
            {skill.downloads}
          </span>
        </div>
        <span className="text-xs text-primary-600 dark:text-primary-400 font-medium">
          评分：{skill.quality_score.toFixed(0)}
        </span>
      </div>
    </a>
  )
}
