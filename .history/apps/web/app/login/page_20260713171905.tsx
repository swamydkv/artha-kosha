"use client"

import { FormEvent, useState } from "react"

export default function LoginPage() {
  const [status, setStatus] = useState("")

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const payload = Object.fromEntries(formData.entries())

    const response = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })

    if (response.ok) {
      const data = await response.json()
      setStatus(data.welcome_message || "Logged in successfully")
      return
    }

    setStatus("Authentication failed. Please check your credentials.")
  }

  return (
    <main style={{ maxWidth: 480, margin: '0 auto', padding: 24, fontFamily: 'sans-serif' }}>
      <h1>Login</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: 12 }}>
        <input name="username" placeholder="Username" required />
        <input name="password" placeholder="Password" type="password" required />
        <button type="submit">Login</button>
      </form>
      {status ? <p>{status}</p> : null}
    </main>
  )
}
