import { useState } from 'react';
import { shortenUrl } from '../utils/api';

const ShortenUrl = ({ setLinks }) => {
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const token = localStorage.getItem('token');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (!url) return;

    setLoading(true);
    try {
      const result = await shortenUrl(url, token);
      if (result.error) {
        setError(result.error);
      } else {
setLinks(prev => [
  {
    shortUrl: result.shorted,  
    originalUrl: result.link || url,
    createdAt: result.createdAt
  },
  ...prev
]);

        setUrl('');
      }
    } catch {
      setError('Failed to shorten URL');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="card" style={{ padding: '1rem', display: 'flex', flexDirection: 'column', gap: '0.75rem', marginBottom: '1.5rem' }}>
      <h3>Shorten a URL</h3>
      <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
        <input
          className="input"
          type="url"
          placeholder="https://example.com"
          value={url}
          onChange={e => setUrl(e.target.value)}
          required
          style={{ flex: '1 1 200px', padding: '0.5rem' }}
        />
        <button className="btn btn-primary" type="submit" disabled={loading}>
          {loading ? 'Shortening…' : 'Shorten'}
        </button>
      </div>
      {error && <p style={{ color: 'red', marginTop: '0.5rem' }}>{error}</p>}
    </form>
  );
};

export default ShortenUrl;
