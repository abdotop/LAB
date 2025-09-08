import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'DGE Admin - E-Parrainages',
  description: 'Interface d\'administration DGE pour la plateforme de parrainage électoral du Sénégal',
  keywords: ['DGE', 'élections', 'Sénégal', 'parrainage', 'administration'],
  authors: [{ name: 'Abdou TOP' }],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="fr">
      <head>
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="theme-color" content="#00853f" />
      </head>
      <body className={inter.className}>
        <div id="root">
          {children}
        </div>
      </body>
    </html>
  )
}