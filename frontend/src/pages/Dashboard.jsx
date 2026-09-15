import { useState, useEffect } from 'react';
import { employeeApi } from '../api/employeeApi';

// Components
import Header from '../components/Header';
import Toast from '../components/Toast';
import SearchBar from '../components/SearchBar';
import EmployeeTable from '../components/EmployeeTable';
import EmployeeModal from '../components/EmployeeModal';

const initialEmployeeState = {
  empId: '',
  name: '',
  email: '',
  contactNo: '',
  role: '',
  department: '',
  salary: ''
};

export default function Dashboard({ auth, onSignOut }) {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  
  // Toast state
  const [toast, setToast] = useState({ message: '', type: 'success', visible: false });

  const showToast = (message, type = 'success') => {
    setToast({ message, type, visible: true });
    setTimeout(() => {
      setToast(prev => ({ ...prev, visible: false }));
    }, 3000);
  };
  
  // Modal state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [currentEmployee, setCurrentEmployee] = useState(initialEmployeeState);
  const [originalEmployee, setOriginalEmployee] = useState(null);
  const [modalError, setModalError] = useState(null);

  useEffect(() => {
    fetchEmployees(); // 1. Fetch employees on mount
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const data = await employeeApi.getAll();
      setEmployees(Array.isArray(data) ? data : data.data || data.body || []);
    } catch (err) {
      console.error(err);
      showToast('Failed to load employees.', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!searchQuery.trim()) {
      fetchEmployees();
      return;
    }
    
    try {
      setLoading(true);
      const data = await employeeApi.search(searchQuery);
      setEmployees(Array.isArray(data) ? data : data.data || data.body || []);
    } catch (err) {
      console.error(err);
      showToast('Search failed.', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleClearSearch = () => {
    setSearchQuery('');
    fetchEmployees();
  };

  const openAddModal = () => {
    setIsEditing(false);
    setOriginalEmployee(null);
    setCurrentEmployee(initialEmployeeState);
    setModalError(null);
    setIsModalOpen(true);
  };

  const openEditModal = (emp) => {
    setIsEditing(true);
    const mappedEmp = {
      empId: emp.empId || '',
      name: emp.name || '',
      email: emp.email || '',
      contactNo: emp.contactNo || '',
      role: emp.role || '',
      department: emp.department || '',
      salary: emp.salary || ''
    };
    setOriginalEmployee(mappedEmp);
    setCurrentEmployee(mappedEmp);
    setModalError(null);
    setIsModalOpen(true);
  };

  const validateForm = (payload) => {
    if (!payload.empId.startsWith("EMP") || payload.empId.length < 6 || !/^\d+$/.test(payload.empId.slice(3))) {
      return "Employee ID must start with 'EMP' followed by digits and be at least 6 characters long.";
    }
    if (!payload.name.trim()) return "Name cannot be empty.";
    if (!/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/.test(payload.email)) {
      return "Invalid email format.";
    }
    if (!/^\d{10,12}$/.test(payload.contactNo)) {
      return "Contact number must be between 10 and 12 digits long.";
    }
    if (!payload.role.trim()) return "Role cannot be empty.";
    if (!payload.department.trim()) return "Department cannot be empty.";
    if (payload.salary <= 0) return "Salary must be greater than 0.";
    return null; // No errors
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setModalError(null);
    try {
      const payload = {
        ...currentEmployee,
        salary: parseFloat(currentEmployee.salary)
      };

      const errorMsg = validateForm(payload);
      if (errorMsg) {
        setModalError(errorMsg);
        return;
      }
      
      if (isEditing) {
        // Calculate partial payload based on changes
        const updatePayload = {};
        Object.keys(payload).forEach(key => {
          if (key !== 'empId' && payload[key] !== originalEmployee[key]) {
             updatePayload[key] = payload[key];
          }
        });

        if (Object.keys(updatePayload).length === 0) {
           setModalError('No changes detected.');
           return;
        }

        await employeeApi.update(currentEmployee.empId, updatePayload);
        showToast('Employee updated successfully!');
      } else {
        await employeeApi.create(payload);
        showToast('Employee added successfully!');
      }
      
      setCurrentEmployee(initialEmployeeState);
      setOriginalEmployee(null);
      setIsModalOpen(false);
      fetchEmployees();
    } catch (err) {
      console.error(err);
      setModalError(err.response?.data?.message || err.response?.data?.error || `Failed to ${isEditing ? 'update' : 'create'} employee.`);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this employee?')) return;
    try {
      await employeeApi.delete(id);
      showToast('Employee deleted successfully!');
      fetchEmployees();
    } catch (err) {
      console.error(err);
      showToast('Failed to delete employee.', 'error');
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCurrentEmployee(prev => ({ ...prev, [name]: value }));
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 relative">
      <Toast 
        toast={toast} 
        onClose={() => setToast(prev => ({ ...prev, visible: false }))} 
      />

      <Header 
        userEmail={auth.user?.profile?.email} 
        onSignOut={onSignOut} 
      />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        <SearchBar 
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          onSearch={handleSearch}
          onClear={handleClearSearch}
          onAddEmployee={openAddModal}
        />

        <EmployeeTable 
          employees={employees}
          loading={loading}
          onEdit={openEditModal}
          onDelete={handleDelete}
        />

        <EmployeeModal 
          isOpen={isModalOpen}
          isEditing={isEditing}
          currentEmployee={currentEmployee}
          handleInputChange={handleInputChange}
          handleSave={handleSave}
          onClose={() => setIsModalOpen(false)}
          modalError={modalError}
        />
      </main>
    </div>
  );
}
