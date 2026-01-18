import { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext.jsx';
import { useNavigate } from 'react-router-dom';
import { fetchLinks, deleteUrl } from '../utils/api.js';
import AllLinks from '../components/AllLinks.jsx';
import ShortenUrl from '../components/ShortenUrl.jsx';
import DashboardLayout from '../layouts/DashboardLayout.jsx';

const Dashboard = () => {
  const { setAuth } = useContext(AuthContext);
  const [links, setLinks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const token = localStorage.getItem('token');
  const navigate = useNavigate();

useEffect(() => {
  const getLinks = async () => {
    setLoading(true);
    setError('');
    try {
      const fetchedLinks = await fetchLinks(token);
      const normalizedLinks = fetchedLinks.map(link => ({
        shortUrl: link.shorted,
        originalUrl: link.link,
        createdAt: link.createdAt
      }));
      setLinks(normalizedLinks);
    } catch {
      setError('Failed to fetch links');
    } finally {
      setLoading(false);
    }
  };
  getLinks();
}, [token]);

const handleDelete = async (shortUrl) => {
  try {
    const result = await deleteUrl(shortUrl, token); // deleteUrl expects the short code
    if (result.error) alert(result.error);
    else setLinks(prev => prev.filter(link => link.shortUrl !== shortUrl));
  } catch {
    alert('Failed to delete URL');
  }
};


  const handleLogout = () => {
    localStorage.removeItem('token');
    setAuth(false);
    navigate('/login');
  };


  return (
    <DashboardLayout onLogout={handleLogout}>
      <div className="dashboard-container">

        <ShortenUrl setLinks={setLinks} />

        {loading && <p>Loading links...</p>}
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <AllLinks links={links} onDelete={handleDelete} />
      </div>
    </DashboardLayout>
  );
};

export default Dashboard;
