import { useState } from "react";
import { userService, User } from "@/api/userService";

interface Props {
  initialData?: User;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function UserForm({ initialData, onSuccess, onCancel }: Props) {
  const [form, setForm] = useState<User>(
    initialData || { name: "", email: "", role: "" }
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.id) {
      await userService.update(form.id, form);
    } else {
      await userService.create(form);
    }
    onSuccess();
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white p-4 rounded shadow space-y-4"
    >
      <input
        type="text"
        name="name"
        value={form.name}
        onChange={handleChange}
        placeholder="Name"
        className="w-full border px-3 py-2 rounded"
      />
      <input
        type="email"
        name="email"
        value={form.email}
        onChange={handleChange}
        placeholder="Email"
        className="w-full border px-3 py-2 rounded"
      />
      <input
        type="text"
        name="role"
        value={form.role}
        onChange={handleChange}
        placeholder="Role"
        className="w-full border px-3 py-2 rounded"
      />

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
        >
          {form.id ? "Update" : "Create"}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="bg-gray-400 hover:bg-gray-500 text-white px-4 py-2 rounded"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
