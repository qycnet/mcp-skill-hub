import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import SkillListPage from './pages/SkillListPage'
import SkillDetailPage from './pages/SkillDetailPage'
import LoginPage from './pages/LoginPage'
import PublishPage from './pages/PublishPage'
import ProfilePage from './pages/ProfilePage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="skills" element={<SkillListPage />} />
        <Route path="skills/:id" element={<SkillDetailPage />} />
        <Route path="publish" element={<PublishPage />} />
        <Route path="profile" element={<ProfilePage />} />
        <Route path="login" element={<LoginPage />} />
      </Route>
    </Routes>
  )
}

export default App
