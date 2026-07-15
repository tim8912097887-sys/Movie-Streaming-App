const PaginationSkeleton = () => {
  return (
    <div className="flex overflow-hidden rounded border border-gray-300 shadow-2xs">
      <div className="w-15 h-10 animate-pulse bg-gray-200" />

      {Array.from({ length: 5 }).map((_, index) => (
        <div
          key={index}
          className="w-10 h-10 animate-pulse bg-gray-200 border-l border-gray-300"
        />
      ))}

      <div className="w-15 h-10 animate-pulse bg-gray-200 border-l border-gray-300" />
    </div>
  );
};

export default PaginationSkeleton;
