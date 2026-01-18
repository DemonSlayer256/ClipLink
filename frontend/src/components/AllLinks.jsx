const AllLinks = ({ links, onDelete }) => {
  if (!links || links.length === 0)
    return <p style={{ textAlign: 'center', marginTop: '1rem' }}>No links yet.</p>;

  return (
    <div className="links-grid" style={{ display: 'grid', gap: '1rem', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))' }}>
      {links.map(link => (
        <div className="card" key={link.shortUrl} style={{ padding: '1rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
          <a href={link.shortUrl} target="_blank" rel="noreferrer" style={{ wordBreak: 'break-all', color: '#2563eb', fontWeight: '500' }}>
            {link.shortUrl}
          </a>
          <span>{link.originalUrl}</span>
          <span>Created at: {link.createdAt}</span>
          <button className="btn btn-danger" onClick={() => onDelete(link.shortUrl)} style={{ marginTop: '0.5rem', alignSelf: 'flex-start' }}>
            Delete
          </button>
        </div>
      ))}
    </div>
  );
};

export default AllLinks;
