"use client";
import { useParams } from "next/navigation";

export default function Page() {
  const params = useParams();
  const userID = params.id;
  return (
    <div className="flex h-screen flex-col items-center justify-center">
      <h1 className="mb-4 text-2xl font-bold">Friend Page {userID}</h1>
      <p className="text-gray-600">This is the friend page.</p>
    </div>
  );
}
