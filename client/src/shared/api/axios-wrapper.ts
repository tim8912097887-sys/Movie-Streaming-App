export const axiosWrapper = <T, R>(func: (data: T) => Promise<R>) => {
  return async (data: T): Promise<R> => {
    try {
      const response = await func(data);
      return response;
    } catch (error: any) {
      if (error.response) {
        console.log("Server Response Error:", error.response.data);

        // Use optional chaining to guarantee a string error message
        const message =
          error.response.data?.error?.message ||
          error.response.data?.message ||
          "Invalid credentials or request error.";

        throw new Error(message);
      } else if (error.request) {
        const errorMessage = getRequestErrorMessage(error);
        console.log("Request Error", errorMessage);
        throw new Error(errorMessage);
      } else {
        console.log("Error", error.message);
        throw new Error(error.message);
      }
    }
  };
};

function getRequestErrorMessage(error: any): string {
  if (error.code === "ECONNABORTED") {
    return "Request timed out. Please try again.";
  }
  if (error.code === "ERR_NETWORK" || !navigator.onLine) {
    return "Network error. Please check your internet connection.";
  }
  return error.message || "No response received from the server.";
}
