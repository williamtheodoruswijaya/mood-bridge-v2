import Sidebar from "~/components/sidebar";
import Navbar from "~/components/navbar";

export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex flex-1 flex-col">
        <Navbar />
        <main className="flex-1 overflow-auto bg-[#E9FEFF] p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
