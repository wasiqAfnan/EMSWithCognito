export default function EmployeeTable({ employees, loading, onEdit, onDelete }) {
  if (loading) {
    return (
      <section className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="p-12 text-center text-gray-500">Loading directory...</div>
      </section>
    );
  }

  if (employees.length === 0) {
    return (
      <section className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="p-12 text-center text-gray-500">No employees found.</div>
      </section>
    );
  }

  return (
    <section className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 text-left">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Emp ID</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Name</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Email</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Contact</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Role</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Department</th>
              <th className="px-6 py-4 text-xs font-semibold text-gray-600 uppercase tracking-wider">Salary</th>
              <th className="px-6 py-4 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {employees.map(emp => (
              <tr key={emp.empId || emp.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{emp.empId}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">{emp.name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{emp.email}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{emp.contactNo}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{emp.role}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{emp.department}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">₹{emp.salary?.toLocaleString('en-IN') || 0}</td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-3">
                  <button
                    onClick={() => onEdit(emp)}
                    className="text-blue-600 hover:text-blue-900 cursor-pointer"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(emp.empId || emp.id)}
                    className="text-red-600 hover:text-red-900 cursor-pointer"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}

