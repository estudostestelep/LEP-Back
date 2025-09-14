import api from "./api";

export interface Customer {
  id?: string;
  name: string;
  phone: string;
  email?: string;
}

export const customerService = {
  getAll: () => api.get<Customer[]>("/customers"),
  getById: (id: string) => api.get<Customer>(`/customers/${id}`),
  create: (data: Customer) => api.post<Customer>("/customers", data),
  update: (id: string, data: Customer) => api.put<Customer>(`/customers/${id}`, data),
  remove: (id: string) => api.delete(`/customers/${id}`),
};