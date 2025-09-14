import { useState } from "react";
import { useAuth } from "../../context/authContext";

export default function Login() {
  const { login } = useAuth();
  const [userId, setUserId] = useState("");
  const [orgId, setOrgId] = useState("");
  const [projId, setProjId] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login({ id: userId, name: "", orgId, projId });
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Login</h1>
      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "1rem", maxWidth: "300px" }}>
        <input placeholder="User ID" value={userId} onChange={(e) => setUserId(e.target.value)} />
        <input placeholder="Org ID" value={orgId} onChange={(e) => setOrgId(e.target.value)} />
        <input placeholder="Project ID" value={projId} onChange={(e) => setProjId(e.target.value)} />
        <button type="submit">Login</button>
      </form>
    </div>
  );
}
