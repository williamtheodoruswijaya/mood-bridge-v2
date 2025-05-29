import CreatePost from "~/components/create-post";

export default function HomePage() {
  return (
    <main className="flex min-h-screen w-full text-white">
      <section className="flex flex-1 flex-col items-center px-6">
        {/* Bagian Atas */}
        <div className="mt-4 w-full max-w-4xl">
          <CreatePost />
        </div>
        {/* Bagian penuh sama postingan */}
      </section>
      <aside className="flex w-[250px] flex-col gap-4 border-r p-4 backdrop-blur-md">
        {/* What's happening Section */}
      </aside>
    </main>
  );
}
