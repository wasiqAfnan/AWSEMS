import { loginUrl } from '../services/api'

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-700">
      <div className="bg-white rounded-2xl shadow-xl p-10 w-full max-w-sm text-center">
        <h1 className="text-3xl font-bold text-slate-800 mb-2">AWSEMS</h1>
        <p className="text-slate-500 text-sm mb-8">Employee Management System</p>
        <a
          href={loginUrl}
          className="block w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 rounded-lg transition-colors"
        >
          Sign in with Cognito
        </a>
      </div>
    </div>
  )
}

