export default function Page() {
  const friends = ["Jonathan", "Juto", "Alice"];
  const messages = [
    { sender: "friend", text: "Hi, how are you?" },
    { sender: "user", text: "I'm good, thanks!" },
    { sender: "friend", text: "Wanna grab coffee later?" },
    { sender: "user", text: "Sure!" },
  ];

  return (
    <main className="grid h-screen w-full grid-cols-[1fr_300px] bg-white text-black">
      {/* Chat Box */}
      <section className="flex flex-col justify-between bg-gradient-to-br from-blue-50 to-white p-4">
        <div className="flex-1 space-y-4 overflow-y-auto rounded border p-4 shadow-inner">
          {messages.map((msg, i) => (
            <div
              key={i}
              className={`max-w-xs rounded-lg px-4 py-2 ${
                msg.sender === "user"
                  ? "self-end bg-blue-500 text-white"
                  : "self-start bg-gray-200 text-black"
              }`}
            >
              {msg.text}
            </div>
          ))}
        </div>
        <form className="mt-4 flex gap-2">
          <input
            type="text"
            placeholder="Type a message..."
            className="flex-1 rounded border px-3 py-2"
          />
          <button
            type="submit"
            className="rounded bg-blue-500 px-4 py-2 text-white"
          >
            Send
          </button>
        </form>
      </section>

      {/* Friends List */}
      <aside className="overflow-y-auto border-l bg-sky-100 p-4">
        <h2 className="mb-4 text-xl font-bold">Friends</h2>
        <ul className="space-y-2">
          {friends.map((name, i) => (
            <li
              key={i}
              className="cursor-pointer rounded bg-white p-2 shadow hover:bg-sky-200"
            >
              {name}
            </li>
          ))}
        </ul>
      </aside>
    </main>
  );
}
