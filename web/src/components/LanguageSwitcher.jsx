import { useTranslation } from 'react-i18next'
import { Globe } from 'lucide-react'

export default function LanguageSwitcher() {
  const { i18n, t } = useTranslation()

  const languages = [
    { code: 'zh', name: '中文' },
    { code: 'en', name: 'English' },
    { code: 'ja', name: '日本語' },
    { code: 'es', name: 'Español' },
    { code: 'fr', name: 'Français' }
  ]

  const changeLanguage = (code) => {
    i18n.changeLanguage(code)
    localStorage.setItem('i18nextLng', code)
  }

  return (
    <div className="relative group">
      <button className="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
        <Globe className="w-5 h-5" />
        <span className="text-sm font-medium">
          {languages.find(l => l.code === i18n.language)?.name || 'English'}
        </span>
      </button>

      {/* 下拉菜单 */}
      <div className="absolute right-0 mt-2 w-40 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
        {languages.map((lang) => (
          <button
            key={lang.code}
            onClick={() => changeLanguage(lang.code)}
            className={`w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700 first:rounded-t-lg last:rounded-b-lg ${
              i18n.language === lang.code ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600' : ''
            }`}
          >
            {lang.name}
          </button>
        ))}
      </div>
    </div>
  )
}
