import * as React from "react";

type SelectOption = {
  value: string | number;
  label: string;
};

type Field = {
  name: string;
  label: string;
  type?: 'text' | 'number' | 'select' | 'checkbox' | string;
  required?: boolean;
  options?: SelectOption[];
};

type Props = {
  title: string;
  open: boolean;
  onClose: () => void;
  fields: Field[];
  initialValues?: Record<string, any>;
  onSubmit: (values: Record<string, any>) => Promise<void> | void;
};


export default function FormModal({ title, open, onClose, fields, initialValues = {}, onSubmit }: Props) {
  const [values, setValues] = React.useState<Record<string, any>>(initialValues);
  const [loading, setLoading] = React.useState(false);


  React.useEffect(() => setValues(initialValues), [initialValues, open]);


  const handleChange = (name: string, v: any) => setValues(prev => ({ ...prev, [name]: v }));


  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await onSubmit(values);
      onClose();
    } catch (err) {
      console.error(err);
      alert((err as any)?.message || "Error");
    } finally {
      setLoading(false);
    }
  };


  if (!open) return null;


  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="bg-white rounded-lg shadow-lg w-full max-w-2xl mx-4">
        <form onSubmit={submit}>
          <div className="flex items-center justify-between px-6 py-4 border-b">
            <h3 className="text-lg font-semibold">{title}</h3>
            <button type="button" onClick={onClose} className="text-gray-500 hover:text-gray-800">✕</button>
          </div>


          <div className="p-6 grid gap-4 grid-cols-1 md:grid-cols-2">
            {fields.map(f => (
              <div key={f.name} className="flex flex-col">
                <label className="text-sm font-medium mb-1">{f.label}{f.required ? ' *' : ''}</label>


                {f.type === 'select' ? (
                    <select
                    value={values[f.name] ?? ''}
                    onChange={e => handleChange(f.name, e.target.value)}
                    className="border rounded px-3 py-2"
                    required={f.required}
                    >
                    <option value="">Select</option>
                    {f.options?.map((o: SelectOption) => (
                      <option key={o.value} value={o.value}>{o.label}</option>
                    ))}
                    </select>
                ) : f.type === 'checkbox' ? (
                  <label className="inline-flex items-center space-x-2">
                    <input type="checkbox" checked={!!values[f.name]} onChange={e => handleChange(f.name, e.target.checked)} />
                    <span className="text-sm">{f.label}</span>
                  </label>
                ) : (
                  <input
                    type={f.type ?? 'text'}
                    value={values[f.name] ?? ''}
                    onChange={e => handleChange(f.name, e.target.value)}
                    className="border rounded px-3 py-2"
                    required={f.required}
                  />
                )}


              </div>
            ))}
          </div>


          <div className="flex justify-end gap-3 px-6 py-4 border-t">
            <button type="button" onClick={onClose} className="px-4 py-2 rounded bg-gray-100">Cancel</button>
            <button type="submit" className="px-4 py-2 rounded bg-blue-600 text-white" disabled={loading}>{loading ? 'Saving...' : 'Save'}</button>
          </div>
        </form>
      </div>
    </div>
  );
}