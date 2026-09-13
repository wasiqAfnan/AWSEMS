import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // send cookies (access_token) on every request
  headers: { 'Content-Type': 'application/json' },
})

// Auth API
// Login and Logout redirect via the backend to interact with Cognito,
// so we just navigate the browser there.
export const loginUrl = `${API_BASE_URL}/auth/login`
export const logoutUrl = `${API_BASE_URL}/auth/logout`

console.log('API BASE URL:', API_BASE_URL)
console.log('LOGIN URL:', loginUrl)
console.log('LOGOUT URL:', logoutUrl)


// Employee API
export const getEmployees = () => api.get('/employees')

export const searchEmployees = (q) => api.get('/employees/search', { params: { q } })

export const createEmployee = (data) => api.post('/employees', data)

export const updateEmployee = (empId, data) => api.patch(`/employees/${empId}`, data)

export const deleteEmployee = (empId) => api.delete(`/employees/${empId}`)

export default api
