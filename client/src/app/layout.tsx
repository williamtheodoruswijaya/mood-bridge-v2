import "~/styles/globals.css";

import { type Metadata } from "next";
import { Geist } from "next/font/google";
import Sidebar from "~/components/sidebar";

export const metadata: Metadata = {
  title: "Mood Bridge",
  description: "Mental health support by everyone, for everyone.",
  icons: [{ rel: "icon", url: "/icon.png" }],
};

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist-sans",
});

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className={geist.className}>
        <div className="flex h-screen">
          <Sidebar />
          <main className="flex-1 overflow-auto bg-[#E9FEFF] p-6">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
