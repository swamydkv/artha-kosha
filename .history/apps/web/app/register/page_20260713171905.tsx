"use client"

import { FormEvent, useState } from "react"

export default function RegisterPage() {
  const [status, setStatus] = useState("")

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const payload = Object.fromEntries(formData.entries())

    const response = await fetch("http://localhost:8080/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    })

    if (response.ok) {
      setStatus("Account created successfully. You can now sign in.")
      return
    }

    setStatus("Registration failed. Please review the entered values.")
  }

  return (
    <main style={{ maxWidth: 480, margin: '0 auto', padding: 24, fontFamily: 'sans-serif' }}>
      <h1>Create Account</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: 12 }}>
        <input name="full_name" placeholder="Full Name" required />
        <input name="date_of_birth" placeholder="Date of Birth" type="date" required />
        <input name="mobile_number" placeholder="Mobile Number" required />
        <input name="email" placeholder="Email" type="email" required />
        <input name="username" placeholder="Username" required />
        <input name="password" placeholder="Password" type="password" required />
        <input name="confirm_password" placeholder="Confirm Password" type="password" required />
        <button type="submit">Create Account</button>
      </form>
      {status ? <p>{status}</p> : null}
    </main>
  )
}
