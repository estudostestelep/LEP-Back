import { useState } from "react";
import { productService, Product } from "@/api/productService";

interface Props {
  initialData?: Product;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function ProductForm({ initialData, onSuccess, onCancel }: Props) {
  const [form, setForm] = useState<Product>(
    initialData || { name: "", price: 0, description: "" }
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: name === "price" ? parseFloat(value) : value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.id) {
      await productService.update(form.id, form);
    } else {
      await productService.create(form);
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
        placeholder="Product name"
        className="w-full border px-3 py-2 rounded"
      />
      <input
        type="number"
        name="price"
        value={form.price}
        onChange={handleChange}
        placeholder="Price"
        className="w-full border px-3 py-2 rounded"
      />
      <input
        type="text"
        name="description"
        value={form.description}
        onChange={handleChange}
        placeholder="Description"
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
