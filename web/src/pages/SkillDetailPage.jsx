import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { skillsApi } from '../api/client'
import { useAuthStore } from '../stores/authStore'
import { 
  Download, Star, Calendar, User, Tag, 
  ExternalLink, Gitlab, FileText, ChevronRight,
  Send, ThumbsUp
} from 'lucide-react'

export default function SkillDetailPage() {
  const { id } = useParams()
  const { isAuthenticated } = useAuthStore()
  const [rating, setRating] = useState(0)
  const [comment, setComment] = useState('')

  // 获取技能详情
  const { data: skill, isLoading } = useQuery({
    queryKey: ['skill', id],
    queryFn: () => skillsApi.getById(id),
    select: (res) => res.data,
  })

  // 评分突变
  const rateMutation = useMutation({
    mutationFn: ({ id, rating, comment }) => skillsApi.rate(id, rating, comment),
    onSuccess: () => {
      alert('评分成功！')
      setRating(0)
      setComment('')
    },
  })

  const handleSubmitRating = (e) => {
    e.preventDefault()
    if (!isAuthenticated) {
      alert('请先登录')
      return
    }
    if (rating === 0) {
      alert('请选择评分')
      return
    }
    rateMutation.mutate({ id, rating, comment })
  }

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
        <p className="mt-4 text-gray-600 dark:text-gray-400">加载中...</p>
      </div>
    )
  }

  if (!skill) {
    return (
      <div className="text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">技能不存在</h1>
        <Link to="/skills" className="btn-primary">返回技能库</Link>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* 面包屑导航 */}
      <nav className="flex items-center text-sm text-gray-600 dark:text-gray-400">
        <Link to="/" className="hover:text-primary-600">首页</Link>
        <ChevronRight className="w-4 h-4 mx-2" />
        <Link to="/skills" className="hover:text-primary-600">技能库</Link>
        <ChevronRight className="w-4 h-4 mx-2" />
        <span className="text-gray-900 dark:text-white">{skill.display_name}</span>
      </nav>

      {/* 技能头部信息 */}
      <div className="card">
        <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
                {skill.display_name}
              </h1>
              {skill.is_verified && (
                <span className="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 text-xs rounded-full">
                  已验证
                </span>
              )}
            </div>
            <p className="text-gray-500 mb-4">{skill.name}</p>
            <p className="text-gray-700 dark:text-gray-300 mb-6">{skill.description}</p>

            {/* 元信息 */}
            <div className="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-400">
              <span className="flex items-center">
                <Tag className="w-4 h-4 mr-1" />
                {skill.category}
              </span>
              <span className="flex items-center">
                <Star className="w-4 h-4 mr-1 text-yellow-500" />
                {skill.rating.toFixed(1)} ({skill.rating_count} 个评分)
              </span>
              <span className="flex items-center">
                <Download className="w-4 h-4 mr-1" />
                {skill.downloads} 次下载
              </span>
              <span className="flex items-center">
                <Calendar className="w-4 h-4 mr-1" />
                {new Date(skill.created_at).toLocaleDateString('zh-CN')}
              </span>
            </div>
          </div>

          {/* 操作按钮 */}
          <div className="flex flex-col gap-3">
            <button className="btn-primary flex items-center justify-center">
              <Download className="w-5 h-5 mr-2" />
              下载安装
            </button>
            {skill.repository && (
              <a
                href={skill.repository}
                target="_blank"
                rel="noopener noreferrer"
                className="btn-secondary flex items-center justify-center"
              >
                <Gitlab className="w-5 h-5 mr-2" />
                源代码
              </a>
            )}
            {skill.homepage && (
              <a
                href={skill.homepage}
                target="_blank"
                rel="noopener noreferrer"
                className="btn-secondary flex items-center justify-center"
              >
                <ExternalLink className="w-5 h-5 mr-2" />
                主页
              </a>
            )}
          </div>
        </div>
      </div>

      {/* 标签 */}
      {skill.tags && skill.tags.length > 0 && (
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 flex items-center">
            <Tag className="w-5 h-5 mr-2" />
            标签
          </h2>
          <div className="flex flex-wrap gap-2">
            {skill.tags.map((tag) => (
              <span
                key={tag.id}
                className="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-full text-sm"
              >
                {tag.name}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* 版本历史 */}
      {skill.versions && skill.versions.length > 0 && (
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 flex items-center">
            <FileText className="w-5 h-5 mr-2" />
            版本历史
          </h2>
          <div className="space-y-3">
            {skill.versions.map((version) => (
              <div
                key={version.id}
                className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-mono font-semibold text-primary-600">
                      v{version.version}
                    </span>
                    {version.is_latest && (
                      <span className="px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 text-xs rounded">
                        最新
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {version.description}
                  </p>
                </div>
                <div className="text-sm text-gray-500">
                  {new Date(version.created_at).toLocaleDateString('zh-CN')}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 作者信息 */}
      {skill.author && (
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 flex items-center">
            <User className="w-5 h-5 mr-2" />
            作者
          </h2>
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center">
              <span className="text-primary-600 dark:text-primary-400 font-semibold text-lg">
                {skill.author.username[0].toUpperCase()}
              </span>
            </div>
            <div>
              <div className="font-semibold text-gray-900 dark:text-white">
                {skill.author.username}
              </div>
              {skill.author.bio && (
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  {skill.author.bio}
                </p>
              )}
            </div>
          </div>
        </div>
      )}

      {/* 评分和评论 */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4 flex items-center">
          <Star className="w-5 h-5 mr-2 text-yellow-500" />
          评分和评论
        </h2>

        {/* 评分表单 */}
        <form onSubmit={handleSubmitRating} className="mb-6">
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              你的评分
            </label>
            <div className="flex gap-2">
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  onClick={() => setRating(star)}
                  className={`text-2xl transition-colors ${
                    star <= rating
                      ? 'text-yellow-500'
                      : 'text-gray-300 hover:text-yellow-400'
                  }`}
                >
                  <Star className="w-8 h-8 fill-current" />
                </button>
              ))}
            </div>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              评论（可选）
            </label>
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              rows={3}
              className="input-field"
              placeholder="分享你的使用体验..."
            />
          </div>
          <button
            type="submit"
            disabled={rateMutation.isPending}
            className="btn-primary flex items-center"
          >
            <Send className="w-4 h-4 mr-2" />
            {rateMutation.isPending ? '提交中...' : '提交评分'}
          </button>
        </form>

        {/* 评论列表 */}
        <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
          <h3 className="font-semibold mb-4">用户评论</h3>
          <div className="space-y-4">
            {/* 示例评论 */}
            <div className="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-8 h-8 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center">
                  <span className="text-sm font-semibold">U</span>
                </div>
                <span className="font-medium">User123</span>
                <div className="flex text-yellow-500">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <Star key={i} className="w-4 h-4 fill-current" />
                  ))}
                </div>
              </div>
              <p className="text-gray-700 dark:text-gray-300">
                非常好用的技能，大大提高了我的工作效率！
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
