import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom'
import { AuthProvider, useAuth } from "@/contexts/auth"
import Dashboard from "@/pages/dashboard"
import ControlsPage from './pages/controls'
import Login from "@/pages/login"
import Register from "@/pages/register"
import { Menubar, MenubarMenu, MenubarTrigger } from "@/components/ui/menubar"
import './App.css'

// Protected Route component
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, isAuthInitialized } = useAuth();

  if (!isAuthInitialized) {
    return <div>Loading...</div>; // Optional: loading indicator
  }

  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
}

function NavigationBar() {
  return (
    <Menubar className="mb-4">
      <MenubarMenu>
        <MenubarTrigger>
          <Link to="/">Dashboard</Link>
        </MenubarTrigger>
      </MenubarMenu>
      <MenubarMenu>
        <MenubarTrigger>
          <Link to="/controls">Controls</Link>
        </MenubarTrigger>
      </MenubarMenu>
    </Menubar>
  )
}

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="container mx-auto p-4">
          <NavigationBar />
          <Routes>
            <Route path="/" element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } />

            <Route path="/controls" element={
              <ProtectedRoute>
                <ControlsPage />
              </ProtectedRoute>
            } />

            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  )
}

export default App
