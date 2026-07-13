import MoviePresenter, {
  type Movie,
} from "../components/layout/movie/MoviePresenter";
import PageContainer from "../components/ui/PageContainer";

const movies: Movie[] = [
  {
    title: "Neon Horizon",
    description:
      "A cybernetic detective uncovers a conspiracy that shakes the foundations of a rain-soaked, futuristic metropolis.",
    image: "https://images.unsplash.com/photo-1578301978693-85fa9c0320b9",
    rating: 8.4,
  },
  {
    title: "Whispers of the Oak",
    description:
      "Two siblings discover an ancient, enchanted forest behind their grandmother's house that holds secrets of the past.",
    image: "https://images.unsplash.com/photo-1511497584788-876760111969",
    rating: 7.9,
  },
  {
    title: "Velocity Zero",
    description:
      "An adrenaline-fueled thriller about an elite heist crew trapped in a hyper-loop train traveling at maximum speed.",
    image: "https://images.unsplash.com/photo-1532103054090-334e6e60b73c",
    rating: 6.8,
  },
  {
    title: "The Last Alchemist",
    description:
      "In an alternate 19th century, a rogue scholar races against time to find the elixir of life before an empire claims it.",
    image: "https://images.unsplash.com/photo-1516979187457-637abb4f9353",
    rating: 8.1,
  },
];

const RecommandationPage = () => {
  return (
    <PageContainer>
      <MoviePresenter movies={movies} />
    </PageContainer>
  );
};

export default RecommandationPage;
