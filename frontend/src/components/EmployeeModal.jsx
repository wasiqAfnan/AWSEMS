import { useState } from 'react'
import { X } from 'lucide-react'
import { createEmployee, updateEmployee } from '../services/api'

const EMPTY_FORM = {
  empId: '',
  name: '',
  email: '',
  contactNo: '',
  role: '',
  department: '',
  salary: '',
}

export default function EmployeeModal({ employee, onClose, showToast }) {
  const isEdit = Boolean(employee)

  const [form, setForm] = useState(
    isEdit
      ? {
          empId: employee.empId ?? '',
          name: employee.name ?? '',
          email: employee.email ?? '',
          contactNo: employee.contactNo ?? '',
          role: employee.role ?? '',
          department: employee.department ?? '',
          salary: employee.salary ?? '',
        }
      : EMPTY_FORM
  )
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState(null)

  const handleChange = (e) =>
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }))

  const validateForm = () => {
    const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/
    const validateString = (field, value) => (!value || String(value).trim() === '' ? `${field} must not be empty` : null)
    
    const validateEmail = (value) => {
      const v = String(value).trim()
      if (!v) return 'email must not be empty'
      if (!emailRegex.test(v)) return 'email is not a valid email address'
      return null
    }

    const validateContactNo = (value) => {
      const v = String(value).trim()
      if (!v) return 'contactNo must not be empty'
      if (!/^\d+$/.test(v)) return 'contactNo must contain only digits'
      if (v.length < 10) return 'contactNo must be at least 10 digits'
      if (v.length > 12) return 'contactNo must not exceed 12 digits'
      return null
    }

    if (!isEdit) {
      // Validate all fields for create
      const empId = String(form.empId).trim()
      if (!empId) return 'empId must not be empty'
      if (!empId.startsWith('EMP')) return "empId must start with 'EMP'"
      const suffix = empId.slice(3)
      if (!suffix) return "empId must contain digits after 'EMP'"
      if (!/^\d+$/.test(suffix)) return "empId must contain only digits after 'EMP'"
      if (empId.length < 6) return 'empId must be at least 6 characters long'

      const strErr = validateString('name', form.name) || 
                     validateString('role', form.role) || 
                     validateString('department', form.department)
      if (strErr) return strErr

      const emailErr = validateEmail(form.email)
      if (emailErr) return emailErr

      const contactErr = validateContactNo(form.contactNo)
      if (contactErr) return contactErr

      if (parseFloat(form.salary) < 0) return 'salary must be greater than or equal to 0'
    } else {
      // Validate changed fields for update
      if (form.name !== employee.name) {
        const err = validateString('name', form.name)
        if (err) return err
      }
      if (form.email !== employee.email) {
        const err = validateEmail(form.email)
        if (err) return err
      }
      if (form.contactNo !== employee.contactNo) {
        const err = validateContactNo(form.contactNo)
        if (err) return err
      }
      if (form.role !== employee.role) {
        const err = validateString('role', form.role)
        if (err) return err
      }
      if (form.department !== employee.department) {
        const err = validateString('department', form.department)
        if (err) return err
      }
      if (String(form.salary) !== String(employee.salary)) {
        if (parseFloat(form.salary) < 0) return 'salary must be greater than or equal to 0'
      }
    }
    return null
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)

    const validationError = validateForm()
    if (validationError) {
      setError(validationError)
      return
    }

    setSubmitting(true)
    try {
      if (isEdit) {
        // Build only changed optional fields for PATCH
        const patch = {}
        if (form.name !== employee.name) patch.name = form.name
        if (form.email !== employee.email) patch.email = form.email
        if (form.contactNo !== employee.contactNo) patch.contactNo = form.contactNo
        if (form.role !== employee.role) patch.role = form.role
        if (form.department !== employee.department) patch.department = form.department
        if (String(form.salary) !== String(employee.salary))
          patch.salary = parseFloat(form.salary)
        await updateEmployee(employee.empId, patch)
      } else {
        await createEmployee({ ...form, salary: parseFloat(form.salary) })
      }
      onClose(true)
    } catch (err) {
      setError(err.response?.data?.message ?? 'Something went wrong')
    } finally {
      setSubmitting(false)
    }
  }

  const fields = [
    { name: 'empId', label: 'Employee ID', placeholder: 'EMP001', disabled: isEdit },
    { name: 'name', label: 'Full Name', placeholder: 'John Doe' },
    { name: 'email', label: 'Email', placeholder: 'john@example.com', type: 'email' },
    { name: 'contactNo', label: 'Contact No', placeholder: '03001234567' },
    { name: 'role', label: 'Role', placeholder: 'Software Engineer' },
    { name: 'department', label: 'Department', placeholder: 'Engineering' },
    { name: 'salary', label: 'Salary', placeholder: '75000', type: 'number' },
  ]

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
          <h3 className="text-lg font-semibold text-slate-800">
            {isEdit ? 'Edit Employee' : 'Add Employee'}
          </h3>
          <button
            onClick={() => onClose(false)}
            className="p-1.5 text-slate-400 hover:text-slate-600 rounded transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="px-6 py-5 space-y-4">
          {error && (
            <div className="bg-red-50 text-red-700 text-sm px-4 py-2.5 rounded-lg border border-red-200">
              {error}
            </div>
          )}

          <div className="grid grid-cols-2 gap-4">
            {fields.map(({ name, label, placeholder, type = 'text', disabled }) => (
              <div key={name} className={name === 'email' ? 'col-span-2' : ''}>
                <label className="block text-xs font-medium text-slate-600 mb-1">
                  {label}
                </label>
                <input
                  type={type}
                  name={name}
                  value={form[name]}
                  onChange={handleChange}
                  placeholder={placeholder}
                  disabled={disabled}
                  required={!isEdit}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-400 disabled:bg-gray-50 disabled:text-gray-400"
                />
              </div>
            ))}
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={() => onClose(false)}
              className="px-4 py-2 text-sm border border-gray-300 rounded-lg text-slate-600 hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-5 py-2 text-sm bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition-colors disabled:opacity-60"
            >
              {submitting ? 'Saving…' : isEdit ? 'Update' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

