import './styles.css'

export const dynamic = 'force-dynamic'
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  if (process.versions.bun) {
    console.log('Bun is running')
    // this code will only run when the file is run with Bun
  } else {
    console.log('Node is running')
  }
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
