import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import AuthLayout from '../layouts/AuthLayout';
import { registerUser } from '../utils/api';

const Register = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const response = await registerUser(email, password);
      if (response.message === "User registered successfully") {
        navigate('/login');
      } else {
        setError(response.message || 'Error during registration');
      }
    } catch (err) {
      setError('Network or server error. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout title="Create Your Account">
      <form 
        onSubmit={handleSubmit} 
        style={{ display: 'flex', flexDirection: 'column', gap: '1rem', marginTop: '1rem' }}
      >
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
        
        {error && (
          <p style={{ color: 'red', fontSize: '0.9rem', textAlign: 'center' }}>{error}</p>
        )}

        <button 
          className="btn btn-primary" 
          type="submit"
          disabled={loading}
          style={{ padding: '0.8rem', fontSize: '1rem' }}
        >
          {loading ? 'Registering...' : 'Register'}
        </button>
      </form>

      <p style={{ marginTop: '1.5rem', textAlign: 'center', fontSize: '0.95rem' }}>
        Already have an account? <Link to="/login" style={{ color: '#2563eb', fontWeight: '500' }}>Sign in</Link>
      </p>
    </AuthLayout>
  );
};

export default Register;
