const DashboardLayout = ({ children, onLogout }) => {
  return (
    <div className="dashboard-container">
      {/* Top Navbar */}
      <header className="navbar">
        <div className="navbar-inner">
          <h1 className="navbar-title">URL Shortener</h1>
          <button className="btn btn-danger" onClick={onLogout}>
            Log out
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="container" style={{ paddingTop: '1.5rem' }}>
        {children}
      </main>
    </div>
  );
};

export default DashboardLayout;
