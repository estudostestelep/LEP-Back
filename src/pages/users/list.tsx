// src/pages/Users/List.tsx
import { useEffect, useState } from "react";
import { userService, User } from "@/api/userService";
import FormModal from "@/components/formModal";
import ConfirmModal from "@/components/confirmModal";

// (Removed duplicate User interface, using imported User type)

export default function UsersList() {
  const [users, setUsers] = useState<User[]>([]);
  const [selected, setSelected] = useState<User | null>(null);
  const [formOpen, setFormOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  const fetchData = async () => {
    const res = await userService.getAll();
    setUsers(res.data);
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleSave = async (data: Partial<User>) => {
    if (selected) {
      // Merge selected user with new data to ensure all required fields are present
      await userService.update(selected.id, { ...selected, ...data });
    } else {
      await userService.create({
        name: data.name ?? "",
        email: data.email ?? "",
        role: data.role ?? "",
      });
    }
    setFormOpen(false);
    setSelected(null);
    fetchData();
  };

  const handleDelete = async () => {
    if (selected) {
      await userService.remove(selected.id);
      setConfirmOpen(false);
      setSelected(null);
      fetchData();
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-xl font-bold">Users</h1>
        <button
          onClick={() => {
            setSelected(null);
            setFormOpen(true);
          }}
          className="bg-blue-600 text-white px-3 py-1.5 rounded-md hover:bg-blue-500"
        >
          Add User
        </button>
      </div>

      <table className="w-full border-collapse border border-gray-200">
        <thead>
          <tr className="bg-gray-100 text-left">
            <th className="border border-gray-200 px-3 py-2">Name</th>
            <th className="border border-gray-200 px-3 py-2">Email</th>
            <th className="border border-gray-200 px-3 py-2">Role</th>
            <th className="border border-gray-200 px-3 py-2">Permission</th>
            <th className="border border-gray-200 px-3 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id} className="hover:bg-gray-50">
              <td className="border px-3 py-2">{u.name}</td>
              <td className="border px-3 py-2">{u.email}</td>
              <td className="border px-3 py-2">{u.role}</td>
              <td className="border px-3 py-2 space-x-2">
                <button
                  onClick={() => {
                    setSelected(u);
                    setFormOpen(true);
                  }}
                  className="text-blue-600 hover:underline"
                >
                  Edit
                </button>
                <button
                  onClick={() => {
                    setSelected(u);
                    setConfirmOpen(true);
                  }}
                  className="text-red-600 hover:underline"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Form Modal */}
      <FormModal
        open={formOpen}
        title={selected ? "Edit User" : "Add User"}
        onClose={() => setFormOpen(false)}
        onSubmit={handleSave}
        fields={[
          { name: "name", label: "Name", type: "text" },
          { name: "email", label: "Email", type: "email" },
          { name: "role", label: "Role", type: "text" },
          {
            name: "permissionLevel",
            label: "Permission Level",
            type: "number",
          },
        ]}
        initialValues={{
          name: selected?.name || "",
          email: selected?.email || "",
          role: selected?.role || "",
        }}
      />

      {/* Confirm Modal */}
      <ConfirmModal
        open={confirmOpen}
        title="Delete User"
        message={`Are you sure you want to delete ${selected?.name}?`}
        onCancel={() => setConfirmOpen(false)}
        onConfirm={handleDelete}
      />
    </div>
  );
}
