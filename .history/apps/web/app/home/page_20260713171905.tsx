"use client"

import { useState } from "react"

export default function HomePage() {
  const [status, setStatus] = useState("You are successfully logged in.")

  async function handleLogout() {
    const response = await fetch("http://localhost:8080/logout", {
      method: "POST",
      headers: { "X-Session-ID": "demo-session" },
    })

    if (response.ok) {
      setStatus("You have been logged out successfully.")
    }
  }

  return (
    <main style={{ maxWidth: 480, margin: '0 auto', padding: 24, fontFamily: 'sans-serif' }}>
      <h1>Home</h1>
      <p>Welcome, Jane!</p>
      <p>{status}</p>
      <button onClick={handleLogout}>Logout</button>
    </main>
  )
}
