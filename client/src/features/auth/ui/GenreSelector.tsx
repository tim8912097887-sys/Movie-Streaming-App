import type { GenresData } from "../../../shared/schema/data/genres";

type GenreSelectorProps = {
  genres: GenresData;
  selectedGenres: GenresData;
  onChange: (genres: GenresData) => void;
};

const GenreSelector = ({
  genres,
  selectedGenres,
  onChange,
}: GenreSelectorProps) => {
  const toggleGenre = (id: number) => {
    if (selectedGenres.some((genre) => genre.genre_id === id)) {
      onChange(selectedGenres.filter((genre) => genre.genre_id !== id));
      return;
    }

    onChange([
      ...selectedGenres,
      {
        genre_id: id,
        name: genres.find((genre) => genre.genre_id === id)?.name || "",
      },
    ]);
  };

  return (
    <div className="flex flex-col gap-5">
      <label className="font-medium text-gray-800">Favorite Genres</label>

      <div className="flex flex-wrap gap-2">
        {genres.map((genre) => {
          const selected = selectedGenres.some(
            (selectedGenre) => selectedGenre.genre_id === genre.genre_id,
          );

          return (
            <button
              key={genre.genre_id}
              type="button"
              onClick={() => toggleGenre(genre.genre_id)}
              className={`
                rounded-full
                border
                px-4
                py-2
                text-sm
                font-medium
                transition-all
                duration-200
                hover:scale-105
                ${
                  selected
                    ? "border-blue-600 bg-blue-600 text-white shadow"
                    : "border-gray-300 bg-white text-gray-700 hover:border-blue-400 hover:bg-blue-50"
                }
              `}
            >
              {genre.name}
            </button>
          );
        })}
      </div>
    </div>
  );
};

export default GenreSelector;
