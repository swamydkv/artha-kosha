import type { Metadata } from "next"
import "./globals.css"

export const metadata: Metadata = {
  title: "ArthaKosha",
  description: "Authentication MVP for ArthaKosha",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
