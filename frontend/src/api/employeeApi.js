import api from '../utils/api';

export const employeeApi = {
  // Get all employees
  getAll: async () => {
    const response = await api.get('/employees');
    return response.data;
  },

  // Search employees
  search: async (query) => {
    const response = await api.get('/employees/search', { params: { q: query } });
    return response.data;
  },

  // Create a new employee
  create: async (employeeData) => {
    const response = await api.post('/employees', employeeData);
    return response.data;
  },

  // Update an employee
  update: async (empId, updateData) => {
    const response = await api.patch(`/employees/${empId}`, updateData);
    return response.data;
  },

  // Delete an employee
  delete: async (empId) => {
    const response = await api.delete(`/employees/${empId}`);
    return response.data;
  }
};

