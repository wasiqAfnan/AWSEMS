import { Pencil, Trash2 } from 'lucide-react'

export default function EmployeeTable({ employees, loading, onEdit, onDelete }) {
  if (loading) {
    return (
      <div className="text-center py-20 text-slate-400 text-sm">Loading…</div>
    )
  }

  if (!employees || employees.length === 0) {
    return (
      <div className="text-center py-20 text-slate-400 text-sm">
        No employees found.
      </div>
    )
  }

  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200 text-sm">
        <thead className="bg-gray-50">
          <tr>
            {['Emp ID', 'Name', 'Email', 'Contact', 'Role', 'Department', 'Salary', ''].map(
              (h) => (
                <th
                  key={h}
                  className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >
                  {h}
                </th>
              )
            )}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {employees.map((emp) => (
            <tr key={emp.empId} className="hover:bg-gray-50 transition-colors">
              <td className="px-4 py-3 font-mono text-indigo-700">{emp.empId}</td>
              <td className="px-4 py-3 font-medium text-slate-800">{emp.name}</td>
              <td className="px-4 py-3 text-slate-600">{emp.email}</td>
              <td className="px-4 py-3 text-slate-600">{emp.contactNo}</td>
              <td className="px-4 py-3 text-slate-600">{emp.role}</td>
              <td className="px-4 py-3 text-slate-600">{emp.department}</td>
              <td className="px-4 py-3 text-slate-600">
                {typeof emp.salary === 'number'
                  ? `₹${emp.salary.toLocaleString()}`
                  : emp.salary}
              </td>
              <td className="px-4 py-3">
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => onEdit(emp)}
                    className="p-1.5 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded transition-colors"
                    title="Edit"
                  >
                    <Pencil size={15} />
                  </button>
                  <button
                    onClick={() => onDelete(emp.empId)}
                    className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
                    title="Delete"
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

