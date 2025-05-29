import Image from "next/image";
import News_1 from "~/assets/news/news-1.jpg";
import News_2 from "~/assets/news/news-2.jpg";
import News_3 from "~/assets/news/news-3.jpg";
import News_4 from "~/assets/news/news-4.jpg";
import News_5 from "~/assets/news/news-5.jpg";

const articles = [
  {
    title: "Be Grateful, They Said",
    author: "Creeping between the Margins",
    excerpt:
      "Growing up, everything in my life was picture-perfect. While other parents became the source of income for divorce lawyers, mine flirted with each other in front of me...",
    url: "https://medium.com/the-virago/be-grateful-they-said-93ad0369ea5d",
    image: News_1,
  },
  {
    title: "What One Leap Into a Lake Taught Me About Living",
    author: "Jakob Ryce",
    excerpt:
      "We all live with limits — limits of what we can take, what we can live with, and how we adapt. Over time, we get so used to our hang-ups that they become the default...",
    url: "https://medium.com/interior-salt/what-one-leap-into-a-lake-taught-me-about-living-d6dce77727ff",
    image: News_2,
  },
  {
    title: "The Empty Chair: Its Role in My Recovery",
    author: "Tom Gavea",
    excerpt:
      "So you walk into a therapist’s office, and there’s an empty chair sitting across from you. You don’t quite understand that the empty chair might just hold the key...",
    url: "https://medium.com/black-bear-recovery/the-empty-chair-its-role-in-my-recovery-4eccf8847bfb",
    image: News_3,
  },
  {
    title:
      "The Expectations Gap: On The Unbearable Weight of (Mindful) Modern Parenting and Managing",
    author: "Carin-Isabel Knoop",
    excerpt:
      "Parents, like managers, are emotional shock absorbers in an increasingly emotionally immature world. They must remain calm in chaos, mediate conflict...",
    url: "https://carinisabelknoop.medium.com/the-expectations-gap-on-the-unbearable-weight-of-mindful-modern-parenting-and-managing-f87b61145469",
    image: News_4,
  },
  {
    title: "When Did UX & Content Get So Hard?",
    author: "Erin Schroeder",
    excerpt:
      "It's a weekday morning and I’m sipping coffee, scanning my calendar for my meetings today, preparing my work, swimming in a slog of newsletters...",
    url: "https://uxdesign.cc/when-did-ux-content-get-so-hard-25edd66ab081",
    image: News_5,
  },
];

export default function News() {
  return (
    <div className="mx-auto max-w-md space-y-4 rounded-xl bg-[#b7f9fd] p-6 shadow-lg">
      <h1 className="text-xl font-semibold text-gray-800">Community Trends</h1>
      {articles.slice(0, 5).map((article, index) => (
        <a
          key={index}
          href={article.url}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-start gap-4 rounded-lg p-2 transition hover:bg-[#f3feff]"
        >
          <Image
            src={article.image}
            alt={article.title}
            className="h-16 w-16 flex-shrink-0 rounded object-cover"
          />
          <div className="min-w-0">
            <h2 className="max-w-[200px] truncate text-sm font-semibold text-gray-800">
              {article.title}
            </h2>
            <p className="line-clamp-2 text-xs text-gray-600">
              {article.excerpt}
            </p>
          </div>
        </a>
      ))}
    </div>
  );
}
