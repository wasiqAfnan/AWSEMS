import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { Search, Plus, LogOut } from 'lucide-react'
import {
  getEmployees,
  searchEmployees,
  deleteEmployee,
  logoutUrl,
} from '../services/api'
import EmployeeTable from '../components/EmployeeTable'
import EmployeeModal from '../components/EmployeeModal'
import ConfirmDialog from '../components/ConfirmDialog'
import Toast from '../components/Toast'

export default function EmployeesPage() {
  const navigate = useNavigate()

  const [employees, setEmployees] = useState([])
  const [loading, setLoading] = useState(false)
  const [query, setQuery] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [editEmployee, setEditEmployee] = useState(null)

  const [confirmDelete, setConfirmDelete] = useState(null) // empId to delete
  const [toast, setToast] = useState(null) // { type, message }

  const showToast = (type, message) => {
    setToast({ type, message })
    setTimeout(() => setToast(null), 3500)
  }

  const fetchEmployees = useCallback(async () => {
    setLoading(true)
    try {
      const res = await getEmployees()
      setEmployees(res.data?.data ?? [])
    } catch (err) {
      if (err.response?.status === 401) navigate('/login')
      else showToast('error', 'Failed to load employees')
    } finally {
      setLoading(false)
    }
  }, [navigate])

  useEffect(() => {
    fetchEmployees()
  }, [fetchEmployees])

  const handleSearch = async (e) => {
    e.preventDefault()
    if (!query.trim()) return fetchEmployees()
    setLoading(true)
    try {
      const res = await searchEmployees(query.trim())
      setEmployees(res.data?.data ?? [])
    } catch {
      showToast('error', 'Search failed')
    } finally {
      setLoading(false)
    }
  }

  const handleSearchClear = () => {
    setQuery('')
    fetchEmployees()
  }

  const handleEdit = (employee) => {
    setEditEmployee(employee)
    setShowModal(true)
  }

  const handleDeleteConfirm = (empId) => setConfirmDelete(empId)

  const handleDelete = async () => {
    try {
      await deleteEmployee(confirmDelete)
      showToast('success', 'Employee deleted')
      fetchEmployees()
    } catch {
      showToast('error', 'Failed to delete employee')
    } finally {
      setConfirmDelete(null)
    }
  }

  const handleModalClose = (didSave) => {
    setShowModal(false)
    setEditEmployee(null)
    if (didSave) {
      fetchEmployees()
      showToast('success', editEmployee ? 'Employee updated' : 'Employee created')
    }
  }

  const handleLogout = () => {
    window.location.href = logoutUrl
  }

  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-800">AWSEMS</h1>
        <button
          onClick={handleLogout}
          className="flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-800 transition-colors"
        >
          <LogOut size={16} />
          Logout
        </button>
      </header>

      <main className="max-w-7xl mx-auto px-6 py-8">
        {/* Title row */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold text-slate-800">Employees</h2>
          <button
            onClick={() => { setEditEmployee(null); setShowModal(true) }}
            className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors"
          >
            <Plus size={16} />
            Add Employee
          </button>
        </div>

        {/* Search bar */}
        <form onSubmit={handleSearch} className="flex gap-2 mb-6">
          <div className="relative flex-1 max-w-md">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search employees…"
              className="w-full pl-9 pr-4 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-400"
            />
          </div>
          <button
            type="submit"
            className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm rounded-lg transition-colors"
          >
            Search
          </button>
          {query && (
            <button
              type="button"
              onClick={handleSearchClear}
              className="px-4 py-2 border border-gray-300 text-sm rounded-lg text-slate-600 hover:bg-gray-50"
            >
              Clear
            </button>
          )}
        </form>

        {/* Table */}
        <EmployeeTable
          employees={employees}
          loading={loading}
          onEdit={handleEdit}
          onDelete={handleDeleteConfirm}
        />
      </main>

      {/* Modals */}
      {showModal && (
        <EmployeeModal
          employee={editEmployee}
          onClose={handleModalClose}
          showToast={showToast}
        />
      )}
      {confirmDelete && (
        <ConfirmDialog
          message={`Delete employee ${confirmDelete}? This cannot be undone.`}
          onConfirm={handleDelete}
          onCancel={() => setConfirmDelete(null)}
        />
      )}
      {toast && <Toast type={toast.type} message={toast.message} />}
    </div>
  )
}

