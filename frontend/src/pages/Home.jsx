import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <div className="container" style={{ paddingTop: '3rem', paddingBottom: '3rem' }}>
      
      {/* Hero Section */}
      <section className="card" style={{ textAlign: 'center', padding: '3rem 2rem', background: 'linear-gradient(135deg, #2563eb 0%, #60a5fa 100%)', color: 'white', marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2.5rem', marginBottom: '1rem' }}>Shorten Your URLs Instantly</h1>
        <p style={{ fontSize: '1.2rem', marginBottom: '2rem' }}>Turn long, messy links into simple, shareable URLs in seconds. Perfect for social media, emails, and campaigns.</p>
        <nav style={{ display: 'flex', justifyContent: 'center', gap: '1rem', flexWrap: 'wrap' }}>
          <Link className="btn btn-primary" to="/login">Get Started</Link>
          <Link className="btn" style={{ background: 'white', color: '#2563eb' }} to="/register">Sign Up Free</Link>
        </nav>
      </section>

      {/* How it Works Section */}
      <section style={{ textAlign: 'center', marginBottom: '3rem' }}>
        <h2 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>How It Works</h2>
        <div className="links-grid" style={{ maxWidth: '900px', margin: '0 auto' }}>
          <div className="card">
            <h3>1. Paste Your Link</h3>
            <p>Copy any URL you want to shorten and paste it into our easy-to-use form.</p>
          </div>
          <div className="card">
            <h3>2. Shorten Instantly</h3>
            <p>Click “Shorten” and get a compact, shareable URL in seconds.</p>
          </div>
          <div className="card">
            <h3>3. Share Anywhere</h3>
            <p>Use your short link on social media, emails, or campaigns and track clicks.</p>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section style={{ textAlign: 'center', marginBottom: '3rem' }}>
        <h2 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>Features That Save You Time</h2>
        <div className="links-grid" style={{ maxWidth: '900px', margin: '0 auto' }}>
          <div className="card">
            <h3>Dashboard</h3>
            <p>Manage all your short links in one place. Edit, delete, or track performance easily.</p>
          </div>
          <div className="card">
            <h3>Secure</h3>
            <p>All links are secured with your account. Your data is private and protected.</p>
          </div>
          <div className="card">
            <h3>Responsive</h3>
            <p>Use it on desktop, tablet, or mobile — it looks great everywhere.</p>
          </div>
        </div>
      </section>

      {/* Call to Action */}
      <section className="card" style={{ textAlign: 'center', padding: '2rem', background: '#2563eb', color: 'white' }}>
        <h2 style={{ fontSize: '2rem', marginBottom: '1rem' }}>Ready to Shorten Your URLs?</h2>
        <p style={{ marginBottom: '1.5rem' }}>Sign up now and start sharing cleaner, shorter links today.</p>
        <nav style={{ display: 'flex', justifyContent: 'center', gap: '1rem', flexWrap: 'wrap' }}>
          <Link className="btn btn-primary" to="/login">Login</Link>
          <Link className="btn" style={{ background: 'white', color: '#2563eb' }} to="/register">Register Free</Link>
        </nav>
      </section>
    </div>
  );
};

export default Home;
