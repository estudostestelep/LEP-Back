import { Link } from 'react-router-dom';
import { useAuth } from '@/context/authContext';

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <div className="space-x-4">
        <Link to="/" className="hover:underline">Home</Link>
        <Link to="/menu" className="hover:underline">Menu</Link>
        {user && (
          <>
            <Link to="/orders" className="hover:underline">Orders</Link>
            <Link to="/users" className="hover:underline">Users</Link>
            <Link to="/products" className="hover:underline">Products</Link>
          </>
        )}
      </div>

      <div>
        {user ? (
          <button
            onClick={logout}
            className="bg-red-500 hover:bg-red-600 px-4 py-2 rounded"
          >
            Logout
          </button>
        ) : (
          <Link
            to="/login"
            className="bg-green-500 hover:bg-green-600 px-4 py-2 rounded"
          >
            Login
          </Link>
        )}
      </div>
    </nav>
  );
}
