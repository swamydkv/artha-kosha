import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ArthaKosha - Personal Finance Management",
  description: "Your personal finance management platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}