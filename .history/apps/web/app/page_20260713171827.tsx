export default function HomePage() {
  return (
    <main style={{ maxWidth: 720, margin: '0 auto', padding: 24, fontFamily: 'sans-serif' }}>
      <h1>ArthaKosha</h1>
      <p>Personal & Family Finance Manager</p>
      <p>Welcome to ArthaKosha.</p>
      <div style={{ display: 'grid', gap: 12 }}>
        <a href="/register">Create Account</a>
        <a href="/login">Login</a>
      </div>
    </main>
  )
}
