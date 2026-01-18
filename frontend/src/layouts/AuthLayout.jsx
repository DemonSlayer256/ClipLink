const AuthLayout = ({ title, children }) => {
  return (
    <div className="auth-container">
      <div className="card auth-card">
        <h1 style={{ textAlign: 'center', marginBottom: '1rem' }}>{title}</h1>
        {children}
      </div>
    </div>
  );
};

export default AuthLayout;
