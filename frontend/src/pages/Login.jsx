import { useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { useNavigate, Link } from 'react-router-dom';
import AuthLayout from '../layouts/AuthLayout';
import { loginUser } from '../utils/api';

const Login = () => {
  const { setAuth } = useContext(AuthContext);
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const token = await loginUser(email, password);
        if (token) {
            localStorage.setItem('token', token);
            setAuth(true);
            navigate('/dashboard');
        } else {
            alert('Invalid credentials'); // Error handling
        }
    };

  return (
    <AuthLayout title="Sign in">
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        <input
          className="input"
          type="email"
          placeholder="Email address"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
        />
        <input
          className="input"
          type="password"
          placeholder="Password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          required
        />
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button className="btn btn-primary" type="submit">Sign In</button>
      </form>
      <p style={{ marginTop: '1rem', textAlign: 'center' }}>
        Don’t have an account? <Link to="/register">Register</Link>
      </p>
    </AuthLayout>
  );
};

export default Login;
