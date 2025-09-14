import * as React from 'react';
import FormModal from '@/components/formModal';
import ConfirmModal from '@/components/confirmModal';
import { customerService, Customer } from '@/api/customerService';


export default function CustomersPage() {
  const [customers, setCustomers] = React.useState<Customer[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [open, setOpen] = React.useState(false);
  const [editing, setEditing] = React.useState<Customer | null>(null);
  const [confirmOpen, setConfirmOpen] = React.useState(false);
  const [toDelete, setToDelete] = React.useState<string | null>(null);


  const load = async () => {
    setLoading(true);
    const res = await customerService.getAll();
    setCustomers(res.data);
    setLoading(false);
  };


  React.useEffect(() => { load(); }, []);


  const handleCreateOrUpdate = async (values: any) => {
    if (editing) {
      await customerService.update(editing.id!, values);
    } else {
      await customerService.create(values);
    }
    await load();
  };


  const handleDelete = async () => {
    if (!toDelete) return;
    await customerService.remove(toDelete);
    setConfirmOpen(false);
    setToDelete(null);
    await load();
  };


  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-xl font-bold">Customers</h1>
        <button onClick={() => { setEditing(null); setOpen(true); }} className="bg-blue-600 text-white px-4 py-2 rounded">New Customer</button>
      </div>


      {loading ? <p>Loading...</p> : (
        <table className="w-full border">
          <thead className="bg-gray-100">
            <tr>
              <th className="p-2 text-left">Name</th>
              <th className="p-2 text-left">Email</th>
              <th className="p-2 text-left">Phone</th>
              <th className="p-2 text-left">Actions</th>
            </tr>
          </thead>
          <tbody>
            {customers.map(c => (
              <tr key={c.id} className="border-t">
                <td className="p-2">{c.name}</td>
                <td className="p-2">{c.email}</td>
                <td className="p-2">{c.phone}</td>
                <td className="p-2">
                  <button onClick={() => { setEditing(c); setOpen(true); }} className="text-blue-600 mr-3">Edit</button>
                  <button onClick={() => { setToDelete(c.id!); setConfirmOpen(true); }} className="text-red-600">Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}


      <FormModal
        title={editing ? 'Edit Customer' : 'New Customer'}
        open={open}
        onClose={() => setOpen(false)}
        fields={[
          { name: 'name', label: 'Name', required: true },
          { name: 'email', label: 'Email', type: 'text' },
          { name: 'phone', label: 'Phone', type: 'text' },
          { name: 'birthday', label: 'Birthday', type: 'date' },
          { name: 'address', label: 'Address', type: 'text' },
        ]}
        initialValues={editing ?? {}}
        onSubmit={handleCreateOrUpdate}
      />


      <ConfirmModal open={confirmOpen} title="Delete Customer" message="Do you want to delete this customer?" onCancel={() => setConfirmOpen(false)} onConfirm={handleDelete} />
    </div>
  );
}