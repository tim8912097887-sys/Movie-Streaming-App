type ErrorFallbackProps = {
  message: string;
  onRetry: () => void;
};

const ErrorFallback = ({ message, onRetry }: ErrorFallbackProps) => {
  return (
    <div className="flex flex-col items-center justify-center p-8 text-center min-h-100 w-full">
      <div className="text-red-500 text-5xl mb-4">⚠️</div>
      <h3 className="text-xl font-bold text-gray-800 dark:text-white mb-2">
        Oops! Something went wrong
      </h3>
      <p className="text-gray-500 mb-6 max-w-md">{message}</p>
      <button
        onClick={onRetry}
        className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-md transition-all active:scale-95"
      >
        Try Again
      </button>
    </div>
  );
};

export default ErrorFallback;
