import Sidebar from "~/components/sidebar";
import Navbar from "~/components/navbar";
import Image from "next/image";
import detectMoodBackground from "~/assets/detect-mood.jpg";

export default function MainLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="relative flex flex-1 flex-col">
        <Image
          src={detectMoodBackground}
          alt="Background"
          fill
          className="object-cover z-0"
        />

        <div className="absolute inset-0 bg-black/60 backdrop-blur-[3px] z-10" />



        <div className="relative z-20 flex flex-col h-full">
          <Navbar />
          <main className="flex-1 overflow-auto p-6">
            {children}
          </main>
        </div>
      </div>
    </div>
  );
}
