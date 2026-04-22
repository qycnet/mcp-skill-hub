import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { skillsApi } from '../api/client'
import { Search, TrendingUp, Star, Download, ArrowRight } from 'lucide-react'

export default function HomePage() {
  const [searchQuery, setSearchQuery] = useState('')

  // 获取热门技能
  const { data: trendingSkills } = useQuery({
    queryKey: ['skills', 'trending'],
    queryFn: () => skillsApi.list({ sort: 'downloads', page_size: 6 }),
    select: (res) => res.data.skills,
  })

  // 获取高评分技能
  const { data: topRatedSkills } = useQuery({
    queryKey: ['skills', 'top-rated'],
    queryFn: () => skillsApi.list({ sort: 'rating', page_size: 6 }),
    select: (res) => res.data.skills,
  })

  const handleSearch = (e) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      window.location.href = `/skills?q=${encodeURIComponent(searchQuery)}`
    }
  }

  return (
    <div className="space-y-12">
      {/* Hero 区域 */}
      <section className="text-center py-16 bg-gradient-to-r from-primary-600 to-primary-700 rounded-2xl text-white">
        <h1 className="text-4xl md:text-5xl font-bold mb-4">
          发现和分享 MCP 技能
        </h1>
        <p className="text-xl mb-8 text-primary-100">
          让 AI Agent 更强大的技能市场
        </p>

        {/* 搜索框 */}
        <form onSubmit={handleSearch} className="max-w-2xl mx-auto">
          <div className="flex">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="搜索技能..."
              className="flex-1 px-6 py-4 rounded-l-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-white"
            />
            <button
              type="submit"
              className="bg-white text-primary-600 px-8 py-4 rounded-r-lg font-semibold hover:bg-primary-50 transition-colors"
            >
              <Search className="w-6 h-6" />
            </button>
          </div>
        </form>

        {/* 统计信息 */}
        <div className="flex justify-center space-x-8 mt-8 text-primary-100">
          <div>
            <div className="text-3xl font-bold text-white">1000+</div>
            <div className="text-sm">技能</div>
          </div>
          <div>
            <div className="text-3xl font-bold text-white">5000+</div>
            <div className="text-sm">开发者</div>
          </div>
          <div>
            <div className="text-3xl font-bold text-white">100K+</div>
            <div className="text-sm">下载量</div>
          </div>
        </div>
      </section>

      {/* 热门技能 */}
      <section>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
            <TrendingUp className="w-6 h-6 mr-2 text-primary-600" />
            热门技能
          </h2>
          <Link
            to="/skills?sort=downloads"
            className="text-primary-600 hover:text-primary-700 flex items-center"
          >
            查看全部 <ArrowRight className="w-4 h-4 ml-1" />
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {trendingSkills?.map((skill) => (
            <SkillCard key={skill.id} skill={skill} />
          ))}
        </div>
      </section>

      {/* 高评分技能 */}
      <section>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
            <Star className="w-6 h-6 mr-2 text-yellow-500" />
            高评分技能
          </h2>
          <Link
            to="/skills?sort=rating"
            className="text-primary-600 hover:text-primary-700 flex items-center"
          >
            查看全部 <ArrowRight className="w-4 h-4 ml-1" />
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {topRatedSkills?.map((skill) => (
            <SkillCard key={skill.id} skill={skill} />
          ))}
        </div>
      </section>

      {/* CTA 区域 */}
      <section className="bg-gray-100 dark:bg-gray-800 rounded-2xl p-12 text-center">
        <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          准备好分享你的技能了吗？
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-8">
          加入数千名开发者，发布你的 MCP 技能，帮助他人提高效率
        </p>
        <Link to="/publish" className="btn-primary inline-flex items-center">
          发布技能 <ArrowRight className="w-5 h-5 ml-2" />
        </Link>
      </section>
    </div>
  )
}

// 技能卡片组件
function SkillCard({ skill }) {
  return (
    <Link
      to={`/skills/${skill.id}`}
      className="card block hover:shadow-xl transition-shadow"
    >
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {skill.display_name}
          </h3>
          <p className="text-sm text-gray-500">{skill.name}</p>
        </div>
        <span className="px-3 py-1 bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 text-xs rounded-full">
          {skill.category}
        </span>
      </div>

      <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">
        {skill.description}
      </p>

      <div className="flex items-center justify-between text-sm text-gray-500">
        <div className="flex items-center space-x-4">
          <span className="flex items-center">
            <Star className="w-4 h-4 mr-1 text-yellow-500" />
            {skill.rating.toFixed(1)}
          </span>
          <span className="flex items-center">
            <Download className="w-4 h-4 mr-1" />
            {skill.downloads}
          </span>
        </div>
        <span className="text-xs">
          评分：{skill.quality_score.toFixed(0)}
        </span>
      </div>
    </Link>
  )
}
