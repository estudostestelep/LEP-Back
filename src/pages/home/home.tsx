import { useAuth } from "../../context/authContext";

export default function Profile() {
  const { user } = useAuth();

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Profile</h1>
      <pre>{JSON.stringify(user, null, 2)}</pre>
    </div>
  );
}
