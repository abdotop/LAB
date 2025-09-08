import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Mandataires - E-Parrainages',
  description: 'Interface mandataires et gestionnaires pour la plateforme de parrainage électoral du Sénégal',
  keywords: ['mandataires', 'élections', 'Sénégal', 'parrainage', 'campagne'],
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
      <body className="font-sans">
        <div id="root">
          {children}
        </div>
      </body>
    </html>
  )
}